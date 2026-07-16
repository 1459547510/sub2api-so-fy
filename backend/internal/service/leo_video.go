package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
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
}

func ParseLeoVideoRequest(body []byte) (LeoVideoRequestInfo, error) {
	if !gjson.ValidBytes(body) {
		return LeoVideoRequestInfo{}, fmt.Errorf("invalid leo video JSON request")
	}
	info := LeoVideoRequestInfo{
		Model:      strings.TrimSpace(gjson.GetBytes(body, "model").String()),
		Prompt:     strings.TrimSpace(gjson.GetBytes(body, "prompt").String()),
		Resolution: strings.TrimSpace(gjson.GetBytes(body, "resolution").String()),
		ImageURL:   strings.TrimSpace(gjson.GetBytes(body, "image_url").String()),
	}
	if duration := gjson.GetBytes(body, "duration"); duration.Exists() && duration.Type == gjson.Number {
		info.DurationSeconds = int(duration.Int())
	}
	return info, nil
}

func (s *OpenAIGatewayService) ForwardLeoVideo(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
) (*OpenAIForwardResult, error) {
	startTime := time.Now()
	if account == nil || !account.IsLeo() {
		return nil, fmt.Errorf("leo account is required")
	}
	requestInfo, err := ParseLeoVideoRequest(body)
	if err != nil {
		return nil, err
	}
	upstreamModel, _ := account.ResolveMappedModel(requestInfo.Model)
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
	defer resp.Body.Close()

	responseBody, err := ReadUpstreamResponseBody(resp.Body, s.cfg, c, nil)
	if err != nil {
		return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return s.handleLeoVideoErrorResponse(ctx, c, account, resp, responseBody)
	}

	writeOpenAIPassthroughResponseHeaders(c.Writer.Header(), resp.Header, s.responseHeaderFilter)
	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/json"
	}
	c.Data(resp.StatusCode, contentType, responseBody)

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
		VideoResolution:      NormalizeVideoBillingResolutionOrDefault(resolution),
		VideoDurationSeconds: NormalizeVideoBillingDurationSecondsOrDefault(durationSeconds),
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

	message := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body)))
	if message == "" {
		message = fmt.Sprintf("LeoStudio rejected the request with HTTP %d", resp.StatusCode)
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
