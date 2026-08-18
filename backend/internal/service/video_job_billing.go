package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	videoHoldRequestPrefix    = "video_hold:"
	videoReleaseRequestPrefix = "video_release:"
	videoUsageRequestPrefix   = "video_usage:"
)

var ErrVideoInsufficientBalance = errors.New("insufficient balance for video hold")

func VideoHoldRequestID(jobID string) string {
	return videoHoldRequestPrefix + strings.TrimSpace(jobID)
}

func VideoReleaseRequestID(jobID string) string {
	return videoReleaseRequestPrefix + strings.TrimSpace(jobID)
}

type VideoJobBillingSnapshot struct {
	Version            int     `json:"version"`
	BillingType        int8    `json:"billing_type"`
	SubscriptionID     *int64  `json:"subscription_id,omitempty"`
	Price400P          float64 `json:"price_400p,omitempty"`
	Price480P          float64 `json:"price_480p"`
	Price544P          float64 `json:"price_544p,omitempty"`
	Price720P          float64 `json:"price_720p"`
	Price960P          float64 `json:"price_960p,omitempty"`
	Price1080P         float64 `json:"price_1080p"`
	Price1440P         float64 `json:"price_1440p,omitempty"`
	Price2160P         float64 `json:"price_2160p,omitempty"`
	RateMultiplier     float64 `json:"rate_multiplier"`
	ChannelID          int64   `json:"channel_id,omitempty"`
	BillingModel       string  `json:"billing_model,omitempty"`
	BillingModelSource string  `json:"billing_model_source,omitempty"`
	ChannelMappedModel string  `json:"channel_mapped_model,omitempty"`
	PricingSource      string  `json:"pricing_source,omitempty"`
}

func (s VideoJobBillingSnapshot) Cost(durationModel, resolution string, durationSeconds, videoCount int) *CostBreakdown {
	if videoCount <= 0 {
		videoCount = 1
	}
	durationSeconds = NormalizeLeoVideoBillingDurationSecondsOrDefault(durationModel, durationSeconds)
	billingResolution := NormalizeVideoBillingResolutionOrDefault(resolution)
	if s.Version >= 4 || (s.Version >= 3 && isLeoLTX23Model(durationModel)) {
		billingResolution = NormalizeLeoVideoBillingResolutionOrDefault(durationModel, resolution)
	}
	var unitPrice float64
	switch billingResolution {
	case VideoBillingResolution400P:
		unitPrice = s.Price400P
	case VideoBillingResolution544P:
		unitPrice = s.Price544P
	case VideoBillingResolution720P:
		unitPrice = s.Price720P
	case VideoBillingResolution960P:
		unitPrice = s.Price960P
	case VideoBillingResolution1080P:
		unitPrice = s.Price1080P
	case VideoBillingResolution1440P:
		unitPrice = s.Price1440P
	case VideoBillingResolution2160P:
		unitPrice = s.Price2160P
	default:
		unitPrice = s.Price480P
	}
	total := unitPrice * float64(durationSeconds) * float64(videoCount)
	return &CostBreakdown{TotalCost: total, ActualCost: total * s.RateMultiplier, BillingMode: string(BillingModeVideo)}
}

type VideoJobBillingUsageRecorder interface {
	RecordUsage(ctx context.Context, input *OpenAIRecordUsageInput) error
}

type VideoJobBillingAPIKeyLoader interface {
	GetByID(ctx context.Context, id int64) (*APIKey, error)
}

type VideoJobBillingUserLoader interface {
	GetByID(ctx context.Context, id int64) (*User, error)
}

type VideoJobBillingAccountLoader interface {
	GetByID(ctx context.Context, id int64) (*Account, error)
}

type VideoJobBillingSubscriptionLoader interface {
	GetByID(ctx context.Context, id int64) (*UserSubscription, error)
}

type VideoJobBillingService struct {
	BillingRepo   UsageBillingRepository
	UsageRecorder VideoJobBillingUsageRecorder
	APIKeyService APIKeyQuotaUpdater
	Gateway       *OpenAIGatewayService
	Pricing       *ModelPricingResolver
	APIKeys       VideoJobBillingAPIKeyLoader
	Users         VideoJobBillingUserLoader
	Accounts      VideoJobBillingAccountLoader
	Subscriptions VideoJobBillingSubscriptionLoader
}

