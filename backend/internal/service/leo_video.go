package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type LeoVideoRequestInfo struct {
	Model           string
	Prompt          string
	Resolution      string
	DurationSeconds int
	ImageURL        string
	ImageURLs       []string
	AspectRatio     string
	Audio           bool
}

const PublicVideoPlatform = "video"

var publicVideoProviderNamePattern = regexp.MustCompile(`(?i)\b(?:leonardo(?:\.ai| ai)?|leo\s*studio|leo)\b`)

func PublicPlatformID(platform string) string {
	if IsMediaPlatform(platform) {
		return PublicVideoPlatform
	}
	return platform
}

func PublicVideoErrorMessage(message string) string {
	return strings.TrimSpace(publicVideoProviderNamePattern.ReplaceAllString(message, "video service"))
}

// PublicLeoVideoResult removes provider-owned metadata before a result crosses
// the customer-facing API boundary. Billing and routing continue to consume
// the original upstream response internally.
func PublicLeoVideoResult(result json.RawMessage) json.RawMessage {
	var payload any
	if err := json.Unmarshal(result, &payload); err != nil {
		return nil
	}
	publicResult, err := json.Marshal(stripPrivateLeoVideoValue(payload))
	if err != nil {
		return nil
	}
	return publicResult
}

func stripPrivateLeoVideoValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		public := make(map[string]any, len(typed))
		for key, nested := range typed {
			if privateLeoVideoResultKey(key) {
				continue
			}
			public[key] = stripPrivateLeoVideoValue(nested)
		}
		return public
	case []any:
		public := make([]any, len(typed))
		for index, nested := range typed {
			public[index] = stripPrivateLeoVideoValue(nested)
		}
		return public
	default:
		return value
	}
}

func privateLeoVideoResultKey(key string) bool {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "provider", "uuid", "source_url", "video_url", "generation_id", "upstream_job_id", "account_id", "api_key", "cookie":
		return true
	default:
		return false
	}
}

func ParseLeoVideoRequest(body []byte) (LeoVideoRequestInfo, error) {
	if !gjson.ValidBytes(body) {
		return LeoVideoRequestInfo{}, fmt.Errorf("invalid video JSON request")
	}
	imageURLs := LeoVideoImageURLs(body)
	info := LeoVideoRequestInfo{
		Model:       strings.TrimSpace(gjson.GetBytes(body, "model").String()),
		Prompt:      strings.TrimSpace(gjson.GetBytes(body, "prompt").String()),
		Resolution:  strings.TrimSpace(gjson.GetBytes(body, "resolution").String()),
		ImageURLs:   imageURLs,
		AspectRatio: strings.TrimSpace(gjson.GetBytes(body, "aspect_ratio").String()),
		Audio:       gjson.GetBytes(body, "audio").Bool(),
	}
	if len(imageURLs) > 0 {
		info.ImageURL = imageURLs[0]
	}
	if duration := gjson.GetBytes(body, "duration"); duration.Exists() && duration.Type == gjson.Number {
		info.DurationSeconds = int(duration.Int())
	}
	return info, nil
}

func LeoVideoReferenceURLs(body []byte) []string {
	urls := LeoVideoImageURLs(body)
	seen := make(map[string]struct{}, len(urls)+2)
	for _, raw := range urls {
		seen[raw] = struct{}{}
	}
	add := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}
		if _, ok := seen[raw]; ok {
			return
		}
		seen[raw] = struct{}{}
		urls = append(urls, raw)
	}
	for _, path := range []string{
		"guidances.video_reference_base.#.video.url",
		"guidances.audio_reference.#.audio.url",
	} {
		for _, item := range gjson.GetBytes(body, path).Array() {
			add(item.String())
		}
	}
	return urls
}

func LeoVideoImageURLs(body []byte) []string {
	seen := make(map[string]struct{})
	urls := make([]string, 0, 6)
	add := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}
		if _, ok := seen[raw]; ok {
			return
		}
		seen[raw] = struct{}{}
		urls = append(urls, raw)
	}
	for _, path := range []string{"image_url", "start_frame_url", "end_frame_url"} {
		add(gjson.GetBytes(body, path).String())
	}
	for _, path := range []string{
		"image_urls",
		"guidances.start_frame.#.image.url",
		"guidances.end_frame.#.image.url",
		"guidances.image_reference.#.image.url",
	} {
		for _, item := range gjson.GetBytes(body, path).Array() {
			add(item.String())
		}
	}
	return urls
}

