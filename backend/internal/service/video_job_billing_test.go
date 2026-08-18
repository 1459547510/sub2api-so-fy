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

type fakeVideoAPIKeyQuotaUpdater struct{}

func (fakeVideoAPIKeyQuotaUpdater) UpdateQuotaUsed(context.Context, int64, float64) error {
	return nil
}

func (fakeVideoAPIKeyQuotaUpdater) UpdateRateLimitUsage(context.Context, int64, float64) error {
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

type fakeVideoBillingChannelRepository struct {
	ChannelRepository
	channels       []Channel
	groupPlatforms map[int64]string
}

func (r *fakeVideoBillingChannelRepository) ListAll(context.Context) ([]Channel, error) {
	return r.channels, nil
}

func (r *fakeVideoBillingChannelRepository) GetGroupPlatforms(context.Context, []int64) (map[int64]string, error) {
	return r.groupPlatforms, nil
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
	require.Equal(t, "group", snapshot.PricingSource)
}

func TestVideoJobBillingPrepareUsesChannelVideoPricingByBillingModelSource(t *testing.T) {
	tests := []struct {
		name               string
		billingModelSource string
		pricingModel       string
		wantBillingModel   string
	}{
		{name: "requested", billingModelSource: BillingModelSourceRequested, pricingModel: "seedance-public", wantBillingModel: "seedance-public"},
		{name: "channel_mapped", billingModelSource: BillingModelSourceChannelMapped, pricingModel: "seedance-channel", wantBillingModel: "seedance-channel"},
		{name: "upstream", billingModelSource: BillingModelSourceUpstream, pricingModel: "seedance-upstream", wantBillingModel: "seedance-upstream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance := &fakeVideoJobBalanceRepo{}
			gateway, pricing := newVideoJobChannelPricingServices(t, tt.billingModelSource, tt.pricingModel, 0.01, 0.02, 0.03)
			svc := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
			groupID := int64(100)
			job := &VideoJob{
				JobID: "vidjob_channel_" + tt.name, UserID: 1, APIKeyID: 2, GroupID: groupID,
				RequestedModel: "seedance-public", UpstreamModel: "seedance-upstream", Resolution: "720p", DurationSeconds: 5,
			}
			apiKey := newVideoJobBillingAPIKey(groupID)

			require.NoError(t, svc.Prepare(context.Background(), job, apiKey, &User{ID: 1}, nil))
			require.NotNil(t, job.HoldAmount)
			require.InDelta(t, 0.1, *job.HoldAmount, 1e-12)
			require.Len(t, balance.reserves, 1)

			var snapshot VideoJobBillingSnapshot
			require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
			require.Equal(t, 4, snapshot.Version)
			require.Equal(t, int64(77), snapshot.ChannelID)
			require.Equal(t, tt.wantBillingModel, snapshot.BillingModel)
			require.Equal(t, tt.billingModelSource, snapshot.BillingModelSource)
			require.Equal(t, "seedance-channel", snapshot.ChannelMappedModel)
			require.Equal(t, PricingSourceChannel, snapshot.PricingSource)
			require.InDelta(t, 0.02, snapshot.Price720P, 1e-12)
		})
	}
}

func TestVideoJobBillingPrepareAcceptsLegacySeedance20ThreeTierPricing(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	gateway, pricing := newVideoJobChannelPricingServices(t, BillingModelSourceRequested, "seedance-2.0", 0.12, 0.25, 0.60)
	svc := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_seedance20_1080p", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "seedance-2.0", UpstreamModel: "seedance-2.0", Resolution: "1080p", DurationSeconds: 15,
	}

	require.NoError(t, svc.Prepare(context.Background(), job, newVideoJobBillingAPIKey(groupID), &User{ID: 1}, nil))
	require.NotNil(t, job.HoldAmount)
	require.InDelta(t, 9.0, *job.HoldAmount, 1e-12)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.Equal(t, PricingSourceChannel, snapshot.PricingSource)
	require.InDelta(t, 0.60, snapshot.Price1080P, 1e-12)
	require.Zero(t, snapshot.Price2160P)
}

