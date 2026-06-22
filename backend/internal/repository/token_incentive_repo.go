package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type tokenIncentiveRepository struct {
	db *sql.DB
}

// NewTokenIncentiveRepository creates a SQL-backed repository for token incentive claims.
func NewTokenIncentiveRepository(db *sql.DB) service.TokenIncentiveRepository {
	return &tokenIncentiveRepository{db: db}
}

func (r *tokenIncentiveRepository) GetWeeklyUsageTokens(ctx context.Context, userID int64, weekStart, weekEnd time.Time) (int64, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("token incentive repository is not initialized")
	}
	var tokens int64
	if err := r.db.QueryRowContext(ctx, tokenIncentiveWeeklyUsageSQL, userID, weekStart, weekEnd).Scan(&tokens); err != nil {
		return 0, fmt.Errorf("sum weekly token incentive usage: %w", err)
	}
	return tokens, nil
}

func (r *tokenIncentiveRepository) GetClaim(ctx context.Context, userID int64, weekStart time.Time) (*service.TokenIncentiveClaim, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("token incentive repository is not initialized")
	}
	claim, err := scanTokenIncentiveClaim(r.db.QueryRowContext(ctx, tokenIncentiveClaimSelectSQL+`
WHERE user_id = $1 AND week_start = $2
LIMIT 1`, userID, weekStart))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get token incentive claim: %w", err)
	}
	return claim, nil
}

func (r *tokenIncentiveRepository) ClaimReward(ctx context.Context, userID int64, weekStart, weekEnd time.Time, _ int64, rewardAmount float64) (*service.TokenIncentiveClaim, float64, error) {
	if r == nil || r.db == nil {
		return nil, 0, fmt.Errorf("token incentive repository is not initialized")
	}
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, 0, fmt.Errorf("begin token incentive claim tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	claim, err := scanTokenIncentiveClaim(tx.QueryRowContext(ctx, tokenIncentiveClaimInsertSQL,
		userID, weekStart, weekEnd, rewardAmount, service.TokenIncentiveThresholdTokens,
	))
	if errors.Is(err, sql.ErrNoRows) {
		existingClaim, getErr := scanTokenIncentiveClaim(tx.QueryRowContext(ctx, tokenIncentiveClaimSelectSQL+`
WHERE user_id = $1 AND week_start = $2
LIMIT 1`, userID, weekStart))
		if getErr == nil && existingClaim != nil {
			return nil, 0, service.ErrTokenIncentiveAlreadyClaimed
		}
		if getErr != nil && !errors.Is(getErr, sql.ErrNoRows) {
			return nil, 0, fmt.Errorf("check existing token incentive claim: %w", getErr)
		}
		return nil, 0, service.ErrTokenIncentiveNotEligible
	}
	if err != nil {
		return nil, 0, fmt.Errorf("insert token incentive claim: %w", err)
	}

	// Token incentive is a cashback/reward credit, not a user recharge. Keep
	// total_recharged unchanged so recharge analytics and downstream rebate logic
	// are not inflated by self-claimed rewards.
	var balanceAfter float64
	if err := tx.QueryRowContext(ctx, `
UPDATE users
SET balance = balance + $1,
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL
RETURNING balance::double precision`, rewardAmount, userID).Scan(&balanceAfter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, service.ErrUserNotFound
		}
		return nil, 0, fmt.Errorf("credit token incentive reward: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("commit token incentive claim tx: %w", err)
	}
	return claim, balanceAfter, nil
}

const tokenIncentiveWeeklyUsageSQL = `
SELECT COALESCE(SUM(
    COALESCE(input_tokens, 0)::bigint +
    COALESCE(output_tokens, 0)::bigint +
    COALESCE(cache_creation_tokens, 0)::bigint +
    COALESCE(cache_read_tokens, 0)::bigint +
    COALESCE(cache_creation_5m_tokens, 0)::bigint +
    COALESCE(cache_creation_1h_tokens, 0)::bigint +
    COALESCE(image_output_tokens, 0)::bigint
), 0) AS tokens
FROM usage_logs
WHERE user_id = $1
  AND created_at >= $2
  AND created_at < $3
  AND actual_cost > 0`

const tokenIncentiveClaimInsertSQL = `
WITH weekly_usage AS (
    ` + tokenIncentiveWeeklyUsageSQL + `
)
INSERT INTO token_incentive_claims (user_id, week_start, week_end, tokens, reward_amount, status)
SELECT $1, $2, $3, tokens, $4, 'claimed'
FROM weekly_usage
WHERE tokens >= $5
ON CONFLICT (user_id, week_start) DO NOTHING
RETURNING id, user_id, week_start, week_end, tokens, reward_amount::double precision, status, claimed_at, created_at, updated_at`

const tokenIncentiveClaimSelectSQL = `
SELECT id, user_id, week_start, week_end, tokens, reward_amount::double precision, status, claimed_at, created_at, updated_at
FROM token_incentive_claims
`

type tokenIncentiveClaimScanner interface {
	Scan(dest ...any) error
}

func scanTokenIncentiveClaim(row tokenIncentiveClaimScanner) (*service.TokenIncentiveClaim, error) {
	var claim service.TokenIncentiveClaim
	if err := row.Scan(
		&claim.ID,
		&claim.UserID,
		&claim.WeekStart,
		&claim.WeekEnd,
		&claim.Tokens,
		&claim.RewardAmount,
		&claim.Status,
		&claim.ClaimedAt,
		&claim.CreatedAt,
		&claim.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &claim, nil
}
