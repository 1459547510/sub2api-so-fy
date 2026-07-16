package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
