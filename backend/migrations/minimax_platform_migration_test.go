package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiniMaxPlatformMigration(t *testing.T) {
	content, err := FS.ReadFile("237_add_minimax_platform.sql")
	require.NoError(t, err)

	sql := strings.Join(strings.Fields(string(content)), " ")
	require.Contains(t, sql, "DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check")
	require.Contains(t, sql, "DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check")
	require.Contains(t, sql, "position('leo' IN quota_constraint_def) = 0")
	require.Contains(t, sql, "position('openai_media' IN quota_constraint_def) = 0")
	require.Contains(t, sql, "position('minimax' IN quota_constraint_def) = 0")
	require.Contains(t, sql, "position('leo' IN route_constraint_def) = 0")
	require.Contains(t, sql, "position('openai_media' IN route_constraint_def) = 0")
	require.Contains(t, sql, "position('minimax' IN route_constraint_def) = 0")
	require.Contains(t, sql,
		"CHECK (platform IN ( 'anthropic', 'openai', 'gemini', 'antigravity', 'grok', 'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax' ))")
	require.Contains(t, sql,
		"CHECK (target_platform IN ( 'anthropic', 'openai', 'gemini', 'antigravity', 'grok', 'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax' ))")
	require.Contains(t, sql,
		"CHECK (provider IN ('openai', 'anthropic', 'gemini', 'grok', 'antigravity', 'kimi', 'zhipu', 'deepseek', 'minimax'))")
	require.NotContains(t, sql,
		"CHECK (platform IN ('anthropic', 'openai', 'gemini', 'antigravity', 'grok', 'kimi', 'zhipu', 'deepseek', 'minimax'))")
}
