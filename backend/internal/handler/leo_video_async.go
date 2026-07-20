package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkghttputil "github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type leoVideoJobResponse struct {
	JobID          string          `json:"job_id"`
	Status         string          `json:"status"`
	StatusURL      string          `json:"status_url"`
	RequestedModel string          `json:"model,omitempty"`
	Prompt         string          `json:"prompt,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Result         json.RawMessage `json:"result,omitempty"`
	Error          *leoVideoJobErr `json:"error,omitempty"`
}

type leoVideoJobErr struct {
	Message string `json:"message"`
}

func (h *OpenAIGatewayHandler) LeoVideoAsyncGeneration(c *gin.Context) {
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
	body, err := pkghttputil.ReadRequestBodyWithPrealloc(c.Request)
	if err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
	}
	if len(body) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
	}
	if contentType := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type"))); contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Content-Type must be application/json")
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
		status, code := http.StatusBadGateway, "upstream_error"
		var upstreamErr *service.LeoAsyncUpstreamError
		if errors.Is(err, service.ErrVideoInsufficientBalance) {
			status, code = http.StatusPaymentRequired, "billing_error"
		} else if errors.As(err, &upstreamErr) && (upstreamErr.StatusCode == http.StatusBadRequest || upstreamErr.StatusCode == http.StatusUnprocessableEntity) {
			status, code = upstreamErr.StatusCode, "invalid_request_error"
		} else if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "valid JSON") {
			status, code = http.StatusBadRequest, "invalid_request_error"
		}
		h.errorResponse(c, status, code, err.Error())
		return
	}
	statusURL := "/v1/videos/jobs/" + job.JobID
	c.Header("Preference-Applied", "respond-async")
	c.Header("Location", statusURL)
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusAccepted, gin.H{"job_id": job.JobID, "status": job.Status, "status_url": statusURL})
}

func (h *OpenAIGatewayHandler) LeoVideoJobs(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	jobs, err := h.videoJobService.List(c.Request.Context(), apiKey.ID, limit, strings.TrimSpace(c.Query("status")))
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video job list unavailable")
		return
	}
	c.Header("Cache-Control", "no-store")
	responses := make([]leoVideoJobResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, publicLeoVideoJob(job))
	}
	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *OpenAIGatewayHandler) LeoVideoJob(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	job, err := h.videoJobService.Get(c.Request.Context(), c.Param("job_id"), apiKey.ID)
	if errors.Is(err, service.ErrVideoJobNotFound) {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Video job not found")
		return
	}
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video job unavailable")
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, publicLeoVideoJob(job))
}

func (h *OpenAIGatewayHandler) LeoVideoJobContent(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	if h.videoJobService == nil || h.videoOutputStore == nil {
		h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "Video service is not configured")
		return
	}
	job, err := h.videoJobService.Get(c.Request.Context(), c.Param("job_id"), apiKey.ID)
	if errors.Is(err, service.ErrVideoJobNotFound) {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Video job not found")
		return
	}
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video job unavailable")
		return
	}
	if job.Status != service.VideoJobCompleted {
		h.errorResponse(c, http.StatusConflict, "conflict_error", "Video output is not available")
		return
	}
	file, err := h.videoOutputStore.Open(job.JobID)
	if errors.Is(err, service.ErrVideoOutputNotFound) {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Video output not found")
		return
	}
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video output unavailable")
		return
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "Video output unavailable")
		return
	}
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", `inline; filename="`+job.JobID+`.mp4"`)
	c.Header("Cache-Control", "private, no-store")
	http.ServeContent(c.Writer, c.Request, job.JobID+".mp4", info.ModTime(), file)
}

func (h *OpenAIGatewayHandler) CancelLeoVideoJob(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}
	job, err := h.videoJobService.Cancel(c.Request.Context(), c.Param("job_id"), apiKey.ID)
	if errors.Is(err, service.ErrVideoJobNotFound) {
		h.errorResponse(c, http.StatusNotFound, "not_found_error", "Video job not found")
		return
	}
	if errors.Is(err, service.ErrVideoJobCancelConflict) {
		h.errorResponse(c, http.StatusConflict, "conflict_error", "Video job cannot be canceled")
		return
	}
	if err != nil {
		h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Video job cancellation failed")
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, publicLeoVideoJob(job))
}

func (h *OpenAIGatewayHandler) LeoVideoUpload(c *gin.Context) {
	if h.videoInputHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{"message": "Video input service is not configured"}})
		return
	}
	h.videoInputHandler.Upload(c)
}

func (h *OpenAIGatewayHandler) LeoVideoInputInternal(c *gin.Context) {
	if h.videoInputHandler == nil {
		c.Status(http.StatusNotFound)
		return
	}
	h.videoInputHandler.GetInternal(c)
}

func publicLeoVideoJob(job *service.VideoJob) leoVideoJobResponse {
	response := leoVideoJobResponse{JobID: job.JobID, Status: job.Status, StatusURL: "/v1/videos/jobs/" + job.JobID,
		RequestedModel: job.RequestedModel, Prompt: job.Prompt, CreatedAt: job.CreatedAt, UpdatedAt: job.UpdatedAt}
	if job.Status == service.VideoJobCompleted && len(job.Result) != 0 {
		response.Result = append(json.RawMessage(nil), job.Result...)
	}
	if job.Status == service.VideoJobFailed && strings.TrimSpace(job.ErrorMessage) != "" {
		response.Error = &leoVideoJobErr{Message: job.ErrorMessage}
	}
	return response
}
