package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildAccountForCreateDefaultsNewOpenAIOAuthToV1Fingerprint(t *testing.T) {
	account, err := buildAccountForCreate(&CreateAccountInput{
		Name:        "new-codex",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"access_token": "token"},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, openAICodexFingerprintModeV1, account.Extra[openAICodexFingerprintModeExtraKey])
}

func TestBuildAccountForCreatePreservesExplicitLegacyFingerprintMode(t *testing.T) {
	account, err := buildAccountForCreate(&CreateAccountInput{
		Name:        "legacy-codex",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"access_token": "token"},
	}, map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeLegacy})
	require.NoError(t, err)
	require.Equal(t, openAICodexFingerprintModeLegacy, account.Extra[openAICodexFingerprintModeExtraKey])
}

func TestMarkOpenAICodexLegacyFingerprintOnlyMarksUnversionedOAuthAccounts(t *testing.T) {
	legacy := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth}
	MarkOpenAICodexLegacyFingerprint(legacy)
	require.Equal(t, openAICodexFingerprintModeLegacy, legacy.Extra[openAICodexFingerprintModeExtraKey])
	require.False(t, openAICodexFingerprintUsesV1(legacy))

	v1 := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Extra:    map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1},
	}
	MarkOpenAICodexLegacyFingerprint(v1)
	require.Equal(t, openAICodexFingerprintModeV1, v1.Extra[openAICodexFingerprintModeExtraKey])
	require.True(t, openAICodexFingerprintUsesV1(v1))

	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey}
	MarkOpenAICodexLegacyFingerprint(apiKey)
	require.Empty(t, apiKey.Extra)
}
