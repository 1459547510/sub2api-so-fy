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
	Version        int     `json:"version"`
	BillingType    int8    `json:"billing_type"`
	SubscriptionID *int64  `json:"subscription_id,omitempty"`
	Price480P      float64 `json:"price_480p"`
	Price720P      float64 `json:"price_720p"`
	Price1080P     float64 `json:"price_1080p"`
	RateMultiplier float64 `json:"rate_multiplier"`
}

func (s VideoJobBillingSnapshot) Cost(resolution string, durationSeconds, videoCount int) *CostBreakdown {
	if videoCount <= 0 {
		videoCount = 1
	}
	durationSeconds = NormalizeVideoBillingDurationSecondsOrDefault(durationSeconds)
	var unitPrice float64
	switch NormalizeVideoBillingResolutionOrDefault(resolution) {
	case VideoBillingResolution720P:
		unitPrice = s.Price720P
	case VideoBillingResolution1080P:
		unitPrice = s.Price1080P
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
	Gateway       *OpenAIGatewayService
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
	if group.VideoPrice480P == nil || group.VideoPrice720P == nil || group.VideoPrice1080P == nil {
		return errors.New("video pricing is incomplete")
	}
	baseMultiplier := group.RateMultiplier
	if s != nil && s.Gateway != nil {
		baseMultiplier = s.Gateway.ResolveUserGroupRateMultiplier(ctx, user.ID, group.ID, group.RateMultiplier)
	}
	snapshot := VideoJobBillingSnapshot{
		Version: 1, BillingType: BillingTypeBalance,
		Price480P: *group.VideoPrice480P, Price720P: *group.VideoPrice720P, Price1080P: *group.VideoPrice1080P,
		RateMultiplier: resolveVideoRateMultiplier(apiKey, baseMultiplier),
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
	cost := snapshot.Cost(job.Resolution, job.DurationSeconds, 1)
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
	cost := snapshot.Cost(resolution, duration, videoCount)
	apiKey, user, account, subscription, err := s.loadUsageContext(ctx, job, snapshot)
	if err != nil {
		return err
	}
	if s.UsageRecorder == nil {
		return errors.New("video usage recorder is not configured")
	}
	usageResult := &OpenAIForwardResult{
		RequestID: videoUsageRequestPrefix + job.JobID, Model: job.RequestedModel,
		BillingModel: job.RequestedModel, UpstreamModel: job.UpstreamModel,
		UpstreamEndpoint: "/v1/videos/generations", VideoCount: videoCount,
		VideoResolution: resolution, VideoDurationSeconds: duration,
	}
	if err := s.UsageRecorder.RecordUsage(ctx, &OpenAIRecordUsageInput{
		Result: usageResult, APIKey: apiKey, User: user, Account: account, Subscription: subscription,
		InboundEndpoint: "/v1/videos/generations", UpstreamEndpoint: "/v1/videos/generations",
		RequestPayloadHash: job.RequestHash, QuotaPlatform: PlatformLeo, CostOverride: cost,
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
	group := &Group{ID: job.GroupID, Platform: PlatformLeo, RateMultiplier: snapshot.RateMultiplier,
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

func float64Pointer(value float64) *float64 { return &value }
func int64Pointer(value int64) *int64       { return &value }
