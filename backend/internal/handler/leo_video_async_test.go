package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type handlerVideoJobRepo struct {
	service.VideoJobRepository
	job *service.VideoJob
}

func (r *handlerVideoJobRepo) CreateVideoJob(_ context.Context, job *service.VideoJob) error {
	copy := *job
	r.job = &copy
	return nil
}

func (r *handlerVideoJobRepo) GetVideoJobForAPIKey(_ context.Context, id string, apiKeyID int64) (*service.VideoJob, error) {
	if r.job == nil || r.job.JobID != id || r.job.APIKeyID != apiKeyID {
		return nil, service.ErrVideoJobNotFound
	}
	copy := *r.job
	return &copy, nil
}

func (r *handlerVideoJobRepo) ListVideoJobsForAPIKey(_ context.Context, apiKeyID int64, _ int, _ string) ([]*service.VideoJob, error) {
	if r.job == nil || r.job.APIKeyID != apiKeyID {
		return []*service.VideoJob{}, nil
	}
	copy := *r.job
	return []*service.VideoJob{&copy}, nil
}

func (r *handlerVideoJobRepo) TransitionVideoJob(_ context.Context, id string, _ []string, status string, transition service.VideoJobTransition) error {
	if r.job == nil || r.job.JobID != id {
		return service.ErrVideoJobNotFound
	}
	r.job.Status = status
	if transition.UpstreamJobID != nil {
		r.job.UpstreamJobID = *transition.UpstreamJobID
	}
	return nil
}

type handlerVideoAccountSelector struct{ account *service.Account }

func (s handlerVideoAccountSelector) Select(context.Context, int64, string, map[int64]struct{}) (*service.Account, error) {
	return s.account, nil
}

func (s handlerVideoAccountSelector) GetByID(context.Context, int64) (*service.Account, error) {
	return s.account, nil
}

type handlerVideoClient struct{}

func (handlerVideoClient) CreateLeoAsyncVideo(context.Context, *service.Account, []byte) (*service.LeoAsyncAccepted, error) {
	return &service.LeoAsyncAccepted{JobID: 777, Status: service.VideoJobPending}, nil
}

func (handlerVideoClient) GetLeoAsyncVideo(context.Context, *service.Account, int64) (*service.LeoAsyncJob, error) {
	return nil, nil
}

func (handlerVideoClient) CancelLeoAsyncVideo(context.Context, *service.Account, int64) (*service.LeoAsyncJob, error) {
	return &service.LeoAsyncJob{JobID: 777, Status: service.VideoJobCanceled}, nil
}

type handlerRejectingVideoClient struct{ handlerVideoClient }

func (handlerRejectingVideoClient) CreateLeoAsyncVideo(context.Context, *service.Account, []byte) (*service.LeoAsyncAccepted, error) {
	return nil, &service.LeoAsyncUpstreamError{StatusCode: http.StatusBadRequest, Message: "guidances.image_reference supports at most 4 items"}
}

type handlerVideoBillingRepo struct{ service.UsageBillingRepository }

func (handlerVideoBillingRepo) ReserveVideoBalance(context.Context, *service.VideoBalanceHoldCommand) (*service.VideoBalanceHoldResult, error) {
	return &service.VideoBalanceHoldResult{Applied: true}, nil
}

func (handlerVideoBillingRepo) ReleaseVideoBalance(context.Context, *service.VideoBalanceHoldCommand) (*service.VideoBalanceHoldResult, error) {
	return &service.VideoBalanceHoldResult{Applied: true}, nil
}

func newHandlerVideoJobService(repo *handlerVideoJobRepo) *service.VideoJobService {
	selector := handlerVideoAccountSelector{account: &service.Account{ID: 9, Platform: service.PlatformLeo, Type: service.AccountTypeAPIKey, Credentials: map[string]any{
		"base_url": "http://leo.internal:8000/v1", "api_key": "secret", "model_mapping": map[string]any{"seedance": "seedance-2.0"},
	}}}
	billing := &service.VideoJobBillingService{BillingRepo: handlerVideoBillingRepo{}}
	return service.NewVideoJobService(repo, selector, handlerVideoClient{}, billing)
}

