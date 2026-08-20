package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
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

func (r *handlerVideoJobRepo) ListVideoJobsForAPIKey(_ context.Context, apiKeyID int64, _ int, _ int, _ string) ([]*service.VideoJob, int, error) {
	if r.job == nil || r.job.APIKeyID != apiKeyID {
		return []*service.VideoJob{}, 0, nil
	}
	copy := *r.job
	return []*service.VideoJob{&copy}, 1, nil
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
	return nil, &service.LeoAsyncUpstreamError{StatusCode: http.StatusBadRequest, Message: "Leonardo: guidances.image_reference supports at most 4 items"}
}

type handlerVideoBillingRepo struct {
	service.UsageBillingRepository
	reserves   int
	releases   int
	reserveErr error
}

func (r *handlerVideoBillingRepo) ReserveVideoBalance(context.Context, *service.VideoBalanceHoldCommand) (*service.VideoBalanceHoldResult, error) {
	if r != nil {
		r.reserves++
		if r.reserveErr != nil {
			return nil, r.reserveErr
		}
	}
	return &service.VideoBalanceHoldResult{Applied: true}, nil
}

func (r *handlerVideoBillingRepo) ReleaseVideoBalance(context.Context, *service.VideoBalanceHoldCommand) (*service.VideoBalanceHoldResult, error) {
	if r != nil {
		r.releases++
	}
	return &service.VideoBalanceHoldResult{Applied: true}, nil
}

func newHandlerVideoJobService(repo *handlerVideoJobRepo) *service.VideoJobService {
	return newHandlerVideoJobServiceWithBalance(repo, &handlerVideoBillingRepo{})
}

func newHandlerVideoJobServiceWithBalance(repo *handlerVideoJobRepo, balance *handlerVideoBillingRepo) *service.VideoJobService {
	selector := handlerVideoAccountSelector{account: &service.Account{ID: 9, Platform: service.PlatformLeo, Type: service.AccountTypeAPIKey, Credentials: map[string]any{
		"base_url": "http://leo.internal:8000/v1", "api_key": "secret", "model_mapping": map[string]any{"seedance": "seedance-2.0"},
	}}}
	if balance == nil {
		balance = &handlerVideoBillingRepo{}
	}
	billing := &service.VideoJobBillingService{BillingRepo: balance}
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

func TestLeoVideoAsyncGenerationRejectsUnsupportedMiniAspect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(&handlerVideoJobRepo{})}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/generations", strings.NewReader(`{"model":"seedance-2.0-mini","prompt":"waves","aspect_ratio":"9:21"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Prefer", "respond-async")
	setHandlerVideoAuth(c)

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "aspect_ratio is not supported")
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
	billing := &service.VideoJobBillingService{BillingRepo: &handlerVideoBillingRepo{}}
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
	require.NotContains(t, strings.ToLower(recorder.Body.String()), "leonardo")
}

func TestLeoVideoJobErrorsHideUpstreamProviderName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_failed", APIKeyID: 2, Status: service.VideoJobFailed,
		ErrorMessage: "LeoStudio failed after Leonardo.ai rejected the request",
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs/vidjob_failed", nil)
	c.Params = gin.Params{{Key: "job_id", Value: "vidjob_failed"}}
	setHandlerVideoAuth(c)

	h.LeoVideoJob(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, strings.ToLower(recorder.Body.String()), "video service")
	require.NotContains(t, strings.ToLower(recorder.Body.String()), "leonardo")
	require.NotContains(t, strings.ToLower(recorder.Body.String()), "leostudio")
}

func TestLeoVideoJobResultHidesUpstreamMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_completed", APIKeyID: 2, Status: service.VideoJobCompleted,
		Result: json.RawMessage(`{"data":[{"mp4_url":"/v1/videos/jobs/vidjob_completed/content","url":"/v1/videos/jobs/vidjob_completed/content","local_url":"/v1/videos/jobs/vidjob_completed/content","video_url":"https://cdn.example/video.mp4","source_url":"https://cdn.example/video.mp4","generation_id":"provider-job-42","provider":{"uuid":"provider-uuid"}}],"provider":{"generation_id":"provider-job-42","account_id":"provider-account"}}`),
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs/vidjob_completed", nil)
	c.Params = gin.Params{{Key: "job_id", Value: "vidjob_completed"}}
	setHandlerVideoAuth(c)

	h.LeoVideoJob(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	body := recorder.Body.String()
	require.Contains(t, body, "/v1/videos/jobs/vidjob_completed/content")
	require.NotContains(t, body, "provider")
	require.NotContains(t, body, "cdn.example")
	require.NotContains(t, body, "generation_id")
	require.NotContains(t, body, "account_id")
	require.NotContains(t, body, "uuid")
}

