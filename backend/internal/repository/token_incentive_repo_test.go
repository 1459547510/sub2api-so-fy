//go:build unit

package repository

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestTokenIncentiveWeeklyUsageSQL_UsesCanonicalTokenTotal(t *testing.T) {
	require.Contains(t, tokenIncentiveWeeklyUsageSQL, "COALESCE(cache_creation_tokens, 0)")
	require.Contains(t, tokenIncentiveWeeklyUsageSQL, "COALESCE(cache_read_tokens, 0)")
	require.NotContains(t, tokenIncentiveWeeklyUsageSQL, "cache_creation_5m_tokens")
	require.NotContains(t, tokenIncentiveWeeklyUsageSQL, "cache_creation_1h_tokens")
	require.NotContains(t, tokenIncentiveWeeklyUsageSQL, "image_output_tokens")
	require.Contains(t, tokenIncentiveWeeklyUsageSQL, "actual_cost > 0")
}

func TestTokenIncentiveRepositoryClaimReward_PassesConfiguredTierAndCreditsBalance(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	defer db.Close()

	repo := NewTokenIncentiveRepository(db)
	ctx := context.Background()
	weekStart := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)
	weekEnd := weekStart.AddDate(0, 0, 7)
	claimedAt := weekStart.Add(12 * time.Hour)
	createdAt := claimedAt
	updatedAt := claimedAt

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO token_incentive_claims")).
		WithArgs(int64(42), weekStart, weekEnd, 5.0, int64(100_000_000)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "week_start", "week_end", "tokens", "reward_amount", "status", "claimed_at", "created_at", "updated_at",
		}).AddRow(int64(9), int64(42), weekStart, weekEnd, int64(120_000_000), 5.0, service.TokenIncentiveClaimedStatus, claimedAt, createdAt, updatedAt))
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE users")).
		WithArgs(5.0, int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(17.0))
	mock.ExpectCommit()

	claim, balanceAfter, err := repo.ClaimReward(ctx, 42, weekStart, weekEnd, 120_000_000, 100_000_000, 5)

	require.NoError(t, err)
	require.EqualValues(t, 9, claim.ID)
	require.EqualValues(t, 120_000_000, claim.Tokens)
	require.Equal(t, 5.0, claim.RewardAmount)
	require.Equal(t, 17.0, balanceAfter)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenIncentiveRepositoryClaimReward_NoRowsMapsAlreadyClaimedBeforeNotEligible(t *testing.T) {
	t.Run("existing claim wins", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		require.NoError(t, err)
		defer db.Close()

		repo := NewTokenIncentiveRepository(db)
		weekStart := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)
		weekEnd := weekStart.AddDate(0, 0, 7)
		now := weekStart.Add(time.Hour)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO token_incentive_claims")).
			WithArgs(int64(42), weekStart, weekEnd, 5.0, int64(100_000_000)).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery("FROM token_incentive_claims").
			WithArgs(int64(42), weekStart).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "user_id", "week_start", "week_end", "tokens", "reward_amount", "status", "claimed_at", "created_at", "updated_at",
			}).AddRow(int64(9), int64(42), weekStart, weekEnd, int64(120_000_000), 5.0, service.TokenIncentiveClaimedStatus, now, now, now))
		mock.ExpectRollback()

		_, _, err = repo.ClaimReward(context.Background(), 42, weekStart, weekEnd, 120_000_000, 100_000_000, 5)

		require.ErrorIs(t, err, service.ErrTokenIncentiveAlreadyClaimed)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no existing claim means not eligible", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		require.NoError(t, err)
		defer db.Close()

		repo := NewTokenIncentiveRepository(db)
		weekStart := time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC)
		weekEnd := weekStart.AddDate(0, 0, 7)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO token_incentive_claims")).
			WithArgs(int64(42), weekStart, weekEnd, 5.0, int64(100_000_000)).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery("FROM token_incentive_claims").
			WithArgs(int64(42), weekStart).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, _, err = repo.ClaimReward(context.Background(), 42, weekStart, weekEnd, 90_000_000, 100_000_000, 5)

		require.ErrorIs(t, err, service.ErrTokenIncentiveNotEligible)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTokenIncentiveClaimInsertSQL_RechecksThresholdInDatabase(t *testing.T) {
	compact := strings.Join(strings.Fields(tokenIncentiveClaimInsertSQL), " ")

	require.Contains(t, compact, "WHERE tokens >= $5")
	require.Contains(t, compact, "ON CONFLICT (user_id, week_start) DO NOTHING")
	require.NotContains(t, compact, "1000000000")
}
