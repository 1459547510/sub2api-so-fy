package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type leoVideoHTTPUpstream struct {
	response    *http.Response
	err         error
	request     *http.Request
	requestBody []byte
	proxyURL    string
	accountID   int64
	concurrency int
}

func (u *leoVideoHTTPUpstream) Do(req *http.Request, proxyURL string, accountID int64, concurrency int) (*http.Response, error) {
	u.request = req
	u.proxyURL = proxyURL
	u.accountID = accountID
	u.concurrency = concurrency
	if req != nil && req.Body != nil {
		u.requestBody, _ = io.ReadAll(req.Body)
	}
	return u.response, u.err
}

func (u *leoVideoHTTPUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, concurrency)
}

func TestForwardLeoVideoMapsModelAndAddsBearer(t *testing.T) {
	body := []byte(`{"model":"seedance","prompt":"city","resolution":"720p","duration":8,"audio":false,"image_urls":["https://cdn.example/ref-1.png","https://cdn.example/ref-2.png"],"guidances":{"video_reference_base":[{"video":{"id":"22222222-2222-2222-2222-222222222222","type":"GENERATED"}}]}}`)
	responseBody := `{"data":[{"url":"https://cdn.example/video.mp4","mp4_url":"https://cdn.example/video.mp4"}],"provider":{"generation_id":"gen-1","resolution":"RESOLUTION_720","duration":12}}`
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, responseBody)}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	account := newLeoVideoTestAccount()
	rec, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, account, body)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, http.MethodPost, upstream.request.Method)
	require.Equal(t, "http://leo.internal:8000/v1/videos/generations", upstream.request.URL.String())
	require.Equal(t, "Bearer leo-secret", upstream.request.Header.Get("Authorization"))
	require.Equal(t, "application/json", upstream.request.Header.Get("Accept"))
	require.Equal(t, "application/json", upstream.request.Header.Get("Content-Type"))
	require.Equal(t, "seedance-2.0", gjson.GetBytes(upstream.requestBody, "model").String())
	require.Equal(t, "city", gjson.GetBytes(upstream.requestBody, "prompt").String())
	require.Equal(t, "720p", gjson.GetBytes(upstream.requestBody, "resolution").String())
	require.Equal(t, int64(8), gjson.GetBytes(upstream.requestBody, "duration").Int())
	require.False(t, gjson.GetBytes(upstream.requestBody, "audio").Bool())
	require.Len(t, gjson.GetBytes(upstream.requestBody, "image_urls").Array(), 2)
	require.Equal(t, "22222222-2222-2222-2222-222222222222", gjson.GetBytes(upstream.requestBody, "guidances.video_reference_base.0.video.id").String())
	require.Equal(t, 1, result.VideoCount)
	require.Equal(t, "720p", result.VideoResolution)
	require.Equal(t, 12, result.VideoDurationSeconds)
	require.Equal(t, "gen-1", result.ResponseID)
	require.Equal(t, "seedance", result.Model)
	require.Equal(t, "seedance-2.0", result.UpstreamModel)
	require.JSONEq(t, responseBody, rec.Body.String())
}

func TestForwardLeoVideoPreservesLTXFastTwentySecondDuration(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, `{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_2160","duration":20}}`)}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	account := newLeoVideoTestAccount()
	account.Credentials["model_mapping"] = map[string]any{"ltx-2.3-fast": "ltxv-2.3-fast"}
	_, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, account, []byte(`{"model":"ltx-2.3-fast","prompt":"city","resolution":"2160p","duration":20}`))

	require.NoError(t, err)
	require.Equal(t, "ltxv-2.3-fast", gjson.GetBytes(upstream.requestBody, "model").String())
	require.Equal(t, "ltx-2.3-fast", result.Model)
	require.Equal(t, VideoBillingResolution2160P, result.VideoResolution)
	require.Equal(t, 20, result.VideoDurationSeconds)
}

func TestForwardLeoVideoPassesMediaAndAudioReferenceURLs(t *testing.T) {
	body := []byte(`{"model":"seedance","prompt":"city","guidances":{"image_reference":[{"image":{"url":"https://cdn.example/ref.png"}}],"video_reference_base":[{"video":{"url":"https://cdn.example/ref.mp4","type":"UPLOADED"}}],"audio_reference":[{"audio":{"url":"https://cdn.example/ref.mp3","type":"UPLOADED"}}]}}`)
	responseBody := `{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_720","duration":8}}`
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, responseBody)}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	_, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), body)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "https://cdn.example/ref.mp4", gjson.GetBytes(upstream.requestBody, "guidances.video_reference_base.0.video.url").String())
	require.Equal(t, "https://cdn.example/ref.mp3", gjson.GetBytes(upstream.requestBody, "guidances.audio_reference.0.audio.url").String())
}

