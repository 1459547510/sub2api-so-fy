package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestReserveUsageBillingVideoBalanceMovesAvailableToFrozen(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	mock.ExpectQuery(`(?s)UPDATE users\s+SET balance = balance - \$1,\s+frozen_balance = COALESCE\(frozen_balance, 0\) \+ \$1`).
		WithArgs(1.2, int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"balance", "frozen_balance"}).AddRow(8.8, 1.2))
	mock.ExpectCommit()

	result, err := reserveUsageBillingVideoBalance(ctx, tx, &service.VideoBalanceHoldCommand{UserID: 42, HoldAmount: 1.2})
	require.NoError(t, err)
	require.InDelta(t, 8.8, *result.NewBalance, 1e-12)
	require.InDelta(t, 1.2, *result.FrozenBalance, 1e-12)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReleaseUsageBillingVideoBalanceRequiresOriginalHold(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	mock.ExpectQuery(`SELECT 1\s+FROM usage_billing_dedup\s+WHERE request_id = \$1 AND api_key_id = \$2`).
		WithArgs(service.VideoHoldRequestID("vidjob_1"), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}).AddRow(1))
	mock.ExpectQuery(`(?s)UPDATE users\s+SET balance = balance \+ \$1,\s+frozen_balance = COALESCE\(frozen_balance, 0\) - \$1`).
		WithArgs(1.2, int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"balance", "frozen_balance"}).AddRow(10.0, 0.0))
	mock.ExpectCommit()

	result, err := releaseUsageBillingVideoBalance(ctx, tx, &service.VideoBalanceHoldCommand{UserID: 42, APIKeyID: 7, JobID: "vidjob_1", HoldAmount: 1.2})
	require.NoError(t, err)
	require.InDelta(t, 10.0, *result.NewBalance, 1e-12)
	require.InDelta(t, 0.0, *result.FrozenBalance, 1e-12)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}