func (s *OpenAIGatewayService) ForwardLeoVideo(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	if account == nil || !account.IsMediaAPIAccount() {
		return nil, fmt.Errorf("leo account is required")
	}
	requestInfo, err := ParseLeoVideoRequest(body)
	if err != nil {
		return nil, err
	}
	upstreamModel, _ := account.ResolveMappedModel(requestInfo.Model)
	requestInfo, err = ValidateLeoVideoRequestForModel(body, upstreamModel)
	if err != nil {
		return nil, err
	}
	body, err = sjson.SetBytes(body, "model", upstreamModel)
	if err != nil {
		return nil, fmt.Errorf("rewrite leo video model: %w", err)
	}
	targetURL, err := BuildLeoVideosGenerationsURL(account.GetLeoBaseURL())
	if err != nil {
		return nil, err
	}
	apiKey := account.GetLeoAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("leo API key is not configured")
	}

	upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
	defer releaseUpstreamCtx()
	upstreamReq, err := http.NewRequestWithContext(upstreamCtx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	upstreamReq.Header.Set("Authorization", "Bearer "+apiKey)
	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Content-Type", "application/json")

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}
	upstreamStart := time.Now()
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, nil)
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return s.handleLeoVideoErrorResponse(ctx, c, account, resp, responseBody)
	}
	if _, _, _, err := parseVideoOutputResult(responseBody); err != nil {
		return nil, fmt.Errorf("video service returned no usable video output")
	}

	writeOpenAIPassthroughResponseHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/json"
	}
	publicResponseBody := PublicLeoVideoResult(responseBody)
	if len(publicResponseBody) == 0 {
		publicResponseBody = responseBody
	}
	c.Data(resp.StatusCode, contentType, publicResponseBody)

	resolution := strings.TrimSpace(gjson.GetBytes(responseBody, "provider.resolution").String())
	if resolution == "" {
		resolution = requestInfo.Resolution
	}
	durationSeconds := 0
	if duration := gjson.GetBytes(responseBody, "provider.duration"); duration.Exists() && duration.Type == gjson.Number {
		durationSeconds = int(duration.Int())
	}
	if durationSeconds <= 0 {
		durationSeconds = requestInfo.DurationSeconds
	}

	return &OpenAIForwardResult{
		RequestID:            strings.TrimSpace(resp.Header.Get("x-request-id")),
		ResponseID:           strings.TrimSpace(gjson.GetBytes(responseBody, "provider.generation_id").String()),
		Model:                requestInfo.Model,
		BillingModel:         requestInfo.Model,
		UpstreamModel:        upstreamModel,
		UpstreamEndpoint:     "/v1/videos/generations",
		ResponseHeaders:      resp.Header.Clone(),
		Duration:             time.Since(startTime),
		VideoCount:           1,
		VideoResolution:      NormalizeLeoVideoBillingResolutionOrDefault(upstreamModel, resolution),
		VideoDurationSeconds: NormalizeLeoVideoBillingDurationSecondsOrDefault(upstreamModel, durationSeconds),
	}, nil
}

func (s *OpenAIGatewayService) handleLeoVideoErrorResponse(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	resp *http.Response,
	body []byte,
) (*OpenAIForwardResult, error) {
	if isLeoVideoFailoverStatus(resp.StatusCode) {
		if s.rateLimitService != nil {
			s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)
		}
		return nil, &UpstreamFailoverError{
			StatusCode:      resp.StatusCode,
			ResponseBody:    body,
			ResponseHeaders: resp.Header.Clone(),
		}
	}

	message := PublicVideoErrorMessage(sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body))))
	if message == "" {
		message = fmt.Sprintf("Video service rejected the request with HTTP %d", resp.StatusCode)
	}
	writeOpenAIPassthroughResponseHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	c.JSON(resp.StatusCode, gin.H{
		"error": gin.H{
			"type":    "invalid_request_error",
			"message": message,
		},
	})
	return nil, nil
}

func isLeoVideoFailoverStatus(statusCode int) bool {
	return statusCode == http.StatusUnauthorized ||
		statusCode == http.StatusForbidden ||
		statusCode == http.StatusTooManyRequests ||
		statusCode >= http.StatusInternalServerError
}