func TestParseLeoVideoRequestCollectsNewGuidanceImageURLs(t *testing.T) {
	body := []byte(`{
		"model":"seedance-2.0","prompt":"city",
		"start_frame_url":"https://cdn.example/start.png",
		"end_frame_url":"https://cdn.example/end.png",
		"image_urls":["https://cdn.example/ref-1.png","https://cdn.example/ref-2.png"],
		"guidances":{"image_reference":[{"image":{"url":"https://cdn.example/ref-2.png"}},{"image":{"url":"https://cdn.example/ref-3.png"}}]}
	}`)

	info, err := ParseLeoVideoRequest(body)

	require.NoError(t, err)
	require.Equal(t, "https://cdn.example/start.png", info.ImageURL)
	require.Equal(t, []string{
		"https://cdn.example/start.png",
		"https://cdn.example/end.png",
		"https://cdn.example/ref-1.png",
		"https://cdn.example/ref-2.png",
		"https://cdn.example/ref-3.png",
	}, info.ImageURLs)
}

func TestParseLeoVideoRequestCollectsMediaAndAudioURLs(t *testing.T) {
	body := []byte(`{
		"model":"seedance-2.0","prompt":"city",
		"guidances":{
			"video_reference_base":[{"video":{"url":"https://cdn.example/ref.mp4"}}],
			"audio_reference":[{"audio":{"url":"https://cdn.example/ref.mp3"}}]
		}
	}`)

	info, err := ParseLeoVideoRequest(body)
	require.NoError(t, err)
	require.Empty(t, info.ImageURLs)
	require.Equal(t, []string{"https://cdn.example/ref.mp4", "https://cdn.example/ref.mp3"}, LeoVideoReferenceURLs(body))
}

func TestForwardLeoVideoFallsBackToRequestMetadata(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, `{"data":[{"url":"https://cdn.example/video.mp4"}]}`)}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	rec, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city","resolution":"1080p","duration":7}`))

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "1080p", result.VideoResolution)
	require.Equal(t, 7, result.VideoDurationSeconds)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestForwardLeoVideoRejectsSuccessfulResponseWithoutVideo(t *testing.T) {
	for _, responseBody := range []string{`{"data":[]}`, `{"data":[{"status":"completed"}]}`} {
		upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(http.StatusOK, responseBody)}
		svc := &OpenAIGatewayService{httpUpstream: upstream}
		recorder, c := newLeoVideoTestContext()

		result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city"}`))

		require.Error(t, err)
		require.Nil(t, result)
		require.False(t, c.Writer.Written())
		require.Empty(t, recorder.Body.String())
	}
}

func TestForwardLeoVideoRejectsInvalidJSONBeforeTransport(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	_, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":`))

	require.Error(t, err)
	require.Nil(t, result)
	require.Nil(t, upstream.request)
}

func TestForwardLeoVideoWritesRequestErrorsWithoutFailover(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusUnprocessableEntity} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(status, `{"error":{"message":"Leonardo AI rejected the bad prompt"}}`)}
			svc := &OpenAIGatewayService{httpUpstream: upstream}
			rec, c := newLeoVideoTestContext()

			result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city"}`))

			require.NoError(t, err)
			require.Nil(t, result)
			require.Equal(t, status, rec.Code)
			require.Contains(t, rec.Body.String(), "bad prompt")
			require.NotContains(t, strings.ToLower(rec.Body.String()), "leonardo")
			require.True(t, c.Writer.Written())
		})
	}
}

func TestPublicVideoErrorMessageHidesProviderNames(t *testing.T) {
	for _, message := range []string{
		"Leonardo rejected the request",
		"LEONARDO AI policy failure",
		"leonardo.ai is unavailable",
		"LeoStudio request failed",
		"Leo Studio response was invalid",
	} {
		public := PublicVideoErrorMessage(message)
		require.NotContains(t, strings.ToLower(public), "leonardo")
		require.NotContains(t, strings.ToLower(public), "leostudio")
		require.NotContains(t, strings.ToLower(public), "leo studio")
		require.Contains(t, public, "video service")
	}
}

func TestForwardLeoVideoFailsOverBeforeResponseCommit(t *testing.T) {
	for _, status := range []int{
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			upstream := &leoVideoHTTPUpstream{response: leoVideoResponse(status, `{"error":{"message":"upstream unavailable"}}`)}
			svc := &OpenAIGatewayService{httpUpstream: upstream}
			rec, c := newLeoVideoTestContext()

			result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city"}`))

			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.True(t, errors.As(err, &failoverErr))
			require.Equal(t, status, failoverErr.StatusCode)
			require.False(t, c.Writer.Written())
			require.Empty(t, rec.Body.String())
		})
	}
}

func TestForwardLeoVideoTransportFailureFailsOver(t *testing.T) {
	upstream := &leoVideoHTTPUpstream{err: errors.New("dial failed")}
	svc := &OpenAIGatewayService{httpUpstream: upstream}
	_, c := newLeoVideoTestContext()

	result, err := svc.ForwardLeoVideo(context.Background(), c, newLeoVideoTestAccount(), []byte(`{"model":"seedance","prompt":"city"}`))

	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.True(t, errors.As(err, &failoverErr))
}

func newLeoVideoTestAccount() *Account {
	return &Account{
		ID:          1,
		Platform:    PlatformLeo,
		Type:        AccountTypeAPIKey,
		Concurrency: 2,
		Credentials: map[string]any{
			"base_url": "http://leo.internal:8000/v1",
			"api_key":  "leo-secret",
			"model_mapping": map[string]any{
				"seedance": "seedance-2.0",
			},
		},
	}
}

func newLeoVideoTestContext() (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(""))
	return rec, c
}

func leoVideoResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