func TestVideoJobBillingPrepareRejectsSeedance20FourKWithout2160Price(t *testing.T) {
	gateway, pricing := newVideoJobChannelPricingServices(t, BillingModelSourceRequested, "seedance-2.0", 0.12, 0.25, 0.60)
	svc := &VideoJobBillingService{BillingRepo: &fakeVideoJobBalanceRepo{}, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_seedance20_4k", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "seedance-2.0", UpstreamModel: "seedance-2.0", Resolution: "4k", DurationSeconds: 15,
	}

	err := svc.Prepare(context.Background(), job, newVideoJobBillingAPIKey(groupID), &User{ID: 1}, nil)
	require.EqualError(t, err, "video pricing is incomplete")
}

func TestVideoJobBillingPrepareFallsBackToGroupVideoPricing(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	gateway, pricing := newVideoJobChannelPricingServices(t, BillingModelSourceRequested, "another-model", 0.01, 0.02, 0.03)
	svc := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_group_fallback", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "seedance-public", UpstreamModel: "seedance-upstream", Resolution: "720p", DurationSeconds: 5,
	}

	require.NoError(t, svc.Prepare(context.Background(), job, newVideoJobBillingAPIKey(groupID), &User{ID: 1}, nil))
	require.NotNil(t, job.HoldAmount)
	require.InDelta(t, 0.5, *job.HoldAmount, 1e-12)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.Equal(t, "group", snapshot.PricingSource)
	require.InDelta(t, 0.1, snapshot.Price720P, 1e-12)
	require.Equal(t, int64(77), snapshot.ChannelID)
}

func TestVideoJobBillingPreparePreservesExplicitZeroChannelPrice(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	gateway, pricing := newVideoJobChannelPricingServices(t, BillingModelSourceRequested, "seedance-public", 0, 0, 0)
	svc := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_channel_zero", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "seedance-public", UpstreamModel: "seedance-upstream", Resolution: "720p", DurationSeconds: 5,
	}

	require.NoError(t, svc.Prepare(context.Background(), job, newVideoJobBillingAPIKey(groupID), &User{ID: 1}, nil))
	require.Nil(t, job.HoldAmount)
	require.Empty(t, balance.reserves)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.Equal(t, PricingSourceChannel, snapshot.PricingSource)
	require.Zero(t, snapshot.Price480P)
	require.Zero(t, snapshot.Price720P)
	require.Zero(t, snapshot.Price1080P)
}

func TestVideoJobBillingPrepareLTXUsesNativeResolutionPrice(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	gateway, pricing := newLTXVideoJobChannelPricingServices(t, "ltx-2.3-fast", 0.06, 0.21, 0.24)
	service := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_ltx_prepare", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "ltx-2.3-fast", UpstreamModel: "ltxv-2.3-fast",
		Resolution: "2160p", DurationSeconds: 20,
	}
	apiKey := newVideoJobBillingAPIKey(groupID)

	require.NoError(t, service.Prepare(context.Background(), job, apiKey, &User{ID: 1}, nil))
	require.NotNil(t, job.HoldAmount)
	require.InDelta(t, 4.8, *job.HoldAmount, 1e-12)
	require.Len(t, balance.reserves, 1)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.Zero(t, snapshot.Price480P)
	require.Zero(t, snapshot.Price720P)
	require.InDelta(t, 0.06, snapshot.Price1080P, 1e-12)
	require.InDelta(t, 0.21, snapshot.Price1440P, 1e-12)
	require.InDelta(t, 0.24, snapshot.Price2160P, 1e-12)
}

func TestVideoJobBillingPrepareGrokUsesNativeResolutionPrice(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	gateway, pricing := newNativeVideoJobChannelPricingServices(t, "grok-imagine-1.5", []PricingInterval{
		{TierLabel: "400p", PerRequestPrice: f64p(0.10)},
		{TierLabel: "544p", PerRequestPrice: f64p(0.10)},
		{TierLabel: "720p", PerRequestPrice: f64p(0.17)},
		{TierLabel: "960p", PerRequestPrice: f64p(0.17)},
	})
	service := &VideoJobBillingService{BillingRepo: balance, Gateway: gateway, Pricing: pricing}
	groupID := int64(100)
	job := &VideoJob{
		JobID: "vidjob_grok_prepare", UserID: 1, APIKeyID: 2, GroupID: groupID,
		RequestedModel: "grok-imagine-1.5", UpstreamModel: "grok-imagine-1.5",
		Resolution: "544p", DurationSeconds: 6,
	}

	require.NoError(t, service.Prepare(context.Background(), job, newVideoJobBillingAPIKey(groupID), &User{ID: 1}, nil))
	require.NotNil(t, job.HoldAmount)
	require.InDelta(t, 0.60, *job.HoldAmount, 1e-12)

	var snapshot VideoJobBillingSnapshot
	require.NoError(t, json.Unmarshal(job.BillingSnapshot, &snapshot))
	require.Equal(t, 4, snapshot.Version)
	require.InDelta(t, 0.10, snapshot.Price400P, 1e-12)
	require.InDelta(t, 0.10, snapshot.Price544P, 1e-12)
	require.InDelta(t, 0.17, snapshot.Price720P, 1e-12)
	require.InDelta(t, 0.17, snapshot.Price960P, 1e-12)
	require.Zero(t, snapshot.Price480P)
	require.Zero(t, snapshot.Price1080P)
}

