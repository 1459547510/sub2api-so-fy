package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/sjson"
)

type LeoAsyncAccepted struct {
	JobID     int64  `json:"job_id"`
	Status    string `json:"status"`
	StatusURL string `json:"status_url"`
}

type LeoAsyncJobError struct {
	Message string `json:"message"`
}

type LeoAsyncJob struct {
	JobID     int64             `json:"job_id"`
	Status    string            `json:"status"`
	StatusURL string            `json:"status_url"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt int64             `json:"updated_at"`
	Result    json.RawMessage   `json:"result,omitempty"`
	Error     *LeoAsyncJobError `json:"error,omitempty"`
}

type LeoAsyncUpstreamError struct {
	StatusCode int
	Message    string
	Retryable  bool
	Ambiguous  bool
}

var videoProviderNamePattern = regexp.MustCompile(`(?i)\b(?:leonardo(?:\.ai| ai)?|leo\s*studio|leo)\b`)

func SanitizeVideoProviderMessage(message string) string {
	return videoProviderNamePattern.ReplaceAllString(message, "Video service")
}

func (e *LeoAsyncUpstreamError) Error() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return "Video service request failed"
	}
	return SanitizeVideoProviderMessage(e.Message)
}

func PrefersLeoRespondAsync(header http.Header) bool {
	for _, value := range header.Values("Prefer") {
		for _, preference := range strings.Split(value, ",") {
			if strings.EqualFold(strings.TrimSpace(preference), "respond-async") {
				return true
			}
		}
	}
	return false
}

func (s *OpenAIGatewayService) CreateLeoAsyncVideo(ctx context.Context, account *Account, body []byte) (*LeoAsyncAccepted, error) {
	if account == nil || !account.IsLeo() {
		return nil, fmt.Errorf("leo account is required")
	}
	request, err := ParseLeoVideoRequest(body)
	if err != nil {
		return nil, err
	}
	upstreamModel, _ := account.ResolveMappedModel(request.Model)
	if _, err := ValidateLeoVideoRequestForModel(body, upstreamModel); err != nil {
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

	responseBody, err := s.doLeoAsyncRequest(ctx, account, http.MethodPost, targetURL, body, true, false)
	if err != nil {
		return nil, err
	}
	var accepted LeoAsyncAccepted
	if err := json.Unmarshal(responseBody, &accepted); err != nil || accepted.JobID <= 0 || strings.TrimSpace(accepted.Status) == "" {
		return nil, &LeoAsyncUpstreamError{StatusCode: http.StatusAccepted, Message: "Video service returned an invalid job response", Ambiguous: true}
	}
	return &accepted, nil
}

func (s *OpenAIGatewayService) GetLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error) {
	return s.readLeoAsyncVideo(ctx, account, http.MethodGet, upstreamJobID)
}

func (s *OpenAIGatewayService) CancelLeoAsyncVideo(ctx context.Context, account *Account, upstreamJobID int64) (*LeoAsyncJob, error) {
	return s.readLeoAsyncVideo(ctx, account, http.MethodDelete, upstreamJobID)
}

func (s *OpenAIGatewayService) readLeoAsyncVideo(ctx context.Context, account *Account, method string, upstreamJobID int64) (*LeoAsyncJob, error) {
	if account == nil || !account.IsLeo() {
		return nil, fmt.Errorf("leo account is required")
	}
	if upstreamJobID <= 0 {
		return nil, fmt.Errorf("leo video job ID must be positive")
	}
	baseURL := account.GetLeoBaseURL()
	if baseURL == "" {
		return nil, fmt.Errorf("leo base URL is not configured")
	}
	targetURL := baseURL + "/videos/jobs/" + strconv.FormatInt(upstreamJobID, 10)
	responseBody, err := s.doLeoAsyncRequest(ctx, account, method, targetURL, nil, false, true)
	if err != nil {
		return nil, err
	}
	var job LeoAsyncJob
	if err := json.Unmarshal(responseBody, &job); err != nil || job.JobID <= 0 || strings.TrimSpace(job.Status) == "" {
		return nil, &LeoAsyncUpstreamError{StatusCode: http.StatusOK, Message: "Video service returned an invalid job response", Retryable: true}
	}
	return &job, nil
}

func (s *OpenAIGatewayService) doLeoAsyncRequest(
	ctx context.Context,
	account *Account,
	method string,
	targetURL string,
	body []byte,
	preferAsync bool,
	retryTransport bool,
) ([]byte, error) {
	apiKey := account.GetLeoAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("leo API key is not configured")
	}
	request, err := http.NewRequestWithContext(ctx, method, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Accept", "application/json")
	if len(body) != 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	if preferAsync {
		request.Header.Set("Prefer", "respond-async")
	}

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}
	response, err := s.httpUpstream.Do(request, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return nil, &LeoAsyncUpstreamError{Message: "Video service request failed", Retryable: retryTransport, Ambiguous: !retryTransport}
	}
	if response == nil || response.Body == nil {
		return nil, &LeoAsyncUpstreamError{Message: "Video service returned an empty response", Retryable: retryTransport, Ambiguous: !retryTransport}
	}
	defer func() { _ = response.Body.Close() }()
	responseBody, err := ReadUpstreamResponseBody(response.Body, s.cfg, nil, nil)
	if err != nil {
		return nil, &LeoAsyncUpstreamError{StatusCode: response.StatusCode, Message: "Video service response could not be read", Retryable: retryTransport, Ambiguous: !retryTransport}
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		message := sanitizeLeoAsyncErrorMessage(responseBody, apiKey, response.StatusCode)
		return nil, &LeoAsyncUpstreamError{
			StatusCode: response.StatusCode,
			Message:    message,
			Retryable:  isLeoVideoFailoverStatus(response.StatusCode),
		}
	}
	return responseBody, nil
}

func sanitizeLeoAsyncErrorMessage(body []byte, apiKey string, statusCode int) string {
	message := sanitizeUpstreamErrorMessage(strings.TrimSpace(extractUpstreamErrorMessage(body)))
	if apiKey != "" {
		message = strings.ReplaceAll(message, apiKey, "***")
	}
	if message == "" {
		return fmt.Sprintf("Video service request failed with HTTP %d", statusCode)
	}
	return SanitizeVideoProviderMessage(message)
}
