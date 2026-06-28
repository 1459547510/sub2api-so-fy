//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type tokenIncentiveRepoStub struct {
	tokens        int64
	claims        []*TokenIncentiveClaim
	claimErr      error
	weeklyErr     error
	capturedClaim tokenIncentiveCapturedClaim
}

type tokenIncentiveCapturedClaim struct {
	userID          int64
	weekStart       time.Time
	weekEnd         time.Time
	tokens          int64
	thresholdTokens int64
	rewardAmount    float64
	called          bool
}

func (r *tokenIncentiveRepoStub) GetWeeklyUsageTokens(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, error) {
	if r.weeklyErr != nil {
		return 0, r.weeklyErr
	}
	return r.tokens, nil
}

func (r *tokenIncentiveRepoStub) GetClaims(ctx context.Context, userID int64, weekStart time.Time) ([]*TokenIncentiveClaim, error) {
	if r.claimErr != nil {
		return nil, r.claimErr
	}
	return r.claims, nil
}

func (r *tokenIncentiveRepoStub) ClaimReward(ctx context.Context, userID int64, weekStart, weekEnd time.Time, tokens int64, thresholdTokens int64, rewardAmount float64) (*TokenIncentiveClaim, float64, error) {
	r.capturedClaim = tokenIncentiveCapturedClaim{
		userID:          userID,
		weekStart:       weekStart,
		weekEnd:         weekEnd,
		tokens:          tokens,
		thresholdTokens: thresholdTokens,
		rewardAmount:    rewardAmount,
		called:          true,
	}
	return &TokenIncentiveClaim{
		ID:              7,
		UserID:          userID,
		WeekStart:       weekStart,
		WeekEnd:         weekEnd,
		Tokens:          tokens,
		ThresholdTokens: thresholdTokens,
		RewardAmount:    rewardAmount,
		Status:          TokenIncentiveClaimedStatus,
		ClaimedAt:       time.Now(),
	}, 12.5, nil
}

type tokenIncentiveSettingRepoStub struct {
	values map[string]string
}

func (s *tokenIncentiveSettingRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *tokenIncentiveSettingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if value, ok := s.values[key]; ok {
		return value, nil
	}
	return "", ErrSettingNotFound
}

func (s *tokenIncentiveSettingRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *tokenIncentiveSettingRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
}

func (s *tokenIncentiveSettingRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
}

func (s *tokenIncentiveSettingRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
}