func (s *VideoJobBillingService) Prepare(ctx context.Context, job *VideoJob, apiKey *APIKey, user *User, subscription *UserSubscription) error {
	if job == nil || apiKey == nil || apiKey.Group == nil || user == nil {
		return errors.New("video billing context is incomplete")
	}
	group := apiKey.Group
	baseMultiplier := group.RateMultiplier
	if s != nil && s.Gateway != nil {
		baseMultiplier = s.Gateway.ResolveUserGroupRateMultiplier(ctx, user.ID, group.ID, group.RateMultiplier)
	}
	mapping := ChannelMappingResult{MappedModel: job.RequestedModel}
	if s != nil && s.Gateway != nil {
		mapping = s.Gateway.ResolveChannelMapping(ctx, group.ID, job.RequestedModel)
	}
	channelMappedModel := strings.TrimSpace(mapping.MappedModel)
	if channelMappedModel == "" {
		channelMappedModel = job.RequestedModel
	}
	billingModel := resolveVideoJobBillingModel(mapping.BillingModelSource, job.RequestedModel, channelMappedModel, job.UpstreamModel)
	price480P, price720P, price1080P := group.VideoPrice480P, group.VideoPrice720P, group.VideoPrice1080P
	var price400P, price544P, price960P, price1440P, price2160P *float64
	pricingSource := "group"
	if s != nil && s.Pricing != nil {
		resolved := s.Pricing.Resolve(ctx, PricingInput{Model: billingModel, GroupID: &group.ID})
		if channelPrices, ok := VideoPriceConfigFromResolvedPricing(resolved); ok {
			price480P, price720P, price1080P = channelPrices.Price480P, channelPrices.Price720P, channelPrices.Price1080P
			price400P, price544P, price960P = channelPrices.Price400P, channelPrices.Price544P, channelPrices.Price960P
			price1440P, price2160P = channelPrices.Price1440P, channelPrices.Price2160P
			pricingSource = resolved.Source
		}
	}
	if !videoPricingIsCompleteForModel(billingModel, job.Resolution, price400P, price480P, price544P, price720P, price960P, price1080P, price1440P, price2160P) {
		return errors.New("video pricing is incomplete")
	}
	if price400P == nil {
		price400P = float64Pointer(0)
	}
	if price480P == nil {
		price480P = float64Pointer(0)
	}
	if price544P == nil {
		price544P = float64Pointer(0)
	}
	if price720P == nil {
		price720P = float64Pointer(0)
	}
	if price960P == nil {
		price960P = float64Pointer(0)
	}
	if price1080P == nil {
		price1080P = float64Pointer(0)
	}
	if price1440P == nil {
		price1440P = float64Pointer(0)
	}
	if price2160P == nil {
		price2160P = float64Pointer(0)
	}
	snapshot := VideoJobBillingSnapshot{
		Version: 4, BillingType: BillingTypeBalance,
		Price400P: *price400P,
		Price480P: *price480P, Price720P: *price720P, Price1080P: *price1080P,
		Price544P: *price544P, Price960P: *price960P,
		Price1440P: *price1440P, Price2160P: *price2160P,
		RateMultiplier: resolveVideoRateMultiplier(apiKey, baseMultiplier), ChannelID: mapping.ChannelID,
		BillingModel: billingModel, BillingModelSource: mapping.BillingModelSource,
		ChannelMappedModel: channelMappedModel, PricingSource: pricingSource,
	}
	if subscription != nil && group.IsSubscriptionType() {
		snapshot.BillingType = BillingTypeSubscription
		snapshot.SubscriptionID = &subscription.ID
	}
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	job.BillingSnapshot = raw
	cost := snapshot.Cost(job.UpstreamModel, job.Resolution, job.DurationSeconds, 1)
	job.HoldAmount = nil
	if snapshot.BillingType == BillingTypeSubscription || cost.ActualCost <= 0 {
		return nil
	}
	job.HoldAmount = float64Pointer(cost.ActualCost)
	if s == nil || s.BillingRepo == nil {
		return errors.New("video billing repository is not configured")
	}
	_, err = s.BillingRepo.ReserveVideoBalance(ctx, &VideoBalanceHoldCommand{
		RequestID: VideoHoldRequestID(job.JobID), APIKeyID: job.APIKeyID, UserID: job.UserID,
		JobID: job.JobID, HoldAmount: cost.ActualCost, RequestPayloadHash: job.RequestHash,
	})
	return err
}

