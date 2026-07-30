package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsolateOpenAISessionID(t *testing.T) {
	t.Run("empty_raw_returns_empty", func(t *testing.T) {
		assert.Equal(t, "", isolateOpenAISessionID(1, ""))
		assert.Equal(t, "", isolateOpenAISessionID(1, "   "))
	})

	t.Run("deterministic", func(t *testing.T) {
		a := isolateOpenAISessionID(42, "sess_abc123")
		b := isolateOpenAISessionID(42, "sess_abc123")
		assert.Equal(t, a, b)
	})

	t.Run("different_apiKeyID_different_result", func(t *testing.T) {
		a := isolateOpenAISessionID(1, "same_session")
		b := isolateOpenAISessionID(2, "same_session")
		require.NotEqual(t, a, b, "不同 API Key 使用相同 session_id 应产生不同隔离值")
	})

	t.Run("different_raw_different_result", func(t *testing.T) {
		a := isolateOpenAISessionID(1, "session_a")
		b := isolateOpenAISessionID(1, "session_b")
		require.NotEqual(t, a, b)
	})

	t.Run("format_is_16_hex_chars", func(t *testing.T) {
		result := isolateOpenAISessionID(99, "test_session")
		assert.Len(t, result, 16, "应为 16 字符的 hex 字符串")
		for _, ch := range result {
			assert.True(t, (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f'),
				"应仅包含 hex 字符: %c", ch)
		}
	})

	t.Run("zero_apiKeyID_still_works", func(t *testing.T) {
		result := isolateOpenAISessionID(0, "session")
		assert.NotEmpty(t, result)
		// apiKeyID=0 与 apiKeyID=1 应产生不同结果
		other := isolateOpenAISessionID(1, "session")
		assert.NotEqual(t, result, other)
	})
}

func TestIsolateOpenAIAccountSessionID(t *testing.T) {
	accountA := &Account{ID: 101, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1}}
	accountB := &Account{ID: 102, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1}}

	a1 := isolateOpenAIAccountSessionID(accountA, 7, "conversation")
	a2 := isolateOpenAIAccountSessionID(accountA, 7, "conversation")
	b := isolateOpenAIAccountSessionID(accountB, 7, "conversation")
	otherTenant := isolateOpenAIAccountSessionID(accountA, 8, "conversation")

	assert.Equal(t, a1, a2)
	assert.NotEqual(t, a1, b)
	assert.NotEqual(t, a1, otherTenant)
	assert.Len(t, a1, 16)
}