func TestVideoJobSettlementIsIdempotent(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	recorder := &fakeVideoUsageRecorder{}
	quotaUpdater := &fakeVideoAPIKeyQuotaUpdater{}
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
		APIKeyService: quotaUpdater,
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
	require.Same(t, quotaUpdater, recorder.inputs[0].APIKeyService)
	require.Len(t, balance.releases, 1)
	require.Equal(t, VideoReleaseRequestID(job.JobID), balance.releases[0].RequestID)
	require.NotNil(t, job.SettledAt)
	require.NotNil(t, job.ActualCost)
	require.InDelta(t, 1.8, *job.ActualCost, 1e-12)
}

func TestVideoJobSettlementLTXUsesNativeResolutionPrice(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	recorder := &fakeVideoUsageRecorder{}
	snapshot, err := json.Marshal(VideoJobBillingSnapshot{
		Version: 3, BillingType: BillingTypeBalance,
		Price1080P: 0.06, Price1440P: 0.21, Price2160P: 0.24, RateMultiplier: 1,
		BillingModel: "ltx-2.3-fast",
	})
	require.NoError(t, err)
	job := &VideoJob{
		JobID: "vidjob_ltx_settle", UserID: 1, APIKeyID: 2, GroupID: 3, AccountID: 9,
		RequestedModel: "ltx-2.3-fast", UpstreamModel: "ltxv-2.3-fast", Resolution: "1440p",
		DurationSeconds: 6, BillingSnapshot: snapshot, HoldAmount: f64p(1.2), RequestHash: "request-hash",
	}
	service := &VideoJobBillingService{
		BillingRepo: balance, UsageRecorder: recorder,
		APIKeyService: fakeVideoAPIKeyQuotaUpdater{},
		APIKeys:       fakeVideoBillingAPIKeyLoader{apiKey: &APIKey{ID: 2}},
		Users:         fakeVideoBillingUserLoader{user: &User{ID: 1}},
		Accounts:      fakeVideoBillingAccountLoader{account: &Account{ID: 9, Type: AccountTypeAPIKey}},
		Subscriptions: fakeVideoBillingSubscriptionLoader{},
	}

	result := json.RawMessage(`{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"2160p","duration":20}}`)
	require.NoError(t, service.SettleCompleted(context.Background(), job, result))
	require.Len(t, recorder.inputs, 1)
	require.InDelta(t, 4.8, recorder.inputs[0].CostOverride.ActualCost, 1e-12)
	require.Len(t, balance.releases, 1)
}

func TestVideoJobSettlementGrokUsesNativeResolutionPrice(t *testing.T) {
	balance := &fakeVideoJobBalanceRepo{}
	recorder := &fakeVideoUsageRecorder{}
	snapshot, err := json.Marshal(VideoJobBillingSnapshot{
		Version: 4, BillingType: BillingTypeBalance,
		Price400P: 0.10, Price544P: 0.10, Price720P: 0.17, Price960P: 0.17, RateMultiplier: 1,
		BillingModel: "grok-imagine-1.5",
	})
	require.NoError(t, err)
	job := &VideoJob{
		JobID: "vidjob_grok_settle", UserID: 1, APIKeyID: 2, GroupID: 3, AccountID: 9,
		RequestedModel: "grok-imagine-1.5", UpstreamModel: "grok-imagine-1.5", Resolution: "720p",
		DurationSeconds: 6, BillingSnapshot: snapshot, HoldAmount: f64p(1.02), RequestHash: "request-hash",
	}
	service := &VideoJobBillingService{
		BillingRepo: balance, UsageRecorder: recorder,
		APIKeyService: fakeVideoAPIKeyQuotaUpdater{},
		APIKeys:       fakeVideoBillingAPIKeyLoader{apiKey: &APIKey{ID: 2}},
		Users:         fakeVideoBillingUserLoader{user: &User{ID: 1}},
		Accounts:      fakeVideoBillingAccountLoader{account: &Account{ID: 9, Type: AccountTypeAPIKey}},
	}

	result := json.RawMessage(`{"data":[{"url":"https://cdn.example/video.mp4"}],"provider":{"resolution":"960p","duration":10}}`)
	require.NoError(t, service.SettleCompleted(context.Background(), job, result))
	require.Len(t, recorder.inputs, 1)
	require.InDelta(t, 1.70, recorder.inputs[0].CostOverride.ActualCost, 1e-12)
}

