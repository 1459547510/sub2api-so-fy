package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
)

const (
	TokenIncentiveClaimedStatus = "claimed"

	TokenIncentiveMaxRules           = 20
	TokenIncentiveMaxRewardAmount    = 10_000
	TokenIncentiveMaxThresholdTokens = math.MaxInt64 / 2
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

type TokenIncentiveRule struct {
	ThresholdTokens int64   `json:"threshold_tokens"`
	RewardAmount    float64 `json:"reward_amount"`
}

type TokenIncentiveRepository interface {
	GetWeeklyUsageTokens(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, error)
	GetClaims(ctx context.Context, userID int64, weekStart time.Time) ([]*TokenIncentiveClaim, error)
	ClaimReward(ctx context.Context, userID int64, weekStart, weekEnd time.Time, tokens int64, thresholdTokens int64, rewardAmount float64) (*TokenIncentiveClaim, float64, error)
}

type TokenIncentiveClaim struct {
	ID              int64
	UserID          int64
	WeekStart       time.Time
	WeekEnd         time.Time
	Tokens          int64
	ThresholdTokens int64
	RewardAmount    float64
	Status          string
	ClaimedAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TokenIncentiveStatus struct {
	Enabled                bool                 `json:"enabled"`
	Eligible               bool                 `json:"eligible"`
	Claimable              bool                 `json:"claimable"`
	Claimed                bool                 `json:"claimed"`
	WeekStart              time.Time            `json:"week_start"`
	WeekEnd                time.Time            `json:"week_end"`
	Tokens                 int64                `json:"tokens"`
	ThresholdTokens        int64                `json:"threshold_tokens"`
	RewardAmount           float64              `json:"reward_amount"`
	ClaimedRewardAmount    float64              `json:"claimed_reward_amount"`
	ClaimedThresholdTokens []int64              `json:"claimed_threshold_tokens,omitempty"`
	Rules                  []TokenIncentiveRule `json:"rules"`
	NextThresholdTokens    int64                `json:"next_threshold_tokens,omitempty"`
	NextRewardAmount       float64              `json:"next_reward_amount,omitempty"`
	ClaimedAt              *time.Time           `json:"claimed_at,omitempty"`
	CurrentBalance         *float64             `json:"current_balance,omitempty"`
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
	rules := s.rules(ctx)
	enabled := s.isEnabled(ctx)
	if !enabled {
		return buildTokenIncentiveStatus(false, weekStart, weekEnd, 0, nil, nil, rules), nil
	}

	tokens, claims, err := s.loadWeekState(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	return buildTokenIncentiveStatus(true, weekStart, weekEnd, tokens, claims, nil, rules), nil
}

func (s *TokenIncentiveService) Claim(ctx context.Context, userID int64) (*TokenIncentiveStatus, error) {
	if !s.isEnabled(ctx) {
		return nil, ErrTokenIncentiveDisabled
	}
	rules := s.rules(ctx)
	weekStart, weekEnd := tokenIncentiveWeekWindow(timezone.Now())
	tokens, existingClaims, err := s.loadWeekState(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	selected, _ := selectClaimableTokenIncentiveRule(tokens, rules, existingClaims)
	if selected == nil {
		if reached, _ := selectTokenIncentiveRule(tokens, rules); reached != nil {
			return nil, ErrTokenIncentiveAlreadyClaimed
		}
		return nil, ErrTokenIncentiveNotEligible
	}

	claim, balanceAfter, err := s.repo.ClaimReward(ctx, userID, weekStart, weekEnd, tokens, selected.ThresholdTokens, selected.RewardAmount)
	if err != nil {
		return nil, err
	}
	s.invalidateUserCaches(ctx, userID)
	claims := append(append([]*TokenIncentiveClaim{}, existingClaims...), claim)
	return buildTokenIncentiveStatus(true, weekStart, weekEnd, tokens, claims, &balanceAfter, rules), nil
}

func (s *TokenIncentiveService) loadWeekState(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, []*TokenIncentiveClaim, error) {
	if s == nil || s.repo == nil {
		return 0, nil, fmt.Errorf("token incentive service is not initialized")
	}
	tokens, err := s.repo.GetWeeklyUsageTokens(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return 0, nil, err
	}
	claims, err := s.repo.GetClaims(ctx, userID, weekStart)
	if err != nil {
		return 0, nil, err
	}
	return tokens, claims, nil
}

func (s *TokenIncentiveService) isEnabled(ctx context.Context) bool {
	if s == nil || s.settingService == nil {
		return false
	}
	return s.settingService.IsTokenIncentiveEnabled(ctx)
}

func (s *TokenIncentiveService) rules(ctx context.Context) []TokenIncentiveRule {
	if s == nil || s.settingService == nil {
		return DefaultTokenIncentiveRules()
	}
	return s.settingService.GetTokenIncentiveRules(ctx)
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

func buildTokenIncentiveStatus(enabled bool, weekStart, weekEnd time.Time, tokens int64, claims []*TokenIncentiveClaim, currentBalance *float64, rules []TokenIncentiveRule) *TokenIncentiveStatus {
	rules, err := NormalizeTokenIncentiveRules(rules)
	if err != nil {
		rules = DefaultTokenIncentiveRules()
	}
	claimable, next := selectClaimableTokenIncentiveRule(tokens, rules, claims)
	reached, _ := selectTokenIncentiveRule(tokens, rules)
	displayRule := tokenIncentiveDisplayRule(claimable, next, rules)
	claimedThresholds, claimedRewardAmount, latestClaimedAt := tokenIncentiveClaimSummary(claims, rules)

	status := &TokenIncentiveStatus{
		Enabled:                enabled,
		Eligible:               enabled && reached != nil,
		Claimable:              enabled && claimable != nil,
		Claimed:                len(claims) > 0,
		WeekStart:              weekStart,
		WeekEnd:                weekEnd,
		Tokens:                 tokens,
		ThresholdTokens:        displayRule.ThresholdTokens,
		RewardAmount:           displayRule.RewardAmount,
		ClaimedRewardAmount:    claimedRewardAmount,
		ClaimedThresholdTokens: claimedThresholds,
		Rules:                  rules,
		CurrentBalance:         currentBalance,
	}
	if next != nil {
		status.NextThresholdTokens = next.ThresholdTokens
		status.NextRewardAmount = next.RewardAmount
	}
	if latestClaimedAt != nil {
		claimedAt := *latestClaimedAt
		status.ClaimedAt = &claimedAt
	}
	return status
}

var defaultTokenIncentiveRules = []TokenIncentiveRule{
	{ThresholdTokens: 50_000_000, RewardAmount: 2},
	{ThresholdTokens: 100_000_000, RewardAmount: 5},
	{ThresholdTokens: 500_000_000, RewardAmount: 10},
}

func DefaultTokenIncentiveRules() []TokenIncentiveRule {
	rules := make([]TokenIncentiveRule, len(defaultTokenIncentiveRules))
	copy(rules, defaultTokenIncentiveRules)
	return rules
}

func NormalizeTokenIncentiveRules(rules []TokenIncentiveRule) ([]TokenIncentiveRule, error) {
	if len(rules) == 0 {
		return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive rules cannot be empty")
	}
	if len(rules) > TokenIncentiveMaxRules {
		return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", fmt.Sprintf("token incentive rules cannot exceed %d tiers", TokenIncentiveMaxRules))
	}
	normalized := make([]TokenIncentiveRule, 0, len(rules))
	for i, rule := range rules {
		if rule.ThresholdTokens <= 0 {
			return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive threshold must be greater than 0").WithMetadata(map[string]string{"index": fmt.Sprintf("%d", i)})
		}
		if rule.ThresholdTokens > TokenIncentiveMaxThresholdTokens {
			return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive threshold is too large").WithMetadata(map[string]string{"index": fmt.Sprintf("%d", i)})
		}
		if rule.RewardAmount <= 0 || math.IsNaN(rule.RewardAmount) || math.IsInf(rule.RewardAmount, 0) {
			return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive reward amount must be greater than 0").WithMetadata(map[string]string{"index": fmt.Sprintf("%d", i)})
		}
		if rule.RewardAmount > TokenIncentiveMaxRewardAmount {
			return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive reward amount is too large").WithMetadata(map[string]string{"index": fmt.Sprintf("%d", i)})
		}
		normalized = append(normalized, rule)
	}
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ThresholdTokens < normalized[j].ThresholdTokens
	})
	for i := 1; i < len(normalized); i++ {
		if normalized[i].ThresholdTokens == normalized[i-1].ThresholdTokens {
			return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive thresholds must be unique").WithMetadata(map[string]string{"threshold_tokens": fmt.Sprintf("%d", normalized[i].ThresholdTokens)})
		}
	}
	return normalized, nil
}

func selectTokenIncentiveRule(tokens int64, rules []TokenIncentiveRule) (*TokenIncentiveRule, *TokenIncentiveRule) {
	normalized, err := NormalizeTokenIncentiveRules(rules)
	if err != nil {
		normalized = DefaultTokenIncentiveRules()
	}
	var selected *TokenIncentiveRule
	var next *TokenIncentiveRule
	for i := range normalized {
		if tokens >= normalized[i].ThresholdTokens {
			selected = &normalized[i]
			continue
		}
		if next == nil {
			next = &normalized[i]
		}
	}
	return selected, next
}

func selectClaimableTokenIncentiveRule(tokens int64, rules []TokenIncentiveRule, claims []*TokenIncentiveClaim) (*TokenIncentiveRule, *TokenIncentiveRule) {
	normalized, err := NormalizeTokenIncentiveRules(rules)
	if err != nil {
		normalized = DefaultTokenIncentiveRules()
	}
	claimed := tokenIncentiveClaimedThresholdSet(claims, normalized)
	var claimable *TokenIncentiveRule
	var next *TokenIncentiveRule
	for i := range normalized {
		if claimed[normalized[i].ThresholdTokens] {
			continue
		}
		if tokens >= normalized[i].ThresholdTokens {
			if claimable == nil {
				claimable = &normalized[i]
			}
			continue
		}
		if next == nil {
			next = &normalized[i]
		}
	}
	return claimable, next
}

func tokenIncentiveClaimedThresholdSet(claims []*TokenIncentiveClaim, rules []TokenIncentiveRule) map[int64]bool {
	claimed := make(map[int64]bool, len(claims))
	for _, claim := range claims {
		if claim == nil {
			continue
		}
		if threshold := resolveTokenIncentiveClaimThreshold(claim, rules, claimed); threshold > 0 {
			claimed[threshold] = true
		}
	}
	return claimed
}

func resolveTokenIncentiveClaimThreshold(claim *TokenIncentiveClaim, rules []TokenIncentiveRule, claimed map[int64]bool) int64 {
	if claim == nil {
		return 0
	}
	for _, rule := range rules {
		if claim.ThresholdTokens == rule.ThresholdTokens {
			return rule.ThresholdTokens
		}
	}
	var fallback int64
	for _, rule := range rules {
		if claimed[rule.ThresholdTokens] {
			continue
		}
		if claim.Tokens >= rule.ThresholdTokens && tokenIncentiveRewardAmountEqual(claim.RewardAmount, rule.RewardAmount) {
			fallback = rule.ThresholdTokens
		}
	}
	if fallback > 0 {
		return fallback
	}
	for _, rule := range rules {
		if claimed[rule.ThresholdTokens] {
			continue
		}
		if claim.Tokens >= rule.ThresholdTokens {
			fallback = rule.ThresholdTokens
		}
	}
	return fallback
}

func tokenIncentiveRewardAmountEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.00000001
}

func tokenIncentiveClaimSummary(claims []*TokenIncentiveClaim, rules []TokenIncentiveRule) ([]int64, float64, *time.Time) {
	thresholds := make([]int64, 0, len(claims))
	claimed := make(map[int64]bool, len(claims))
	var amount float64
	var latest *time.Time
	for _, claim := range claims {
		if claim == nil {
			continue
		}
		if threshold := resolveTokenIncentiveClaimThreshold(claim, rules, claimed); threshold > 0 {
			claimed[threshold] = true
			thresholds = append(thresholds, threshold)
		}
		amount += claim.RewardAmount
		if latest == nil || claim.ClaimedAt.After(*latest) {
			claimedAt := claim.ClaimedAt
			latest = &claimedAt
		}
	}
	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i] < thresholds[j]
	})
	return thresholds, amount, latest
}

func tokenIncentiveDisplayRule(selected, next *TokenIncentiveRule, rules []TokenIncentiveRule) TokenIncentiveRule {
	if selected != nil {
		return *selected
	}
	if next != nil {
		return *next
	}
	normalized, err := NormalizeTokenIncentiveRules(rules)
	if err == nil && len(normalized) > 0 {
		return normalized[len(normalized)-1]
	}
	return DefaultTokenIncentiveRules()[0]
}
