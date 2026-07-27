package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/securityaudit"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLeoVideoGenerationValidatesRequestBeforeDependencies(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "invalid JSON", body: `{"model":`, want: "valid JSON"},
		{name: "missing model", body: `{"prompt":"waves"}`, want: "model is required"},
		{name: "missing prompt", body: `{"model":"seedance-2.0"}`, want: "prompt is required"},
		{name: "fast unsupported resolution", body: `{"model":"seedance-2.0-fast","prompt":"waves","resolution":"1080p"}`, want: "resolution is not supported"},
		{name: "mini unsupported resolution", body: `{"model":"seedance-2.0-mini","prompt":"waves","resolution":"1080p"}`, want: "resolution is not supported"},
		{name: "mini unsupported aspect", body: `{"model":"seedance-2.0-mini","prompt":"waves","aspect_ratio":"4:3"}`, want: "aspect_ratio is not supported"},
		{name: "unsupported duration", body: `{"model":"seedance-2.0","prompt":"waves","duration":3}`, want: "duration must be a whole number"},
		{name: "standard 1080p duration above range", body: `{"model":"seedance-2.0","prompt":"waves","resolution":"1080p","duration":13}`, want: "duration must be a whole number from 4 through 12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, c := newLeoVideoHandlerTestContext(tt.body, true)

			(&OpenAIGatewayHandler{}).LeoVideoGeneration(c)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.Contains(t, rec.Body.String(), tt.want)
		})
	}
}

func TestLeoVideoGenerationRequiresVideoPermission(t *testing.T) {
	rec, c := newLeoVideoHandlerTestContext(`{"model":"seedance-2.0","prompt":"waves"}`, false)

	(&OpenAIGatewayHandler{}).LeoVideoGeneration(c)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), service.ImageGenerationPermissionMessage())
}

func TestLeoVideoGenerationSecurityAuditRunsBeforeScheduling(t *testing.T) {
	rec, c := newLeoVideoHandlerTestContext(`{"model":"seedance-2.0","prompt":"blocked video prompt"}`, true)
	engine := blockingHandlerPromptEngine()
	h := &OpenAIGatewayHandler{
		gatewayService:           &service.OpenAIGatewayService{},
		billingCacheService:      &service.BillingCacheService{},
		apiKeyService:            &service.APIKeyService{},
		concurrencyHelper:        &ConcurrencyHelper{concurrencyService: &service.ConcurrencyService{}},
		securityAuditCoordinator: securityaudit.NewCoordinator(nil, engine),
	}

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), securityaudit.ErrorCodeBlocked)
	evaluated, _, requests := engine.snapshot()
	require.Equal(t, 1, evaluated)
	require.Len(t, requests, 1)
	require.Contains(t, string(requests[0].Body), "blocked video prompt")
	require.Equal(t, service.PlatformLeo, requests[0].Provider)
}

func newLeoVideoHandlerTestContext(body string, allowVideo bool) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	groupID := int64(1)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{
		ID:      2,
		GroupID: &groupID,
		Group: &service.Group{
			ID:                   groupID,
			Platform:             service.PlatformLeo,
			AllowImageGeneration: allowVideo,
		},
	})
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 3, Concurrency: 1})
	return rec, c
}