func (s *tokenIncentiveSettingRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

func newTokenIncentiveTestService(repo *tokenIncentiveRepoStub, rules string, enabled bool) *TokenIncentiveService {
	values := map[string]string{
		SettingKeyTokenIncentiveRules: rules,
	}
	if enabled {
		values[SettingKeyTokenIncentiveEnabled] = "true"
	} else {
		values[SettingKeyTokenIncentiveEnabled] = "false"
	}
	return NewTokenIncentiveService(
		repo,
		NewSettingService(&tokenIncentiveSettingRepoStub{values: values}, nil),
		nil,
		nil,
	)
}

func TestNormalizeTokenIncentiveRules_SortsAndValidates(t *testing.T) {
	got, err := NormalizeTokenIncentiveRules([]TokenIncentiveRule{
		{ThresholdTokens: 500_000_000, RewardAmount: 10},
		{ThresholdTokens: 50_000_000, RewardAmount: 2},
		{ThresholdTokens: 100_000_000, RewardAmount: 5},
	})
	require.NoError(t, err)
	require.Equal(t, []TokenIncentiveRule{
		{ThresholdTokens: 50_000_000, RewardAmount: 2},
		{ThresholdTokens: 100_000_000, RewardAmount: 5},
		{ThresholdTokens: 500_000_000, RewardAmount: 10},
	}, got)

	_, err = NormalizeTokenIncentiveRules([]TokenIncentiveRule{{ThresholdTokens: 0, RewardAmount: 2}})
	require.Error(t, err)

	_, err = NormalizeTokenIncentiveRules([]TokenIncentiveRule{
		{ThresholdTokens: 50_000_000, RewardAmount: 2},
		{ThresholdTokens: 50_000_000, RewardAmount: 3},
	})
	require.Error(t, err)

	_, err = NormalizeTokenIncentiveRules([]TokenIncentiveRule{{ThresholdTokens: 1, RewardAmount: TokenIncentiveMaxRewardAmount + 0.01}})
	require.Error(t, err)

	_, err = NormalizeTokenIncentiveRules([]TokenIncentiveRule{{ThresholdTokens: TokenIncentiveMaxThresholdTokens + 1, RewardAmount: 1}})
	require.Error(t, err)

	tooMany := make([]TokenIncentiveRule, TokenIncentiveMaxRules+1)
	for i := range tooMany {
		tooMany[i] = TokenIncentiveRule{ThresholdTokens: int64(i + 1), RewardAmount: 1}
	}
	_, err = NormalizeTokenIncentiveRules(tooMany)
	require.Error(t, err)
}

func TestSelectTokenIncentiveRule_SelectsHighestReachedTierAndNext(t *testing.T) {
	rules := DefaultTokenIncentiveRules()

	selected, next := selectTokenIncentiveRule(60_000_000, rules)
	require.NotNil(t, selected)
	require.EqualValues(t, 50_000_000, selected.ThresholdTokens)
	require.NotNil(t, next)
	require.EqualValues(t, 100_000_000, next.ThresholdTokens)

	selected, next = selectTokenIncentiveRule(600_000_000, rules)
	require.NotNil(t, selected)
	require.EqualValues(t, 500_000_000, selected.ThresholdTokens)
	require.Nil(t, next)
}

func TestTokenIncentiveServiceClaim_UsesFirstReachedUnclaimedTier(t *testing.T) {
	repo := &tokenIncentiveRepoStub{tokens: 500_000_000}
	svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2},{"threshold_tokens":100000000,"reward_amount":5},{"threshold_tokens":500000000,"reward_amount":10}]`, true)

	status, err := svc.Claim(context.Background(), 42)

	require.NoError(t, err)
	require.True(t, repo.capturedClaim.called)
	require.EqualValues(t, 42, repo.capturedClaim.userID)
	require.EqualValues(t, 500_000_000, repo.capturedClaim.tokens)
	require.EqualValues(t, 50_000_000, repo.capturedClaim.thresholdTokens)
	require.Equal(t, 2.0, repo.capturedClaim.rewardAmount)
	require.True(t, status.Claimed)
	require.True(t, status.Claimable)
	require.Equal(t, 5.0, status.RewardAmount)
	require.NotNil(t, status.CurrentBalance)
	require.Equal(t, 12.5, *status.CurrentBalance)
}

func TestTokenIncentiveServiceClaim_AllowsNextTierFullReward(t *testing.T) {
	repo := &tokenIncentiveRepoStub{
		tokens: 100_000_000,
		claims: []*TokenIncentiveClaim{
			{ID: 1, UserID: 42, Tokens: 50_000_000, ThresholdTokens: 50_000_000, RewardAmount: 2, ClaimedAt: time.Now()},
		},
	}
	svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2},{"threshold_tokens":100000000,"reward_amount":5},{"threshold_tokens":500000000,"reward_amount":10}]`, true)

	status, err := svc.Claim(context.Background(), 42)

	require.NoError(t, err)
	require.True(t, repo.capturedClaim.called)
	require.EqualValues(t, 100_000_000, repo.capturedClaim.thresholdTokens)
	require.Equal(t, 5.0, repo.capturedClaim.rewardAmount)
	require.True(t, status.Claimed)
	require.True(t, status.Eligible)
	require.False(t, status.Claimable)
	require.Equal(t, 7.0, status.ClaimedRewardAmount)
	require.ElementsMatch(t, []int64{50_000_000, 100_000_000}, status.ClaimedThresholdTokens)
}

