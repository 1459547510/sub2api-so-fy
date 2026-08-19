package handler

// OpenAI-compatible (Sora-format) video entry points.
//
// Relays such as NewAPI submit multipart/form-data to POST /v1/videos and poll
// GET /v1/videos/{id} expecting an OpenAI video object. This layer translates
// those requests onto the existing async video-job pipeline so billing,
// scheduling, and group permissions stay unchanged. The native JSON contract
// on /v1/videos/generations and /v1/videos/jobs/{id} is untouched.

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type soraCompatRequestError struct {
	status  int
	message string
}

func (e *soraCompatRequestError) Error() string { return e.message }

func newSoraCompatError(status int, format string, args ...any) *soraCompatRequestError {
	return &soraCompatRequestError{status: status, message: fmt.Sprintf(format, args...)}
}

// SoraVideoCompatCreate accepts a Sora-format multipart create request and
// starts an async video job. It always behaves asynchronously regardless of
// the Prefer header and answers with an OpenAI video object.
func (h *OpenAIGatewayHandler) SoraVideoCompatCreate(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
	}
	if apiKey.Group == nil || !service.GroupAllowsImageGeneration(apiKey.Group) {
		h.errorResponse(c, http.StatusForbidden, "permission_error", service.ImageGenerationPermissionMessage())
		return
	}
	if h.videoJobService == nil {
		h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "Video service is not configured")
		return
	}
	if err := c.Request.ParseMultipartForm(64 << 20); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "multipart form is invalid")
		return
	}
	body, err := h.soraCompatRequestBody(c)
	if err != nil {
		var reqErr *soraCompatRequestError
		if errors.As(err, &reqErr) {
			errType := "invalid_request_error"
			if reqErr.status >= http.StatusInternalServerError {
				errType = "api_error"
			}
			h.errorResponse(c, reqErr.status, errType, reqErr.message)
			return
		}
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}
	user := apiKey.User
	if user == nil {
		user = &service.User{ID: subject.UserID}
	}
	subscription, _ := middleware2.GetSubscriptionFromContext(c)
	localInputName := ""
	if h.videoInputHandler != nil && h.videoInputHandler.store != nil {
		localInputName = strings.Join(h.videoInputHandler.store.TokensFromVideoRequest(body), ",")
	}
	job, err := h.videoJobService.Create(c.Request.Context(), service.CreateVideoJobInput{
		APIKey: apiKey, User: user, Subscription: subscription, Body: body, LocalInputName: localInputName,
	})
	if err != nil {
		h.leoVideoCreateErrorResponse(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, soraVideoObject(job))
}

// LeoVideoJobSora answers GET /v1/videos/{vidjob_id} with the OpenAI video
// object shape. The native shape stays on /v1/videos/jobs/{id}.
func (h *OpenAIGatewayHandler) LeoVideoJobSora(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	if h.videoJobService == nil {
		h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "Video service is not configured")
		return
	}
	jobID := strings.TrimSpace(c.Param("request_id"))
	if jobID == "" {
		jobID = strings.TrimSpace(c.Param("job_id"))
	}
	job, err := h.videoJobService.Get(c.Request.Context(), jobID, apiKey.ID)
	if errors.Is(err, service.ErrVideoJobNotFound) {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Video job not found")
		return
	}
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video job unavailable")
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, soraVideoObject(job))
}

// soraCompatRequestBody translates the multipart form into the native JSON
// video create body.
func (h *OpenAIGatewayHandler) soraCompatRequestBody(c *gin.Context) ([]byte, error) {
	form := c.Request
	body := map[string]any{}
	if model := strings.TrimSpace(form.FormValue("model")); model != "" {
		body["model"] = model
	}
	if prompt := strings.TrimSpace(form.FormValue("prompt")); prompt != "" {
		body["prompt"] = prompt
	}

	seconds := strings.TrimSpace(form.FormValue("seconds"))
	if seconds == "" {
		seconds = strings.TrimSpace(form.FormValue("duration"))
	}
	if seconds != "" {
		value, err := strconv.ParseFloat(seconds, 64)
		if err != nil || value <= 0 || value != math.Trunc(value) {
			return nil, newSoraCompatError(http.StatusBadRequest, "seconds must be a positive whole number")
		}
		body["duration"] = int(value)
	}

	if size := strings.TrimSpace(form.FormValue("size")); size != "" {
		resolution, aspect, err := parseSoraVideoSize(size)
		if err != nil {
			return nil, err
		}
		if resolution != "" {
			body["resolution"] = resolution
		}
		if aspect != "" {
			body["aspect_ratio"] = aspect
		}
	}

	imageURLs, err := applySoraCompatMetadata(body, strings.TrimSpace(form.FormValue("metadata")))
	if err != nil {
		return nil, err
	}

	if fileURL, err := h.storeSoraInputReference(c); err != nil {
		return nil, err
	} else if fileURL != "" {
		imageURLs = append(imageURLs, fileURL)
	}
	if len(imageURLs) > 0 {
		body["image_urls"] = imageURLs
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, newSoraCompatError(http.StatusBadRequest, "request could not be encoded")
	}
	return encoded, nil
}

