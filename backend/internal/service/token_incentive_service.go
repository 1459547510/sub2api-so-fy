package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
)

const (
	TokenIncentiveThresholdTokens int64   = 1_000_000_000
	TokenIncentiveRewardAmount    float64 = 10
	TokenIncentiveClaimedStatus           = "claimed"
)

var (
	ErrTokenIncentiveDisabled = infraerrors.Forbidden(
		"TOKEN_INCENTIVE_DISABLED",
		"token incentive plan is disabled",
	)
	ErrTokenIncentiveNotEligible = infraerrors.BadRequest(
		"TOKEN_INCENTIVE_NOT_ELIGIBLE",
		"weekly token usage has not reached the incentive threshold",
	)
	ErrTokenIncentiveAlreadyClaimed = infraerrors.Conflict(
		"TOKEN_INCENTIVE_ALREADY_CLAIMED",
		"token incentive reward already claimed for this week",
	)
)

type TokenIncentiveRepository interface {
	GetWeeklyUsageTokens(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, error)
	GetClaim(ctx context.Context, userID int64, weekStart time.Time) (*TokenIncentiveClaim, error)
	ClaimReward(ctx context.Context, userID int64, weekStart, weekEnd time.Time, tokens int64, rewardAmount float64) (*TokenIncentiveClaim, float64, error)
}

type TokenIncentiveClaim struct {
	ID           int64
	UserID       int64
	WeekStart    time.Time
	WeekEnd      time.Time
	Tokens       int64
	RewardAmount float64
	Status       string
	ClaimedAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TokenIncentiveStatus struct {
	Enabled         bool       `json:"enabled"`
	Eligible        bool       `json:"eligible"`
	Claimed         bool       `json:"claimed"`
	WeekStart       time.Time  `json:"week_start"`
	WeekEnd         time.Time  `json:"week_end"`
	Tokens          int64      `json:"tokens"`
	ThresholdTokens int64      `json:"threshold_tokens"`
	RewardAmount    float64    `json:"reward_amount"`
	ClaimedAt       *time.Time `json:"claimed_at,omitempty"`
	CurrentBalance  *float64   `json:"current_balance,omitempty"`
}

type TokenIncentiveService struct {
	repo                 TokenIncentiveRepository
	settingService       *SettingService
	authCacheInvalidator APIKeyAuthCacheInvalidator
	billingCache         BillingCache
}

func NewTokenIncentiveService(
	repo TokenIncentiveRepository,
	settingService *SettingService,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
	billingCache BillingCache,
) *TokenIncentiveService {
	return &TokenIncentiveService{
		repo:                 repo,
		settingService:       settingService,
		authCacheInvalidator: authCacheInvalidator,
		billingCache:         billingCache,
	}
}

func (s *TokenIncentiveService) GetStatus(ctx context.Context, userID int64) (*TokenIncentiveStatus, error) {
	weekStart, weekEnd := tokenIncentiveWeekWindow(timezone.Now())
	enabled := s.isEnabled(ctx)
	if !enabled {
		return &TokenIncentiveStatus{
			Enabled:         false,
			Eligible:        false,
			Claimed:         false,
			WeekStart:       weekStart,
			WeekEnd:         weekEnd,
			ThresholdTokens: TokenIncentiveThresholdTokens,
			RewardAmount:    TokenIncentiveRewardAmount,
		}, nil
	}

	tokens, claim, err := s.loadWeekState(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	return buildTokenIncentiveStatus(true, weekStart, weekEnd, tokens, claim, nil), nil
}

func (s *TokenIncentiveService) Claim(ctx context.Context, userID int64) (*TokenIncentiveStatus, error) {
	if !s.isEnabled(ctx) {
		return nil, ErrTokenIncentiveDisabled
	}
	weekStart, weekEnd := tokenIncentiveWeekWindow(timezone.Now())
	tokens, existingClaim, err := s.loadWeekState(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	if existingClaim != nil {
		return nil, ErrTokenIncentiveAlreadyClaimed
	}
	if tokens < TokenIncentiveThresholdTokens {
		return nil, ErrTokenIncentiveNotEligible
	}

	claim, balanceAfter, err := s.repo.ClaimReward(ctx, userID, weekStart, weekEnd, tokens, TokenIncentiveRewardAmount)
	if err != nil {
		return nil, err
	}
	s.invalidateUserCaches(ctx, userID)
	return buildTokenIncentiveStatus(true, weekStart, weekEnd, tokens, claim, &balanceAfter), nil
}

func (s *TokenIncentiveService) loadWeekState(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, *TokenIncentiveClaim, error) {
	if s == nil || s.repo == nil {
		return 0, nil, fmt.Errorf("token incentive service is not initialized")
	}
	tokens, err := s.repo.GetWeeklyUsageTokens(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return 0, nil, err
	}
	claim, err := s.repo.GetClaim(ctx, userID, weekStart)
	if err != nil {
		return 0, nil, err
	}
	return tokens, claim, nil
}

func (s *TokenIncentiveService) isEnabled(ctx context.Context) bool {
	if s == nil || s.settingService == nil {
		return false
	}
	return s.settingService.IsTokenIncentiveEnabled(ctx)
}

func (s *TokenIncentiveService) invalidateUserCaches(ctx context.Context, userID int64) {
	if s == nil {
		return
	}
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.billingCache == nil {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in token incentive balance cache invalidation", "user_id", userID, "recover", r)
			}
		}()
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.billingCache.InvalidateUserBalance(cacheCtx, userID); err != nil {
			slog.Error("invalidate token incentive user balance cache failed", "user_id", userID, "error", err)
		}
	}()
}

func tokenIncentiveWeekWindow(now time.Time) (time.Time, time.Time) {
	weekStart := timezone.StartOfWeek(now)
	return weekStart, weekStart.AddDate(0, 0, 7)
}

func buildTokenIncentiveStatus(enabled bool, weekStart, weekEnd time.Time, tokens int64, claim *TokenIncentiveClaim, currentBalance *float64) *TokenIncentiveStatus {
	status := &TokenIncentiveStatus{
		Enabled:         enabled,
		Eligible:        enabled && tokens >= TokenIncentiveThresholdTokens,
		Claimed:         claim != nil,
		WeekStart:       weekStart,
		WeekEnd:         weekEnd,
		Tokens:          tokens,
		ThresholdTokens: TokenIncentiveThresholdTokens,
		RewardAmount:    TokenIncentiveRewardAmount,
		CurrentBalance:  currentBalance,
	}
	if claim != nil {
		claimedAt := claim.ClaimedAt
		status.ClaimedAt = &claimedAt
		status.Tokens = claim.Tokens
		status.RewardAmount = claim.RewardAmount
		status.Eligible = true
	}
	return status
}
