package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type leoIntegrationAccountRepo struct {
	service.AccountRepository
	account service.Account
}

func (r *leoIntegrationAccountRepo) GetByID(_ context.Context, id int64) (*service.Account, error) {
	if id != r.account.ID {
		return nil, service.ErrNoAvailableAccounts
	}
	account := r.account
	return &account, nil
}

func (r *leoIntegrationAccountRepo) ListSchedulableByGroupIDAndPlatform(_ context.Context, _ int64, platform string) ([]service.Account, error) {
	return r.accountsForPlatform(platform), nil
}

func (r *leoIntegrationAccountRepo) ListSchedulableByPlatform(_ context.Context, platform string) ([]service.Account, error) {
	return r.accountsForPlatform(platform), nil
}

func (r *leoIntegrationAccountRepo) ListSchedulableUngroupedByPlatform(_ context.Context, platform string) ([]service.Account, error) {
	return r.accountsForPlatform(platform), nil
}

func (r *leoIntegrationAccountRepo) accountsForPlatform(platform string) []service.Account {
	if platform != service.PlatformLeo {
		return nil
	}
	return []service.Account{r.account}
}

type leoIntegrationUsageRepo struct {
	service.UsageLogRepository
	lastLog *service.UsageLog
}

func (r *leoIntegrationUsageRepo) Create(_ context.Context, log *service.UsageLog) (bool, error) {
	r.lastLog = log
	return true, nil
}

type leoIntegrationHTTPUpstream struct {
	client *http.Client
}

func (u *leoIntegrationHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return u.client.Do(req)
}

func (u *leoIntegrationHTTPUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, "", 0, 0)
}

func TestLeoVideoGenerationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const leoAPIKey = "leo-integration-secret"
	upstreamBody := make(chan []byte, 1)

	leoStudio := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+leoAPIKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/health":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		case "/v1/videos/generations":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			body, _ := json.Marshal(payload)
			upstreamBody <- body
			if payload["model"] != "seedance-2.0-fast" {
				http.Error(w, "unexpected mapped model", http.StatusUnprocessableEntity)
				return
			}
			time.Sleep(10 * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "req_leo_integration")
			_, _ = w.Write([]byte(`{"data":[{"mp4_url":"https://cdn.example/leo.mp4"}],"provider":{"generation_id":"leo_gen_123","resolution":"RESOLUTION_720","duration":12}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer leoStudio.Close()

	wrongHealthReq, err := http.NewRequest(http.MethodGet, leoStudio.URL+"/health", nil)
	require.NoError(t, err)
	wrongHealthReq.Header.Set("Authorization", "Bearer wrong-key")
	wrongHealthResp, err := leoStudio.Client().Do(wrongHealthReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, wrongHealthResp.StatusCode)
	require.NoError(t, wrongHealthResp.Body.Close())

	healthReq, err := http.NewRequest(http.MethodGet, leoStudio.URL+"/health", nil)
	require.NoError(t, err)
	healthReq.Header.Set("Authorization", "Bearer "+leoAPIKey)
	healthResp, err := leoStudio.Client().Do(healthReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, healthResp.StatusCode)
	require.NoError(t, healthResp.Body.Close())

	groupID := int64(9101)
	price480P := 0.08
	price720P := 0.125
	price1080P := 0.25
	accountRepo := &leoIntegrationAccountRepo{account: service.Account{
		ID:          9201,
		Name:        "leo-integration",
		Platform:    service.PlatformLeo,
		Type:        service.AccountTypeAPIKey,
		Status:      service.StatusActive,
		Schedulable: true,
		Concurrency: 1,
		Credentials: map[string]any{
			"base_url": leoStudio.URL + "/v1",
			"api_key":  leoAPIKey,
			"model_mapping": map[string]any{
				"seedance-2.0": "seedance-2.0-fast",
			},
		},
	}}
	usageRepo := &leoIntegrationUsageRepo{}
	cfg := &config.Config{RunMode: config.RunModeSimple}
	cfg.Default.RateMultiplier = 1
	billingCacheService := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheService.Stop)
	httpUpstream := &leoIntegrationHTTPUpstream{client: leoStudio.Client()}
	gatewayService := service.NewOpenAIGatewayService(
		accountRepo,
		usageRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg,
		nil,
		nil,
		service.NewBillingService(cfg, nil),
		nil,
		billingCacheService,
		httpUpstream,
		&service.DeferredService{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	require.Equal(t, leoStudio.URL+"/v1", accountRepo.account.GetLeoBaseURL())
	require.Equal(t, leoAPIKey, accountRepo.account.GetLeoAPIKey())
	require.Equal(t, "seedance-2.0-fast", accountRepo.account.GetMappedModel("seedance-2.0"))
	handler := NewOpenAIGatewayHandler(
		gatewayService,
		service.NewConcurrencyService(nil),
		billingCacheService,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg),
		nil,
		nil,
		nil,
		nil,
		cfg,
	)
	apiKey := &service.APIKey{
		ID:      9301,
		GroupID: &groupID,
		User:    &service.User{ID: 9401, Status: service.StatusActive},
		Group: &service.Group{
			ID:                   groupID,
			Platform:             service.PlatformLeo,
			Status:               service.StatusActive,
			RateMultiplier:       1.5,
			AllowImageGeneration: true,
			VideoPrice480P:       &price480P,
			VideoPrice720P:       &price720P,
			VideoPrice1080P:      &price1080P,
		},
	}

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: apiKey.User.ID, Concurrency: 0})
		c.Next()
	})
	router.POST("/v1/videos/generations", handler.LeoVideoGeneration)

	requestBody := []byte(`{"model":"seedance-2.0","prompt":"A cinematic city at night","resolution":"480p","duration":8,"audio":false}`)
	request := httptest.NewRequest(http.MethodPost, "/v1/videos/generations", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	var forwardedBody []byte
	select {
	case forwardedBody = <-upstreamBody:
	case <-time.After(time.Second):
		t.Fatal("LeoStudio did not receive the generation request")
	}

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, "https://cdn.example/leo.mp4", gjson.GetBytes(response.Body.Bytes(), "data.0.mp4_url").String())
	require.Equal(t, "leo_gen_123", gjson.GetBytes(response.Body.Bytes(), "provider.generation_id").String())

	require.NotEmpty(t, forwardedBody)
	require.Equal(t, "seedance-2.0-fast", gjson.GetBytes(forwardedBody, "model").String())
	require.Equal(t, "A cinematic city at night", gjson.GetBytes(forwardedBody, "prompt").String())
	require.Equal(t, "480p", gjson.GetBytes(forwardedBody, "resolution").String())
	require.Equal(t, int64(8), gjson.GetBytes(forwardedBody, "duration").Int())

	require.Equal(t, service.PlatformLeo, service.PlatformFromAPIKey(apiKey))
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "seedance-2.0", usageRepo.lastLog.RequestedModel)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, "seedance-2.0-fast", *usageRepo.lastLog.UpstreamModel)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(service.BillingModeVideo), *usageRepo.lastLog.BillingMode)
	require.NotNil(t, usageRepo.lastLog.VideoResolution)
	require.Equal(t, "720p", *usageRepo.lastLog.VideoResolution)
	require.NotNil(t, usageRepo.lastLog.VideoDurationSeconds)
	require.Equal(t, 12, *usageRepo.lastLog.VideoDurationSeconds)
	require.InDelta(t, price720P*12*1.5, usageRepo.lastLog.ActualCost, 1e-12)
}