func TestVideoJobBillingSnapshotV3GrokKeepsCompatibilityPrices(t *testing.T) {
	snapshot := VideoJobBillingSnapshot{
		Version: 3, Price480P: 0.10, Price720P: 0.17, Price1080P: 0.17, RateMultiplier: 1,
	}

	require.InDelta(t, 0.60, snapshot.Cost("grok-imagine-1.5", "544p", 6, 1).ActualCost, 1e-12)
	require.InDelta(t, 1.02, snapshot.Cost("grok-imagine-1.5", "960p", 6, 1).ActualCost, 1e-12)
}

func TestVideoJobBillingSnapshotV2LTXUses1080PCompatibilityPrice(t *testing.T) {
	snapshot := VideoJobBillingSnapshot{
		Version: 2, Price1080P: 0.2, RateMultiplier: 1,
	}

	cost := snapshot.Cost("ltxv-2.3-fast", "2160p", 20, 1)

	require.InDelta(t, 4, cost.ActualCost, 1e-12)
}

func TestVideoJobSettlementRecordsFrozenChannelUsageFields(t *testing.T) {
	recorder := &fakeVideoUsageRecorder{}
	snapshot, err := json.Marshal(VideoJobBillingSnapshot{
		Version: 2, BillingType: BillingTypeBalance, Price480P: 0.01, Price720P: 0.02, Price1080P: 0.03, RateMultiplier: 1,
		ChannelID: 77, BillingModel: "seedance-channel", BillingModelSource: BillingModelSourceChannelMapped,
		ChannelMappedModel: "seedance-channel", PricingSource: PricingSourceChannel,
	})
	require.NoError(t, err)
	job := &VideoJob{
		JobID: "vidjob_channel_settle", UserID: 1, APIKeyID: 2, GroupID: 100, AccountID: 9,
		RequestedModel: "seedance-public", UpstreamModel: "seedance-upstream", Resolution: "720p",
		DurationSeconds: 5, BillingSnapshot: snapshot, RequestHash: "request-hash",
	}
	svc := &VideoJobBillingService{
		UsageRecorder: recorder,
		APIKeys:       fakeVideoBillingAPIKeyLoader{apiKey: &APIKey{ID: 2}},
		Users:         fakeVideoBillingUserLoader{user: &User{ID: 1}},
		Accounts:      fakeVideoBillingAccountLoader{account: &Account{ID: 9, Type: AccountTypeAPIKey}},
		Subscriptions: fakeVideoBillingSubscriptionLoader{},
	}

	require.NoError(t, svc.SettleCompleted(context.Background(), job, json.RawMessage(`{"data":[{"url":"https://cdn.example/video.mp4"}]}`)))
	require.Len(t, recorder.inputs, 1)
	input := recorder.inputs[0]
	require.Equal(t, "seedance-channel", input.Result.BillingModel)
	require.Equal(t, int64(77), input.ChannelID)
	require.Equal(t, "seedance-public", input.OriginalModel)
	require.Equal(t, "seedance-channel", input.ChannelMappedModel)
	require.Equal(t, BillingModelSourceChannelMapped, input.BillingModelSource)
	require.NotEmpty(t, input.ModelMappingChain)
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

func TestVideoJobSettlementRejectsEmptyResult(t *testing.T) {
	recorder := &fakeVideoUsageRecorder{}
	snapshot, err := json.Marshal(VideoJobBillingSnapshot{BillingType: BillingTypeBalance, Price720P: 0.1, RateMultiplier: 1})
	require.NoError(t, err)
	job := &VideoJob{JobID: "vidjob_empty", BillingSnapshot: snapshot}
	service := &VideoJobBillingService{UsageRecorder: recorder}

	err = service.SettleCompleted(context.Background(), job, json.RawMessage(`{"data":[]}`))

	require.ErrorIs(t, err, ErrVideoOutputURLMissing)
	require.Empty(t, recorder.inputs)
	require.Nil(t, job.SettledAt)
}

func newVideoJobBillingAPIKey(groupID int64) *APIKey {
	return &APIKey{ID: 2, GroupID: &groupID, Group: &Group{
		ID: groupID, Platform: PlatformLeo, RateMultiplier: 1,
		VideoPrice480P: f64p(0.05), VideoPrice720P: f64p(0.1), VideoPrice1080P: f64p(0.2),
	}}
}

func newVideoJobChannelPricingServices(
	t *testing.T,
	billingModelSource string,
	pricingModel string,
	price480P, price720P, price1080P float64,
) (*OpenAIGatewayService, *ModelPricingResolver) {
	t.Helper()
	const groupID int64 = 100
	repo := &fakeVideoBillingChannelRepository{
		channels: []Channel{{
			ID: 77, Name: "leo-video", Status: StatusActive, GroupIDs: []int64{groupID},
			BillingModelSource: billingModelSource,
			ModelMapping:       map[string]map[string]string{PlatformLeo: {"seedance-public": "seedance-channel"}},
			ModelPricing: []ChannelModelPricing{{
				Platform: PlatformLeo, Models: []string{pricingModel}, BillingMode: BillingModeVideo,
				Intervals: []PricingInterval{
					{TierLabel: "480p", PerRequestPrice: f64p(price480P)},
					{TierLabel: "720p", PerRequestPrice: f64p(price720P)},
					{TierLabel: "1080p", PerRequestPrice: f64p(price1080P)},
				},
			}},
		}},
		groupPlatforms: map[int64]string{groupID: PlatformLeo},
	}
	channelService := NewChannelService(repo, nil, nil, nil)
	resolver := NewModelPricingResolver(channelService, &BillingService{fallbackPrices: map[string]*ModelPricing{}})
	return &OpenAIGatewayService{channelService: channelService}, resolver
}

func newLTXVideoJobChannelPricingServices(
	t *testing.T,
	model string,
	price1080P, price1440P, price2160P float64,
) (*OpenAIGatewayService, *ModelPricingResolver) {
	t.Helper()
	const groupID int64 = 100
	repo := &fakeVideoBillingChannelRepository{
		channels: []Channel{{
			ID: 77, Name: "leo-video", Status: StatusActive, GroupIDs: []int64{groupID},
			BillingModelSource: BillingModelSourceRequested,
			ModelPricing: []ChannelModelPricing{{
				Platform: PlatformLeo, Models: []string{model}, BillingMode: BillingModeVideo,
				Intervals: []PricingInterval{
					{TierLabel: "1080p", PerRequestPrice: f64p(price1080P)},
					{TierLabel: "1440p", PerRequestPrice: f64p(price1440P)},
					{TierLabel: "2160p", PerRequestPrice: f64p(price2160P)},
				},
			}},
		}},
		groupPlatforms: map[int64]string{groupID: PlatformLeo},
	}
	channelService := NewChannelService(repo, nil, nil, nil)
	resolver := NewModelPricingResolver(channelService, &BillingService{fallbackPrices: map[string]*ModelPricing{}})
	return &OpenAIGatewayService{channelService: channelService}, resolver
}

func newNativeVideoJobChannelPricingServices(
	t *testing.T,
	model string,
	intervals []PricingInterval,
) (*OpenAIGatewayService, *ModelPricingResolver) {
	t.Helper()
	const groupID int64 = 100
	repo := &fakeVideoBillingChannelRepository{
		channels: []Channel{{
			ID: 77, Name: "leo-video", Status: StatusActive, GroupIDs: []int64{groupID},
			BillingModelSource: BillingModelSourceRequested,
			ModelPricing: []ChannelModelPricing{{
				Platform: PlatformLeo, Models: []string{model}, BillingMode: BillingModeVideo,
				Intervals: intervals,
			}},
		}},
		groupPlatforms: map[int64]string{groupID: PlatformLeo},
	}
	channelService := NewChannelService(repo, nil, nil, nil)
	resolver := NewModelPricingResolver(channelService, &BillingService{fallbackPrices: map[string]*ModelPricing{}})
	return &OpenAIGatewayService{channelService: channelService}, resolver
}