// soraCompatMetadataKeys are the native body fields a metadata JSON object may
// set or override. Everything else is ignored.
var soraCompatMetadataKeys = []string{
	"resolution", "aspect_ratio", "audio", "duration",
	"start_frame_url", "end_frame_url", "guidances", "prompt_enhance",
}

func applySoraCompatMetadata(body map[string]any, metadata string) ([]string, error) {
	if metadata == "" {
		return nil, nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(metadata), &fields); err != nil {
		return nil, newSoraCompatError(http.StatusBadRequest, "metadata must be a JSON object")
	}
	for _, key := range soraCompatMetadataKeys {
		if raw, ok := fields[key]; ok {
			body[key] = raw
		}
	}
	var imageURLs []string
	if raw, ok := fields["image_urls"]; ok {
		if err := json.Unmarshal(raw, &imageURLs); err != nil {
			return nil, newSoraCompatError(http.StatusBadRequest, "metadata.image_urls must be an array of URLs")
		}
	}
	return imageURLs, nil
}

// storeSoraInputReference stores an input_reference image upload and returns
// its media URL, or "" when the field is absent.
func (h *OpenAIGatewayHandler) storeSoraInputReference(c *gin.Context) (string, error) {
	if c.Request.MultipartForm == nil || len(c.Request.MultipartForm.File["input_reference"]) == 0 {
		return "", nil
	}
	if h.videoInputHandler == nil || h.videoInputHandler.store == nil {
		return "", newSoraCompatError(http.StatusServiceUnavailable, "Video input service is not configured")
	}
	file, header, err := c.Request.FormFile("input_reference")
	if err != nil {
		return "", newSoraCompatError(http.StatusBadRequest, "input_reference could not be opened")
	}
	defer func() { _ = file.Close() }()
	input, err := h.videoInputHandler.store.SaveMedia(file, service.VideoInputKindImage, header.Filename)
	if err != nil {
		if errors.Is(err, service.ErrVideoInputTooLarge) {
			return "", newSoraCompatError(http.StatusRequestEntityTooLarge, "input_reference is too large")
		}
		return "", newSoraCompatError(http.StatusBadRequest, "input_reference must be a PNG, JPEG, or WebP image")
	}
	return input.URL, nil
}

var soraResolutionLabelPattern = regexp.MustCompile(`^\d{3,4}p$`)

// parseSoraVideoSize maps a Sora-format size ("1920x1080", "720x1280", or a
// direct resolution label like "720p" / "4k") onto the native resolution and
// aspect-ratio fields. Model-specific support is enforced by validation later.
func parseSoraVideoSize(size string) (resolution, aspect string, err error) {
	normalized := strings.ToLower(strings.TrimSpace(size))
	switch normalized {
	case "":
		return "", "", nil
	case "4k", "uhd":
		return "2160p", "", nil
	}
	if soraResolutionLabelPattern.MatchString(normalized) {
		return normalized, "", nil
	}
	separator := ""
	for _, candidate := range []string{"x", "*", "×"} {
		if strings.Contains(normalized, candidate) {
			separator = candidate
			break
		}
	}
	if separator == "" {
		return "", "", newSoraCompatError(http.StatusBadRequest, "size must look like 1280x720 or a resolution label like 720p")
	}
	parts := strings.SplitN(normalized, separator, 2)
	width, wErr := strconv.Atoi(strings.TrimSpace(parts[0]))
	height, hErr := strconv.Atoi(strings.TrimSpace(parts[1]))
	if wErr != nil || hErr != nil || width <= 0 || height <= 0 {
		return "", "", newSoraCompatError(http.StatusBadRequest, "size must look like 1280x720 or a resolution label like 720p")
	}
	return nearestSoraResolution(width, height), nearestSoraAspectRatio(width, height), nil
}

