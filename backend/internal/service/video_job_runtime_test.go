package service

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type runtimeVideoClient struct {
	job       *LeoAsyncJob
	err       error
	started   chan struct{}
	continueC chan struct{}
	calls     int
	mu        sync.Mutex
}

func (c *runtimeVideoClient) CreateLeoAsyncVideo(context.Context, *Account, []byte) (*LeoAsyncAccepted, error) {
	return nil, nil
}

func (c *runtimeVideoClient) GetLeoAsyncVideo(_ context.Context, _ *Account, _ int64) (*LeoAsyncJob, error) {
	c.mu.Lock()
	c.calls++
	c.mu.Unlock()
	if c.started != nil {
		select {
		case <-c.started:
		default:
			close(c.started)
		}
	}
	if c.continueC != nil {
		<-c.continueC
	}
	return c.job, c.err
}

func (c *runtimeVideoClient) CancelLeoAsyncVideo(context.Context, *Account, int64) (*LeoAsyncJob, error) {
	return &LeoAsyncJob{JobID: 42, Status: VideoJobCanceled}, nil
}

func newRuntimeBilling() *VideoJobBillingService {
	snapshot, _ := json.Marshal(VideoJobBillingSnapshot{Version: 1, BillingType: BillingTypeBalance, Price720P: 0.1, RateMultiplier: 1})
	_ = snapshot
	return &VideoJobBillingService{
		BillingRepo:   &fakeVideoJobBalanceRepo{},
		UsageRecorder: &fakeVideoUsageRecorder{},
		APIKeys:       fakeVideoBillingAPIKeyLoader{apiKey: &APIKey{ID: 2}},
		Users:         fakeVideoBillingUserLoader{user: &User{ID: 1}},
		Accounts:      fakeVideoBillingAccountLoader{account: &Account{ID: 9, Type: AccountTypeAPIKey}},
		Subscriptions: fakeVideoBillingSubscriptionLoader{},
	}
}

func newRuntimeJob(status string) *VideoJob {
	snapshot, _ := json.Marshal(VideoJobBillingSnapshot{Version: 1, BillingType: BillingTypeBalance, Price720P: 0.1, RateMultiplier: 1})
	return &VideoJob{JobID: "vidjob_runtime", UserID: 1, APIKeyID: 2, GroupID: 3, AccountID: 9, UpstreamJobID: 42,
		Status: status, RequestedModel: "seedance-2.0", UpstreamModel: "seedance-2.0", Resolution: "720p", DurationSeconds: 8,
		BillingSnapshot: snapshot, HoldAmount: f64p(0.8), RequestHash: "runtime-hash"}
}

func TestVideoJobRuntimeReconcilesCompletedJobThroughSettling(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{job: newRuntimeJob(VideoJobPending)}
	client := &runtimeVideoClient{job: &LeoAsyncJob{JobID: 42, Status: VideoJobCompleted, Result: json.RawMessage(`{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_720","duration":8}}`)}}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	billing := newRuntimeBilling()
	runtime := &VideoJobRuntime{Repo: repo, Accounts: selector, Client: client, Billing: billing, PollInterval: time.Hour}

	require.NoError(t, runtime.RunOnce(context.Background()))
	require.Equal(t, VideoJobCompleted, repo.job.Status)
	require.NotNil(t, repo.job.SettledAt)
	require.GreaterOrEqual(t, len(repo.transitions), 2)
	require.Equal(t, VideoJobSettling, repo.transitions[0].to)
	require.Equal(t, VideoJobCompleted, repo.transitions[1].to)
}

func TestVideoJobRuntimeReleasesFailedAndLeavesTransientErrorsPending(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{job: newRuntimeJob(VideoJobRunning)}
	client := &runtimeVideoClient{job: &LeoAsyncJob{JobID: 42, Status: VideoJobFailed, Error: &LeoAsyncJobError{Message: "generation failed"}}}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	runtime := &VideoJobRuntime{Repo: repo, Accounts: selector, Client: client, Billing: newRuntimeBilling()}

	require.NoError(t, runtime.RunOnce(context.Background()))
	require.Equal(t, VideoJobFailed, repo.job.Status)
	require.NotNil(t, repo.job.SettledAt)

	repo = &fakeVideoJobServiceRepo{job: newRuntimeJob(VideoJobPending)}
	client = &runtimeVideoClient{err: &LeoAsyncUpstreamError{Message: "temporary", Retryable: true}}
	runtime = &VideoJobRuntime{Repo: repo, Accounts: selector, Client: client, Billing: newRuntimeBilling()}
	require.NoError(t, runtime.RunOnce(context.Background()))
	require.Equal(t, VideoJobPending, repo.job.Status)
}

func TestVideoJobRuntimeStartRecoversAfterRestartAndStopWaitsForPoll(t *testing.T) {
	repo := &fakeVideoJobServiceRepo{job: newRuntimeJob(VideoJobPending)}
	started := make(chan struct{})
	continueC := make(chan struct{})
	client := &runtimeVideoClient{job: &LeoAsyncJob{JobID: 42, Status: VideoJobPending}, started: started, continueC: continueC}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	runtime := &VideoJobRuntime{Repo: repo, Accounts: selector, Client: client, Billing: newRuntimeBilling(), PollInterval: time.Millisecond}
	runtime.Start(context.Background())
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("runtime did not start polling")
	}
	stopDone := make(chan struct{})
	go func() {
		runtime.Stop()
		close(stopDone)
	}()
	select {
	case <-stopDone:
		t.Fatal("Stop returned before in-flight poll finished")
	case <-time.After(30 * time.Millisecond):
	}
	close(continueC)
	select {
	case <-stopDone:
	case <-time.After(time.Second):
		t.Fatal("Stop did not wait for runtime shutdown")
	}
}

func TestVideoJobRuntimeMarksLocalInputsTerminalForDelayedCleanup(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)
	input, err := store.Save(bytes.NewReader([]byte{137, 80, 78, 71, 13, 10, 26, 10}))
	require.NoError(t, err)

	repo := &fakeVideoJobServiceRepo{job: newRuntimeJob(VideoJobPending)}
	repo.job.LocalInputName = input.Token
	client := &runtimeVideoClient{job: &LeoAsyncJob{JobID: 42, Status: VideoJobFailed, Error: &LeoAsyncJobError{Message: "failed"}}}
	selector := &fakeVideoJobSelector{accounts: []*Account{newVideoJobServiceTestAccount(9)}}
	runtime := &VideoJobRuntime{Repo: repo, Accounts: selector, Client: client, Billing: newRuntimeBilling(), InputStore: store}

	require.NoError(t, runtime.RunOnce(context.Background()))
	removed, err := store.Cleanup(time.Now().Add(30 * time.Minute))
	require.NoError(t, err)
	require.Equal(t, 0, removed)
	removed, err = store.Cleanup(time.Now().Add(2 * time.Hour))
	require.NoError(t, err)
	require.Equal(t, 1, removed)
}
