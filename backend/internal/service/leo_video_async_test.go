package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestPrefersLeoRespondAsync(t *testing.T) {
	header := http.Header{}
	header.Add("Prefer", "wait=5, Respond-Async")
	require.True(t, PrefersLeoRespondAsync(header))

	header.Set("Prefer", "respond-async=false")
	require.False(t, PrefersLeoRespondAsync(header))

	header.Set("Prefer", "respond-async; handling=lenient")
	require.False(t, PrefersLeoRespondAsync(header))
}

func TestCreateLeoAsyncVideoMapsModelAndAuthenticates(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusAccepted, `{"job_id":42,"status":"pending","status_url":"/v1/videos/jobs/42"}`)}
	service := &OpenAIGatewayService{httpUpstream: upstream}

	accepted, err := service.CreateLeoAsyncVideo(context.Background(), newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city"}`))

	require.NoError(t, err)
	require.Equal(t, int64(42), accepted.JobID)
	require.Equal(t, "pending", accepted.Status)
	require.Equal(t, http.MethodPost, upstream.request.Method)
	require.Equal(t, "http://leo.internal:8000/v1/videos/generations", upstream.request.URL.String())
	require.Equal(t, "Bearer leo-secret", upstream.request.Header.Get("Authorization"))
	require.Equal(t, "respond-async", upstream.request.Header.Get("Prefer"))
	require.Equal(t, "seedance-2.0", gjson.GetBytes(upstream.requestBody, "model").String())
}

func TestGetLeoAsyncVideoDecodesCompletedResult(t *testing.T) {
	response := `{"job_id":42,"status":"completed","status_url":"/v1/videos/jobs/42","result":{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_720","duration":12}}}`
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, response)}
	service := &OpenAIGatewayService{httpUpstream: upstream}

	job, err := service.GetLeoAsyncVideo(context.Background(), newLeoVideoTestAccount(), 42)

	require.NoError(t, err)
	require.Equal(t, int64(42), job.JobID)
	require.Equal(t, "completed", job.Status)
	require.JSONEq(t, `{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_720","duration":12}}`, string(job.Result))
	require.Equal(t, http.MethodGet, upstream.request.Method)
	require.Equal(t, "http://leo.internal:8000/v1/videos/jobs/42", upstream.request.URL.String())
	require.Empty(t, upstream.request.Header.Get("Prefer"))
	require.Equal(t, "Bearer leo-secret", upstream.request.Header.Get("Authorization"))
}

func TestCancelLeoAsyncVideoDecodesCanceledJob(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, `{"job_id":42,"status":"canceled","status_url":"/v1/videos/jobs/42"}`)}
	service := &OpenAIGatewayService{httpUpstream: upstream}

	job, err := service.CancelLeoAsyncVideo(context.Background(), newLeoVideoTestAccount(), 42)

	require.NoError(t, err)
	require.Equal(t, "canceled", job.Status)
	require.Equal(t, http.MethodDelete, upstream.request.Method)
	require.Equal(t, "http://leo.internal:8000/v1/videos/jobs/42", upstream.request.URL.String())
}

func TestLeoAsyncVideoErrorsAreTypedAndDoNotExposeSecret(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusUnauthorized, `{"error":{"message":"bad token leo-secret at https://leo.test/?access_token=raw-secret"}}`)}
	service := &OpenAIGatewayService{httpUpstream: upstream}

	_, err := service.GetLeoAsyncVideo(context.Background(), newLeoVideoTestAccount(), 42)

	var upstreamErr *LeoAsyncUpstreamError
	require.True(t, errors.As(err, &upstreamErr))
	require.Equal(t, http.StatusUnauthorized, upstreamErr.StatusCode)
	require.True(t, upstreamErr.Retryable)
	require.NotContains(t, err.Error(), "leo-secret")
	require.NotContains(t, err.Error(), "raw-secret")
	require.False(t, strings.Contains(err.Error(), "Authorization"))
}

func TestSanitizeVideoProviderMessageHidesUpstreamNames(t *testing.T) {
	message := SanitizeVideoProviderMessage("Leonardo request failed; Leonardo AI, Leonardo.ai, LeoStudio and Leo Studio rejected it")

	require.Equal(t, "Video provider request failed; Video provider, Video provider, Video provider and Video provider rejected it", message)
	for _, name := range []string{"Leonardo", "Leonardo AI", "Leonardo.ai", "LeoStudio", "Leo Studio"} {
		require.NotContains(t, message, name)
	}
}
