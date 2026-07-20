package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeVideoJobServiceRepo struct {
	VideoJobRepository
	job         *VideoJob
	transitions []struct {
		jobID string
		to    string
		data  VideoJobTransition
	}
	list []*VideoJob
}

func (r *fakeVideoJobServiceRepo) CreateVideoJob(_ context.Context, job *VideoJob) error {
	r.job = cloneVideoJob(job)
	return nil
}

func (r *fakeVideoJobServiceRepo) GetVideoJob(_ context.Context, jobID string) (*VideoJob, error) {
	if r.job == nil || r.job.JobID != jobID {
		return nil, ErrVideoJobNotFound
	}
	return cloneVideoJob(r.job), nil
}

func (r *fakeVideoJobServiceRepo) GetVideoJobForAPIKey(_ context.Context, jobID string, apiKeyID int64) (*VideoJob, error) {
	if r.job == nil || r.job.JobID != jobID || r.job.APIKeyID != apiKeyID {
		return nil, ErrVideoJobNotFound
	}
	return cloneVideoJob(r.job), nil
}

func (r *fakeVideoJobServiceRepo) ListVideoJobsForAPIKey(_ context.Context, apiKeyID int64, _ int, _ string) ([]*VideoJob, error) {
	if r.job != nil && r.job.APIKeyID == apiKeyID {
		return []*VideoJob{cloneVideoJob(r.job)}, nil
	}
	return r.list, nil
}

func (r *fakeVideoJobServiceRepo) ListActiveVideoJobs(_ context.Context, _ int) ([]*VideoJob, error) {
	if r.job == nil || IsTerminalVideoJobStatus(r.job.Status) {
		return nil, nil
	}
	return []*VideoJob{cloneVideoJob(r.job)}, nil
}

func (r *fakeVideoJobServiceRepo) TransitionVideoJob(_ context.Context, jobID string, _ []string, status string, transition VideoJobTransition) error {
	r.transitions = append(r.transitions, struct {
		jobID string
		to    string
		data  VideoJobTransition
	}{jobID: jobID, to: status, data: transition})
	if r.job != nil && r.job.JobID == jobID {
		r.job.Status = status
		if transition.AccountID != nil {
			r.job.AccountID = *transition.AccountID
		}
		if transition.UpstreamJobID != nil {
			r.job.UpstreamJobID = *transition.UpstreamJobID
		}
		if transition.SettledAt != nil {
			r.job.SettledAt = transition.SettledAt
		}
		if transition.FinishedAt != nil {
			r.job.FinishedAt = transition.FinishedAt
		}
		if transition.ErrorMessage != nil {
			r.job.ErrorMessage = *transition.ErrorMessage
		}
		if len(transition.Result) != 0 {
			r.job.Result = append(json.RawMessage(nil), transition.Result...)
		}
		if transition.ActualCost != nil {
			r.job.ActualCost = transition.ActualCost
		}
	}
	return nil
}

type fakeVideoJobSelector struct {
	accounts []*Account
	index    int
}

func (s *fakeVideoJobSelector) Select(context.Context, int64, string, map[int64]struct{}) (*Account, error) {
	if s.index >= len(s.accounts) {
		return nil, errors.New("no leo account")
	}
	account := s.accounts[s.index]
	s.index++
	return account, nil
}

