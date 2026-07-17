package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVideoJobsMigration(t *testing.T) {
	content, err := FS.ReadFile("182_video_jobs.sql")
	require.NoError(t, err)

	sql := strings.Join(strings.Fields(string(content)), " ")
	require.Contains(t, sql, "CREATE TABLE IF NOT EXISTS video_jobs")
	require.Contains(t, sql, "upstream_job_id BIGINT")
	require.Contains(t, sql, "billing_snapshot JSONB")
	require.Contains(t, sql, "CREATE UNIQUE INDEX IF NOT EXISTS idx_video_jobs_account_upstream_job")
	require.Contains(t, sql, "WHERE upstream_job_id IS NOT NULL")
	require.Contains(t, sql, "CREATE INDEX IF NOT EXISTS idx_video_jobs_api_key_created")
	require.Contains(t, sql, "CREATE INDEX IF NOT EXISTS idx_video_jobs_status_updated")
}