var soraResolutionShortSides = []int{400, 480, 544, 720, 960, 1080, 1440, 2160}

func nearestSoraResolution(width, height int) string {
	short := width
	if height < short {
		short = height
	}
	best := soraResolutionShortSides[0]
	for _, candidate := range soraResolutionShortSides[1:] {
		if abs(short-candidate) < abs(short-best) {
			best = candidate
		}
	}
	return strconv.Itoa(best) + "p"
}

var soraAspectRatios = [][2]int{{16, 9}, {9, 16}, {1, 1}, {4, 3}, {3, 4}, {21, 9}, {9, 21}}

func nearestSoraAspectRatio(width, height int) string {
	ratio := math.Log(float64(width) / float64(height))
	best := soraAspectRatios[0]
	bestDistance := math.Abs(ratio - math.Log(float64(best[0])/float64(best[1])))
	for _, candidate := range soraAspectRatios[1:] {
		distance := math.Abs(ratio - math.Log(float64(candidate[0])/float64(candidate[1])))
		if distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	return strconv.Itoa(best[0]) + ":" + strconv.Itoa(best[1])
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// soraVideoStatus maps native job statuses onto OpenAI video statuses. Every
// native status must map to a recognized value; unknown values break relays.
func soraVideoStatus(status string) string {
	switch status {
	case service.VideoJobPending:
		return "queued"
	case service.VideoJobRunning, service.VideoJobSettling:
		return "in_progress"
	case service.VideoJobCompleted:
		return "completed"
	case service.VideoJobFailed, service.VideoJobCanceled:
		return "failed"
	default:
		return "queued"
	}
}

func soraVideoObject(job *service.VideoJob) gin.H {
	progress := 0
	if job.Status == service.VideoJobCompleted {
		progress = 100
	}
	out := gin.H{
		"id":         job.JobID,
		"object":     "video",
		"model":      job.RequestedModel,
		"status":     soraVideoStatus(job.Status),
		"progress":   progress,
		"created_at": job.CreatedAt.Unix(),
		"seconds":    strconv.Itoa(job.DurationSeconds),
		"size":       soraVideoSize(job.Resolution, job.AspectRatio),
	}
	if job.FinishedAt != nil {
		out["completed_at"] = job.FinishedAt.Unix()
	}
	switch job.Status {
	case service.VideoJobFailed:
		message := strings.TrimSpace(service.SanitizeVideoProviderMessage(job.ErrorMessage))
		if message == "" {
			message = "Video generation failed"
		}
		out["error"] = gin.H{"code": "generation_failed", "message": message}
	case service.VideoJobCanceled:
		out["error"] = gin.H{"code": "canceled", "message": "Video job was canceled"}
	}
	return out
}

// soraVideoSize rebuilds a WxH size string from the stored resolution label
// and aspect ratio, falling back to the plain label when the ratio is unknown.
func soraVideoSize(resolution, aspectRatio string) string {
	label := strings.ToLower(strings.TrimSpace(resolution))
	short, err := strconv.Atoi(strings.TrimSuffix(label, "p"))
	if err != nil || short <= 0 {
		return label
	}
	parts := strings.SplitN(strings.TrimSpace(aspectRatio), ":", 2)
	if len(parts) != 2 {
		return label
	}
	ratioW, wErr := strconv.Atoi(strings.TrimSpace(parts[0]))
	ratioH, hErr := strconv.Atoi(strings.TrimSpace(parts[1]))
	if wErr != nil || hErr != nil || ratioW <= 0 || ratioH <= 0 {
		return label
	}
	if ratioW >= ratioH {
		width := int(math.Round(float64(short) * float64(ratioW) / float64(ratioH)))
		return strconv.Itoa(width) + "x" + strconv.Itoa(short)
	}
	height := int(math.Round(float64(short) * float64(ratioH) / float64(ratioW)))
	return strconv.Itoa(short) + "x" + strconv.Itoa(height)
}