func (s *fakeVideoJobSelector) GetByID(_ context.Context, id int64) (*Account, error) {
	for _, account := range s.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

type fakeVideoJobClient struct {
	responses []*LeoAsyncAccepted
	errors    []error
	creates   int
	cancels   int
}

func (c *fakeVideoJobClient) CreateLeoAsyncVideo(context.Context, *Account, []byte) (*LeoAsyncAccepted, error) {
	index := c.creates
	c.creates++
	if index < len(c.errors) && c.errors[index] != nil {
		return nil, c.errors[index]
	}
	return c.responses[index], nil
}

func (c *fakeVideoJobClient) GetLeoAsyncVideo(context.Context, *Account, int64) (*LeoAsyncJob, error) {
	return &LeoAsyncJob{JobID: 42, Status: VideoJobPending}, nil
}

func (c *fakeVideoJobClient) CancelLeoAsyncVideo(context.Context, *Account, int64) (*LeoAsyncJob, error) {
	c.cancels++
	return &LeoAsyncJob{JobID: 42, Status: VideoJobCanceled}, nil
}

func newVideoJobServiceTestBilling() *VideoJobBillingService {
	return &VideoJobBillingService{BillingRepo: &fakeVideoJobBalanceRepo{}}
}

func newVideoJobServiceTestAPIKey() *APIKey {
	groupID := int64(3)
	return &APIKey{ID: 2, GroupID: &groupID, Group: &Group{
		ID: groupID, Platform: PlatformLeo, VideoPrice480P: f64p(0.05), VideoPrice720P: f64p(0.1), VideoPrice1080P: f64p(0.2), RateMultiplier: 1,
	}}
}

func newVideoJobServiceTestAccount(id int64) *Account {
	return &Account{ID: id, Platform: PlatformLeo, Type: AccountTypeAPIKey, Credentials: map[string]any{
		"base_url": "http://leo.internal:8000/v1", "api_key": "secret",
		"model_mapping": map[string]any{"seedance": "seedance-2.0", "seedance-2.0": "seedance-2.0"},
	}}
}

func TestVideoJobServiceCreateHoldsBeforeUpstreamAndPersistsMapping(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{}
	client := &fakeVideoJobClient{responses: []*LeoAsyncAccepted{{JobID: 42, Status: VideoJobPending}}}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	service := &VideoJobService{Repo: repo, Selector: selector, Client: client, Billing: newVideoJobServiceTestBilling()}

	job, err := service.Create(context.Background(), CreateVideoJobInput{APIKey: newVideoJobServiceTestAPIKey(), User: &User{ID: 1}, Body: []byte(`{"model":"seedance","prompt":"waves","resolution":"720p","duration":8}`)})

	require.NoError(t, err)
	require.NotNil(t, job)
	require.Equal(t, "seedance", job.RequestedModel)
	require.Equal(t, "seedance-2.0", job.UpstreamModel)
	require.Equal(t, int64(42), job.UpstreamJobID)
	require.Equal(t, int64(9), job.AccountID)
	require.Equal(t, VideoJobPending, job.Status)
	require.NotNil(t, repo.job)
	require.NotEmpty(t, repo.job.BillingSnapshot)
	require.Len(t, repo.transitions, 1)
	require.Equal(t, int64(42), *repo.transitions[0].data.UpstreamJobID)
	require.Equal(t, 1, client.creates)
}

func TestVideoJobServiceFailsOverBeforeAcceptedButNotAfterAmbiguousFailure(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{}
	client := &fakeVideoJobClient{
		responses: []*LeoAsyncAccepted{nil, &LeoAsyncAccepted{JobID: 43, Status: VideoJobPending}},
		errors:    []error{&LeoAsyncUpstreamError{StatusCode: 503, Message: "temporary", Retryable: true}},
	}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9), newVideoJobServiceTestAccount(10)}}
	service := &VideoJobService{Repo: repo, Selector: selector, Client: client, Billing: newVideoJobServiceTestBilling()}

	job, err := service.Create(context.Background(), CreateVideoJobInput{APIKey: newVideoJobServiceTestAPIKey(), User: &User{ID: 1}, Body: []byte(`{"model":"seedance-2.0","prompt":"waves"}`)})
	require.NoError(t, err)
	require.Equal(t, int64(43), job.UpstreamJobID)
	require.Equal(t, int64(10), job.AccountID)
	require.GreaterOrEqual(t, len(repo.transitions), 2)
	require.Equal(t, int64(10), *repo.transitions[0].data.AccountID)

	repo = &fakeVideoJobServiceRepo{}
	client = &fakeVideoJobClient{errors: []error{&LeoAsyncUpstreamError{Message: "ambiguous", Ambiguous: true}}}
	selector = &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9), newVideoJobServiceTestAccount(10)}}
	service = &VideoJobService{Repo: repo, Selector: selector, Client: client, Billing: newVideoJobServiceTestBilling()}
	_, err = service.Create(context.Background(), CreateVideoJobInput{APIKey: newVideoJobServiceTestAPIKey(), User: &User{ID: 1}, Body: []byte(`{"model":"seedance-2.0","prompt":"waves"}`)})
	require.Error(t, err)
	require.Equal(t, 1, client.creates)
	require.Equal(t, VideoJobFailed, repo.job.Status)
}

func TestVideoJobServiceScopesAccessAndOnlyCancelsPending(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{job: &VideoJob{JobID: "vidjob_1", APIKeyID: 2, UserID: 1, AccountID: 9, UpstreamJobID: 42, Status: VideoJobPending}}
	client := &fakeVideoJobClient{}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	service := &VideoJobService{Repo: repo, Selector: selector, Client: client, Billing: newVideoJobServiceTestBilling()}

	_, err := service.Get(context.Background(), "vidjob_1", 999)
	require.ErrorIs(t, err, ErrVideoJobNotFound)
	_, err = service.Cancel(context.Background(), "vidjob_1", 2)
	require.NoError(t, err)
	require.Equal(t, 1, client.cancels)
	require.Equal(t, VideoJobCanceled, repo.job.Status)

	repo.job.Status = VideoJobRunning
	_, err = service.Cancel(context.Background(), "vidjob_1", 2)
	require.ErrorIs(t, err, ErrVideoJobCancelConflict)
}

func TestVideoJobServiceMarksCanceledLocalInputTerminal(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)
	input, err := store.Save(bytes.NewReader([]byte{137, 80, 78, 71, 13, 10, 26, 10}))
	require.NoError(t, err)
	repo := &fakeVideoJobServiceRepo{job: &VideoJob{JobID: "vidjob_local", APIKeyID: 2, UserID: 1, AccountID: 9, UpstreamJobID: 42, Status: VideoJobPending, LocalInputName: input.Token}}
	service := &VideoJobService{Repo: repo, Selector: &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}, Client: &fakeVideoJobClient{}, Billing: newVideoJobServiceTestBilling()}
	service.SetVideoInputStore(store)

	_, err = service.Cancel(context.Background(), "vidjob_local", 2)
	require.NoError(t, err)
	removed, err := store.Cleanup(time.Now().Add(30 * time.Minute))
	require.NoError(t, err)
	require.Equal(t, 0, removed)
	removed, err = store.Cleanup(time.Now().Add(2 * time.Hour))
	require.NoError(t, err)
	require.Equal(t, 1, removed)
}

func TestVideoJobServiceValidatesPromptBeforeSelectingAccount(t *testing.T) {
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	service := &VideoJobService{Repo: &fakeVideoJobServiceRepo{}, Selector: selector, Client: &fakeVideoJobClient{}, Billing: newVideoJobServiceTestBilling()}

	_, err := service.Create(context.Background(), CreateVideoJobInput{APIKey: newVideoJobServiceTestAPIKey(), User: &User{ID: 1}, Body: []byte(`{"model":"seedance-2.0"}`)})

	require.Error(t, err)
	require.Equal(t, 0, selector.index)
}

func cloneVideoJob(job *VideoJob) *VideoJob {
	if job == nil {
		return nil
	}
	copy := *job
	copy.BillingSnapshot = append([]byte(nil), job.BillingSnapshot...)
	return &copy
}