func videoPricingIsCompleteForModel(model, resolution string, price400P, price480P, price544P, price720P, price960P, price1080P, price1440P, price2160P *float64) bool {
	switch NormalizeLeoVideoBillingResolutionOrDefault(model, resolution) {
	case VideoBillingResolution400P:
		return price400P != nil
	case VideoBillingResolution480P:
		return price480P != nil
	case VideoBillingResolution544P:
		return price544P != nil
	case VideoBillingResolution720P:
		return price720P != nil
	case VideoBillingResolution960P:
		return price960P != nil
	case VideoBillingResolution1080P:
		return price1080P != nil
	case VideoBillingResolution1440P:
		return price1440P != nil
	case VideoBillingResolution2160P:
		return price2160P != nil
	default:
		return price480P != nil
	}
}

func (s *VideoJobBillingService) SettleCompleted(ctx context.Context, job *VideoJob, result json.RawMessage) error {
	if job == nil {
		return ErrVideoJobNotFound
	}
	if job.SettledAt != nil {
		return nil
	}
	var snapshot VideoJobBillingSnapshot
	if err := json.Unmarshal(job.BillingSnapshot, &snapshot); err != nil {
		return fmt.Errorf("decode video billing snapshot: %w", err)
	}
	resolution := strings.TrimSpace(gjson.GetBytes(result, "provider.resolution").String())
	if resolution == "" {
		resolution = job.Resolution
	}
	duration := int(gjson.GetBytes(result, "provider.duration").Int())
	if duration <= 0 {
		duration = job.DurationSeconds
	}
	videoCount := len(gjson.GetBytes(result, "data").Array())
	if videoCount <= 0 {
		return ErrVideoOutputURLMissing
	}
	cost := snapshot.Cost(job.UpstreamModel, resolution, duration, videoCount)
	apiKey, user, account, subscription, err := s.loadUsageContext(ctx, job, snapshot)
	if err != nil {
		return err
	}
	if s.UsageRecorder == nil {
		return errors.New("video usage recorder is not configured")
	}
	usageResult := &OpenAIForwardResult{
		RequestID: videoUsageRequestPrefix + job.JobID, Model: job.RequestedModel,
		BillingModel: videoJobSnapshotBillingModel(snapshot, job), UpstreamModel: job.UpstreamModel,
		UpstreamEndpoint: "/v1/videos/generations", VideoCount: videoCount,
		VideoResolution: resolution, VideoDurationSeconds: duration,
	}
	if err := s.UsageRecorder.RecordUsage(ctx, &OpenAIRecordUsageInput{
		Result: usageResult, APIKey: apiKey, User: user, Account: account, Subscription: subscription,
		InboundEndpoint: "/v1/videos/generations", UpstreamEndpoint: "/v1/videos/generations",
		RequestPayloadHash: job.RequestHash, QuotaPlatform: account.Platform, CostOverride: cost,
		APIKeyService:      s.APIKeyService,
		ChannelUsageFields: videoJobSnapshotChannelUsageFields(snapshot, job),
	}); err != nil {
		return err
	}
	if snapshot.BillingType == BillingTypeBalance {
		if err := s.release(ctx, job); err != nil {
			return err
		}
	}
	now := time.Now()
	job.ActualCost = float64Pointer(cost.ActualCost)
	job.SettledAt = &now
	return nil
}

func (s *VideoJobBillingService) SettleWithoutCharge(ctx context.Context, job *VideoJob) error {
	if job == nil || job.SettledAt != nil {
		return nil
	}
	if err := s.release(ctx, job); err != nil {
		return err
	}
	now := time.Now()
	job.SettledAt = &now
	return nil
}

