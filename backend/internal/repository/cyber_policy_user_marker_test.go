package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestUserRepositoryCyberPolicyUserMarker(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := &userRepository{sql: db}

	mock.ExpectExec("INSERT INTO cyber_policy_user_marks").
		WithArgs(int64(7), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	marked, err := repo.MarkCyberPolicyUser(context.Background(), 7, time.Unix(100, 0))
	require.NoError(t, err)
	require.True(t, marked)

	mock.ExpectExec("INSERT INTO cyber_policy_user_marks").
		WithArgs(int64(8), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))
	marked, err = repo.MarkCyberPolicyUser(context.Background(), 8, time.Unix(101, 0))
	require.NoError(t, err)
	require.False(t, marked)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	has, err := repo.HasCyberPolicyUser(context.Background(), 7)
	require.NoError(t, err)
	require.True(t, has)

	require.NoError(t, mock.ExpectationsWereMet())
}