func TestLeoVideoAsyncGenerationReturnsPublic202Mapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(`{"model":"seedance","prompt":"waves"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Prefer", "respond-async")
	setHandlerVideoAuth(c)

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusAccepted, recorder.Code)
	require.Equal(t, "respond-async", recorder.Header().Get("Preference-Applied"))
	require.Equal(t, "/v1/videos/jobs/"+repo.job.JobID, recorder.Header().Get("Location"))
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, repo.job.JobID, body["job_id"])
	require.NotContains(t, recorder.Body.String(), "777")
	require.NotContains(t, recorder.Body.String(), "account_id")
}

func TestLeoVideoAsyncGenerationTracksMultipleLocalGuidanceInputs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	first := "01234567890123456789012345678901"
	second := "abcdefghijklmnopqrstuvwxyzABCDEF"
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo), videoInputHandler: NewVideoInputHandler(store)}
	body := `{"model":"seedance","prompt":"waves","image_urls":["` + store.InternalURL(first) + `","` + store.InternalURL(second) + `"]}`
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Prefer", "respond-async")
	setHandlerVideoAuth(c)

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusAccepted, recorder.Code)
	require.Equal(t, first+","+second, repo.job.LocalInputName)
}

func TestLeoVideoAsyncGenerationPreservesNewUpstreamValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	selector := handlerVideoAccountSelector{account: &service.Account{ID: 9, Platform: service.PlatformLeo, Type: service.AccountTypeAPIKey, Credentials: map[string]any{
		"base_url": "http://leo.internal:8000/v1", "api_key": "secret", "model_mapping": map[string]any{"seedance": "seedance-2.0"},
	}}}
	billing := &service.VideoJobBillingService{BillingRepo: handlerVideoBillingRepo{}}
	h := &OpenAIGatewayHandler{videoJobService: service.NewVideoJobService(repo, selector, handlerRejectingVideoClient{}, billing)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(`{"model":"seedance","prompt":"waves"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Prefer", "respond-async")
	setHandlerVideoAuth(c)

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "guidances.image_reference supports at most 4 items")
}

func TestLeoVideoAsyncJobEndpointsStayAPIKeyScoped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{JobID: "vidjob_public", APIKeyID: 2, Status: service.VideoJobPending, UpstreamJobID: 777}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}

	for _, apiKeyID := range []int64{999, 2} {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs/vidjob_public", nil)
		c.Params = gin.Params{{Key: "job_id", Value: "vidjob_public"}}
		setHandlerVideoAuthWithKey(c, apiKeyID)
		h.LeoVideoJob(c)
		if apiKeyID == 999 {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		} else {
			require.Equal(t, http.StatusOK, recorder.Code)
			require.NotContains(t, recorder.Body.String(), "777")
		}
	}
}

func TestLeoVideoJobContentStaysAPIKeyScopedAndServesMP4(t *testing.T) {
	gin.SetMode(gin.TestMode)
	video := []byte{0, 0, 0, 16, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0, 0, 0}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(video)
	}))
	defer server.Close()
	store := service.NewVideoOutputStore(t.TempDir())
	_, err := store.Save(context.Background(), "vidjob_content", json.RawMessage(`{"data":[{"url":"`+server.URL+`/video.mp4"}]}`))
	require.NoError(t, err)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{JobID: "vidjob_content", APIKeyID: 2, Status: service.VideoJobCompleted}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo), videoOutputStore: store}

	for _, apiKeyID := range []int64{999, 2} {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs/vidjob_content/content", nil)
		c.Params = gin.Params{{Key: "job_id", Value: "vidjob_content"}}
		setHandlerVideoAuthWithKey(c, apiKeyID)
		h.LeoVideoJobContent(c)
		if apiKeyID == 999 {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		} else {
			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, "video/mp4", recorder.Header().Get("Content-Type"))
			require.Equal(t, video, recorder.Body.Bytes())
		}
	}
}

func setHandlerVideoAuth(c *gin.Context) { setHandlerVideoAuthWithKey(c, 2) }

func setHandlerVideoAuthWithKey(c *gin.Context, apiKeyID int64) {
	groupID := int64(3)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{ID: apiKeyID, UserID: 1, GroupID: &groupID, Group: &service.Group{
		ID: groupID, Platform: service.PlatformLeo, AllowImageGeneration: true,
		VideoPrice480P: handlerFloatPointer(0.05), VideoPrice720P: handlerFloatPointer(0.1), VideoPrice1080P: handlerFloatPointer(0.2),
	}})
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 1})
}

func handlerFloatPointer(value float64) *float64 { return &value }
