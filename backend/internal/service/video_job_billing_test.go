package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeVideoJobBalanceRepo struct {
	UsageBillingRepository
	reserves []*VideoBalanceHoldCommand
	releases []*VideoBalanceHoldCommand
}

func (r *fakeVideoJobBalanceRepo) ReserveVideoBalance(_ context.Context, cmd *VideoBalanceHoldCommand) (*VideoBalanceHoldResult, error) {
	r.reserves = append(r.reserves, cmd)
	return &VideoBalanceHoldResult{Applied: true}, nil
}

func (r *fakeVideoJobBalanceRepo) ReleaseVideoBalance(_ context.Context, cmd *VideoBalanceHoldCommand) (*VideoBalanceHoldResult, error) {
	r.releases = append(r.releases, cmd)
	return &VideoBalanceHoldResult{Applied: true}, nil
}

type fakeVideoUsageRecorder struct {
	inputs []*OpenAIRecordUsageInput
}

func (r *fakeVideoUsageRecorder) RecordUsage(_ context.Context, input *OpenAIRecordUsageInput) error {
	r.inputs = append(r.inputs, input)
	return nil
}

type fakeVideoBillingAPIKeyLoader struct{ apiKey *APIKey }

func (l fakeVideoBillingAPIKeyLoader) GetByID(context.Context, int64) (*APIKey, error) {
	return l.apiKey, nil
}

type fakeVideoBillingUserLoader struct{ user *User }

func (l fakeVideoBillingUserLoader) GetByID(context.Context, int64) (*User, error) {
	return l.user, nil
}

type fakeVideoBillingAccountLoader struct{ account *Account }

func (l fakeVideoBillingAccountLoader) GetByID(context.Context, int64) (*Account, error) {
	return l.account, nil
}

type fakeVideoBillingSubscriptionLoader struct{ subscription *UserSubscription }

func (l fakeVideoBillingSubscriptionLoader) GetByID(context.Context, int64) (*UserSubscription, error) {
	return l.subscription, nil
}

func TestVideoJobBillingPrepareReservesSnapshotCost(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	service := &VideoJobBillingService{BillingRepo: balance}
	groupID := int64(3)
	job := &VideoJob{JobID: "vidjob_prepare", UserID: 1, APIKeyID: 2, GroupID: groupID, Resolution: "720p", DurationSeconds: 8}
	apiKey := &APIKey{ID: 2, GroupID: &groupID, Group: &Group{
		ID: groupID, RateMultiplier: 1.5,
		VideoPrice480P: f64p(0.05), VideoPrice720P: f64p(0.1), VideoPrice1080P: f64p(0.2),
	}}

	require.NoError(t, service.Prepare(context.Background(), job, apiKey, &User{ID: 1}, nil))
	require.NotNil(t, job.HoldAmount)
	require.InDelta(t, 1.2, *job.HoldAmount, 1e-12)
	require.Len(t, balance.reserves, 1)
	require.Equal(t, VideoHoldRequestID(job.JobID), balance.reserves[0].RequestID)
	require.InDelta(t, 1.2, balance.reserves[0].HoldAmount, 1e-12)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.InDelta(t, 0.1, snapshot.Price720P, 1e-12)
	require.InDelta(t, 1.5, snapshot.RateMultiplier, 1e-12)
}

func TestVideoJobSettlementIsIdempotent(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	recorder := &fakeVideoUsageRecorder{}
	snapshot, err := json.Marshal(VideoJobBillingSnapshot{
		BillingType: BillingTypeBalance, Price480P: 0.05, Price720P: 0.1, Price1080P: 0.2, RateMultiplier: 1.5,
	})
	require.NoError(t, err)
	job := &VideoJob{
		JobID: "vidjob_settle", UserID: 1, APIKeyID: 2, GroupID: 3, AccountID: 9,
		RequestedModel: "seedance-2.0", UpstreamModel: "seedance-2.0", Resolution: "720p",
		DurationSeconds: 8, BillingSnapshot: snapshot, HoldAmount: f64p(1.2), RequestHash: "request-hash",
	}
	service := &VideoJobBillingService{
		BillingRepo: balance, UsageRecorder: recorder,
		APIKeys:       fakeVideoBillingAPIKeyLoader{apiKey: &APIKey{ID: 2}},
		Users:         fakeVideoBillingUserLoader{user: &User{ID: 1}},
		Accounts:      fakeVideoBillingAccountLoader{account: &Account{ID: 9, Type: AccountTypeAPIKey}},
		Subscriptions: fakeVideoBillingSubscriptionLoader{},
	}
	result := json.RawMessage(`{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"RESOLUTION_720","duration":12}}`)

	require.NoError(t, service.SettleCompleted(context.Background(), job, result))
	require.NoError(t, service.SettleCompleted(context.Background(), job, result))

	require.Len(t, recorder.inputs, 1)
	require.NotNil(t, recorder.inputs[0].CostOverride)
	require.InDelta(t, 1.8, recorder.inputs[0].CostOverride.ActualCost, 1e-12)
	require.Equal(t, "video_usage:"+job.JobID, recorder.inputs[0].Result.RequestID)
	require.Len(t, balance.releases, 1)
	require.Equal(t, VideoReleaseRequestID(job.JobID), balance.releases[0].RequestID)
	require.NotNil(t, job.SettledAt)
	require.NotNil(t, job.ActualCost)
	require.InDelta(t, 1.8, *job.ActualCost, 1e-12)
}

func TestVideoJobFailureOnlyReleasesHold(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	recorder := &fakeVideoUsageRecorder{}
	job := &VideoJob{JobID: "vidjob_failed", UserID: 1, APIKeyID: 2, HoldAmount: f64p(0.8)}
	service := &VideoJobBillingService{BillingRepo: balance, UsageRecorder: recorder}

	require.NoError(t, service.SettleWithoutCharge(context.Background(), job))
	require.Empty(t, recorder.inputs)
	require.Len(t, balance.releases, 1)
	require.NotNil(t, job.SettledAt)
}