func TestTokenIncentiveServiceClaim_ThirdTierUsesFullConfiguredReward(t *testing.T) {
	repo := &tokenIncentiveRepoStub{
		tokens: 500_000_000,
		claims: []*TokenIncentiveClaim{
			{ID: 1, UserID: 42, Tokens: 50_000_000, ThresholdTokens: 50_000_000, RewardAmount: 2, ClaimedAt: time.Now()},
			{ID: 2, UserID: 42, Tokens: 100_000_000, ThresholdTokens: 100_000_000, RewardAmount: 5, ClaimedAt: time.Now()},
		},
	}
	svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2},{"threshold_tokens":100000000,"reward_amount":5},{"threshold_tokens":500000000,"reward_amount":10}]`, true)

	status, err := svc.Claim(context.Background(), 42)

	require.NoError(t, err)
	require.True(t, repo.capturedClaim.called)
	require.EqualValues(t, 500_000_000, repo.capturedClaim.thresholdTokens)
	require.Equal(t, 10.0, repo.capturedClaim.rewardAmount)
	require.True(t, status.Claimed)
	require.True(t, status.Eligible)
	require.False(t, status.Claimable)
	require.Equal(t, 17.0, status.ClaimedRewardAmount)
	require.ElementsMatch(t, []int64{50_000_000, 100_000_000, 500_000_000}, status.ClaimedThresholdTokens)
}

func TestTokenIncentiveStatus_LegacyClaimUsesHighestMatchingRewardTier(t *testing.T) {
	rules := []TokenIncentiveRule{
		{ThresholdTokens: 50_000_000, RewardAmount: 2},
		{ThresholdTokens: 100_000_000, RewardAmount: 5},
		{ThresholdTokens: 500_000_000, RewardAmount: 5},
	}

	status := buildTokenIncentiveStatus(
		true,
		time.Now(),
		time.Now().AddDate(0, 0, 7),
		500_000_000,
		[]*TokenIncentiveClaim{{
			ID:           1,
			UserID:       42,
			Tokens:       500_000_000,
			RewardAmount: 5,
			ClaimedAt:    time.Now(),
		}},
		nil,
		rules,
	)

	require.True(t, status.Claimed)
	require.True(t, status.Claimable)
	require.EqualValues(t, 50_000_000, status.ThresholdTokens)
	require.Equal(t, 2.0, status.RewardAmount)
	require.ElementsMatch(t, []int64{500_000_000}, status.ClaimedThresholdTokens)
}

func TestTokenIncentiveServiceClaim_NotEligibleAndAlreadyClaimed(t *testing.T) {
	t.Run("not eligible below first tier", func(t *testing.T) {
		repo := &tokenIncentiveRepoStub{tokens: 49_999_999}
		svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2}]`, true)

		_, err := svc.Claim(context.Background(), 42)

		require.ErrorIs(t, err, ErrTokenIncentiveNotEligible)
		require.False(t, repo.capturedClaim.called)
	})

	t.Run("already claimed before credit", func(t *testing.T) {
		repo := &tokenIncentiveRepoStub{
			tokens: 50_000_000,
			claims: []*TokenIncentiveClaim{{ID: 1, UserID: 42, Tokens: 50_000_000, ThresholdTokens: 50_000_000, RewardAmount: 2, ClaimedAt: time.Now()}},
		}
		svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2}]`, true)

		_, err := svc.Claim(context.Background(), 42)

		require.ErrorIs(t, err, ErrTokenIncentiveAlreadyClaimed)
		require.False(t, repo.capturedClaim.called)
	})
}

func TestTokenIncentiveServiceClaim_DisabledDoesNotReadUsage(t *testing.T) {
	repo := &tokenIncentiveRepoStub{tokens: 500_000_000}
	svc := newTokenIncentiveTestService(repo, `[{"threshold_tokens":50000000,"reward_amount":2}]`, false)

	_, err := svc.Claim(context.Background(), 42)

	require.ErrorIs(t, err, ErrTokenIncentiveDisabled)
	require.False(t, repo.capturedClaim.called)
}

func TestTokenIncentiveStatus_ClaimedKeepsLiveTokensAndClaimedAmount(t *testing.T) {
	claimedAt := time.Now()
	status := buildTokenIncentiveStatus(
		true,
		time.Now(),
		time.Now().AddDate(0, 0, 7),
		600_000_000,
		[]*TokenIncentiveClaim{{
			ID:              1,
			UserID:          42,
			Tokens:          50_000_000,
			ThresholdTokens: 50_000_000,
			RewardAmount:    2,
			ClaimedAt:       claimedAt,
		}, {
			ID:              2,
			UserID:          42,
			Tokens:          100_000_000,
			ThresholdTokens: 100_000_000,
			RewardAmount:    5,
			ClaimedAt:       claimedAt,
		}},
		nil,
		DefaultTokenIncentiveRules(),
	)

	require.True(t, status.Eligible)
	require.True(t, status.Claimable)
	require.True(t, status.Claimed)
	require.EqualValues(t, 600_000_000, status.Tokens)
	require.Equal(t, 10.0, status.RewardAmount)
	require.Equal(t, 7.0, status.ClaimedRewardAmount)
	require.NotNil(t, status.ClaimedAt)
	require.Equal(t, claimedAt, *status.ClaimedAt)
}

func TestParseTokenIncentiveRules_InvalidDoesNotFallThroughSilentlyOnUpdate(t *testing.T) {
	repo := &settingUpdateRepoStub{}
	svc := NewSettingService(repo, &config.Config{})

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		TokenIncentiveRules: []TokenIncentiveRule{{ThresholdTokens: 50_000_000, RewardAmount: TokenIncentiveMaxRewardAmount + 1}},
	})

	require.Error(t, err)
	var appErr *infraerrors.ApplicationError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, "INVALID_TOKEN_INCENTIVE_RULES", appErr.Reason)
}
