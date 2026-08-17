package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type liveHTTPUpstream struct {
	lastURL  string
	lastBody []byte
}

func (u *liveHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	if req != nil && req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
		u.lastURL = req.URL.String()
		u.lastBody = body
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	client := &http.Client{Timeout: 3 * time.Minute}
	return client.Do(req)
}

func (u *liveHTTPUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, concurrency)
}

func TestLiveLeoImageEditsRewriteAgainstLeoStudio(t *testing.T) {
	key := strings.TrimSpace(os.Getenv("LEO_LIVE_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("LEO_LIVE_BASE_URL"))
	model := strings.TrimSpace(os.Getenv("LEO_LIVE_IMAGE_MODEL"))
	refURL := strings.TrimSpace(os.Getenv("LEO_LIVE_REF_URL"))
	if key == "" || baseURL == "" || model == "" {
		t.Skip("set LEO_LIVE_API_KEY, LEO_LIVE_BASE_URL, and LEO_LIVE_IMAGE_MODEL to run the live Leo image test")
	}
	if refURL == "" {
		refURL = "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Fronalpstock_big.jpg/640px-Fronalpstock_big.jpg"
	}

	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"` + model + `","prompt":"keep the same scene, photorealistic, no text","n":1,"aspect_ratio":"1:1","images":[{"image_url":"` + refURL + `"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(WithOpenAIImagesPlatform(req.Context(), PlatformLeo))
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	upstream := &liveHTTPUpstream{}
	svc := &OpenAIGatewayService{
		cfg: &config.Config{
			Security: config.SecurityConfig{
				URLAllowlist: config.URLAllowlistConfig{AllowInsecureHTTP: true},
			},
		},
		httpUpstream: upstream,
	}
	parsed, err := svc.ParseOpenAIImagesRequest(c, body)
	require.NoError(t, err)

	account := &Account{
		ID:       9001,
		Name:     "live-leo-image",
		Platform: PlatformLeo,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  key,
			"base_url": baseURL,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	result, err := svc.ForwardImages(ctx, c, account, body, parsed, "")
	if err != nil {
		t.Fatalf("live Leo edits rewrite failed: %v recorder=%s upstream_url=%s upstream_body=%s", err, rec.Body.String(), upstream.lastURL, string(upstream.lastBody))
	}
	require.NotNil(t, result)
	require.Equal(t, baseURL+"/images/generations", upstream.lastURL)
	require.Equal(t, refURL, gjson.GetBytes(upstream.lastBody, "image_urls.0").String())
	require.False(t, gjson.GetBytes(upstream.lastBody, "images").Exists())
	require.False(t, gjson.GetBytes(upstream.lastBody, "image_url").Exists())
	require.Greater(t, result.ImageCount, 0)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotEmpty(t, gjson.Get(rec.Body.String(), "data.0.url").String()+gjson.Get(rec.Body.String(), "data.0.b64_json").String())
}