func (s *VideoJobBillingService) release(ctx context.Context, job *VideoJob) error {
	if job == nil || job.HoldAmount == nil || *job.HoldAmount <= 0 {
		return nil
	}
	if s == nil || s.BillingRepo == nil {
		return errors.New("video billing repository is not configured")
	}
	_, err := s.BillingRepo.ReleaseVideoBalance(ctx, &VideoBalanceHoldCommand{
		RequestID: VideoReleaseRequestID(job.JobID), APIKeyID: job.APIKeyID, UserID: job.UserID,
		JobID: job.JobID, HoldAmount: *job.HoldAmount, RequestPayloadHash: job.RequestHash,
	})
	return err
}

func (s *VideoJobBillingService) loadUsageContext(ctx context.Context, job *VideoJob, snapshot VideoJobBillingSnapshot) (*APIKey, *User, *Account, *UserSubscription, error) {
	if s.APIKeys == nil || s.Users == nil || s.Accounts == nil {
		return nil, nil, nil, nil, errors.New("video billing loaders are not configured")
	}
	apiKey, err := s.APIKeys.GetByID(ctx, job.APIKeyID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	user, err := s.Users.GetByID(ctx, job.UserID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	account, err := s.Accounts.GetByID(ctx, job.AccountID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	platform := PlatformLeo
	if account != nil && IsMediaPlatform(account.Platform) {
		platform = account.Platform
	}
	group := &Group{ID: job.GroupID, Platform: platform, RateMultiplier: snapshot.RateMultiplier,
		VideoPrice480P: float64Pointer(snapshot.Price480P), VideoPrice720P: float64Pointer(snapshot.Price720P), VideoPrice1080P: float64Pointer(snapshot.Price1080P)}
	var subscription *UserSubscription
	if snapshot.BillingType == BillingTypeSubscription {
		if snapshot.SubscriptionID == nil || s.Subscriptions == nil {
			return nil, nil, nil, nil, errors.New("video subscription snapshot is incomplete")
		}
		group.SubscriptionType = SubscriptionTypeSubscription
		subscription, err = s.Subscriptions.GetByID(ctx, *snapshot.SubscriptionID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	apiKeyCopy := *apiKey
	apiKeyCopy.GroupID = int64Pointer(job.GroupID)
	apiKeyCopy.Group = group
	apiKeyCopy.User = user
	return &apiKeyCopy, user, account, subscription, nil
}

func resolveVideoJobBillingModel(source, requestedModel, channelMappedModel, upstreamModel string) string {
	requestedModel = strings.TrimSpace(requestedModel)
	channelMappedModel = strings.TrimSpace(channelMappedModel)
	upstreamModel = strings.TrimSpace(upstreamModel)
	switch source {
	case BillingModelSourceRequested:
		return requestedModel
	case BillingModelSourceUpstream:
		if upstreamModel != "" {
			return upstreamModel
		}
	case BillingModelSourceChannelMapped:
		if channelMappedModel != "" {
			return channelMappedModel
		}
	default:
		if channelMappedModel != "" {
			return channelMappedModel
		}
	}
	return requestedModel
}

func videoJobSnapshotBillingModel(snapshot VideoJobBillingSnapshot, job *VideoJob) string {
	if billingModel := strings.TrimSpace(snapshot.BillingModel); billingModel != "" {
		return billingModel
	}
	return strings.TrimSpace(job.RequestedModel)
}

func videoJobSnapshotChannelUsageFields(snapshot VideoJobBillingSnapshot, job *VideoJob) ChannelUsageFields {
	if snapshot.ChannelID <= 0 {
		return ChannelUsageFields{}
	}
	mappedModel := strings.TrimSpace(snapshot.ChannelMappedModel)
	if mappedModel == "" {
		mappedModel = job.RequestedModel
	}
	mapping := ChannelMappingResult{
		MappedModel: mappedModel, ChannelID: snapshot.ChannelID,
		Mapped: mappedModel != job.RequestedModel, BillingModelSource: snapshot.BillingModelSource,
	}
	return mapping.ToUsageFields(job.RequestedModel, job.UpstreamModel)
}

func float64Pointer(value float64) *float64 { return &value }
func int64Pointer(value int64) *int64       { return &value }
