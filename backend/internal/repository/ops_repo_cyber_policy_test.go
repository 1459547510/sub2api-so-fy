package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestFindCyberPolicyBanCandidateRanksByCountThenLatestHit(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &opsRepository{db: db}

	since := time.Date(2026, 6, 24, 8, 0, 0, 0, time.UTC)
	until := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	lastHit := until.Add(-20 * time.Minute)
	apiKeyID := int64(88)

	mock.ExpectQuery(`(?s)e\.error_type = 'cyber_policy'.*e\.created_at >= \$2.*e\.created_at <= \$3.*u\.role <> \$4.*ORDER BY hit_count DESC, last_hit_at DESC, e\.user_id ASC.*latest_hit.*moderation\.action = 'cyber_policy'`).
		WithArgs(int64(1714), since, until, service.RoleAdmin).
		WillReturnRows(cyberPolicyCandidateRows().AddRow(
			int64(147), int64(3), lastHit, int64(9001), "req-cyber", "client-cyber", apiKeyID,
			"sk-test1", "203.0.113.9", "claude-code/1.0", "gpt-5.5", "/v1/responses", "retained prompt excerpt",
		))

	candidate, err := repo.FindCyberPolicyBanCandidate(context.Background(), 1714, since, until)
	require.NoError(t, err)
	require.Equal(t, &service.CyberPolicyBanCandidate{
		UserID:          147,
		HitCount:        3,
		LastHitAt:       lastHit,
		OpsErrorLogID:   9001,
		RequestID:       "req-cyber",
		ClientRequestID: "client-cyber",
		APIKeyID:        &apiKeyID,
		APIKeyPrefix:    "sk-test1",
		ClientIP:        "203.0.113.9",
		UserAgent:       "claude-code/1.0",
		Model:           "gpt-5.5",
		RequestPath:     "/v1/responses",
		InputExcerpt:    "retained prompt excerpt",
	}, candidate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCyberPolicyBanCandidateIgnoresSessionBlockedRows(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &opsRepository{db: db}

	since := time.Date(2026, 6, 24, 8, 0, 0, 0, time.UTC)
	until := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)e\.error_type = 'cyber_policy'.*LIMIT 1`).
		WithArgs(int64(1714), since, until, service.RoleAdmin).
		WillReturnRows(cyberPolicyCandidateRows())

	candidate, err := repo.FindCyberPolicyBanCandidate(context.Background(), 1714, since, until)
	require.NoError(t, err)
	require.Nil(t, candidate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func cyberPolicyCandidateRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"user_id", "hit_count", "last_hit_at", "ops_error_log_id", "request_id", "client_request_id",
		"api_key_id", "api_key_prefix", "client_ip", "user_agent", "model", "request_path", "input_excerpt",
	})
}
