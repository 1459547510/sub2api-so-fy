package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserPlatformQuotasKeepMediaPlatformsMigration(t *testing.T) {
	content, err := FS.ReadFile("227_user_platform_quotas_keep_media_platforms.sql")
	require.NoError(t, err)

	sql := strings.Join(strings.Fields(string(content)), " ")
	require.Contains(t, sql, "DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check")
	require.Contains(t, sql, "position('leo' IN constraint_def) = 0")
	require.Contains(t, sql, "position('openai_media' IN constraint_def) = 0")
	require.Contains(t, sql, "position('kimi' IN constraint_def) = 0")
	require.Contains(t, sql,
		"CHECK (platform IN ( 'anthropic', 'openai', 'gemini', 'antigravity', 'grok', 'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek' ))")
}