func TestLeoVideoPassthroughMessagesHideUpstreamProviderName(t *testing.T) {
	message := "Leonardo.ai rejected the request through LeoStudio"

	require.Equal(t, "video service rejected the request through video service", publicUpstreamErrorMessage(service.PlatformLeo, message))
	require.Equal(t, message, publicUpstreamErrorMessage(service.PlatformOpenAI, message))
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

func TestLeoVideoAsyncJobHidesLegacyProviderNames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_failed", APIKeyID: 2, Status: service.VideoJobFailed,
		ErrorMessage: "Leonardo.ai request rejected by LeoStudio",
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs/vidjob_failed", nil)
	c.Params = gin.Params{{Key: "job_id", Value: "vidjob_failed"}}
	setHandlerVideoAuth(c)

	h.LeoVideoJob(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), "Video service request rejected by Video service")
	require.NotContains(t, recorder.Body.String(), "Leonardo")
	require.NotContains(t, recorder.Body.String(), "LeoStudio")
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

func TestLeoVideoGenerationRejectsInsufficientSyncHold(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{RunMode: config.RunModeSimple}
	billingCache := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCache.Stop)
	h := NewOpenAIGatewayHandler(
		&service.OpenAIGatewayService{},
		service.NewConcurrencyService(nil),
		billingCache,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg),
		nil, nil, nil, nil, cfg,
	)
	jobRepo := &handlerVideoJobRepo{}
	balance := &handlerVideoBillingRepo{reserveErr: service.ErrVideoInsufficientBalance}
	videoJobs := newHandlerVideoJobServiceWithBalance(jobRepo, balance)
	h.SetVideoServices(videoJobs, nil, nil)
	rec, c := newLeoVideoHandlerTestContext(`{"model":"seedance-2.5","prompt":"waves","resolution":"720p","duration":30}`, true)
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 3, Concurrency: 0})
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	require.True(t, ok)
	price := 0.8
	apiKey.User = &service.User{ID: 3}
	apiKey.Group.RateMultiplier = 1
	apiKey.Group.VideoPrice480P = &price
	apiKey.Group.VideoPrice720P = &price
	apiKey.Group.VideoPrice1080P = &price
	h.LeoVideoGeneration(c)

	require.Equal(t, 1, balance.reserves, rec.Body.String())
	require.Equal(t, http.StatusPaymentRequired, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "insufficient balance")
	require.NotContains(t, strings.ToLower(rec.Body.String()), "leonardo")
	require.Nil(t, jobRepo.job)
	require.Equal(t, 0, balance.releases)
}

func TestPersistSyncVideoJobWritesFailedLogWithoutVendorNames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	_, c := newLeoVideoHandlerTestContext(`{"model":"seedance","prompt":"waves"}`, true)
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	require.True(t, ok)
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	require.True(t, ok)

	h.persistSyncVideoJob(c, apiKey, subject, &service.Account{ID: 9, Platform: service.PlatformLeo},
		[]byte(`{"model":"seedance","prompt":"waves"}`), service.VideoJobFailed, nil,
		&service.LeoVideoRejectedError{StatusCode: http.StatusUnprocessableEntity, Message: "Leonardo rejected the prompt"})

	require.NotNil(t, repo.job)
	require.Equal(t, service.VideoJobFailed, repo.job.Status)
	require.Equal(t, "waves", repo.job.Prompt)
	require.Nil(t, repo.job.HoldAmount)
	require.NotContains(t, repo.job.ErrorMessage, "Leonardo")
	require.Contains(t, strings.ToLower(repo.job.ErrorMessage), "video service")
}

func TestLeoVideoJobsReturnsPagedTotal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_list", APIKeyID: 2, Status: service.VideoJobPending, Prompt: "waves",
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/jobs?limit=20&offset=0", nil)
	setHandlerVideoAuth(c)

	h.LeoVideoJobs(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, float64(1), body["total"])
	require.Equal(t, float64(20), body["limit"])
	require.Equal(t, float64(0), body["offset"])
	require.Len(t, body["data"], 1)
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
