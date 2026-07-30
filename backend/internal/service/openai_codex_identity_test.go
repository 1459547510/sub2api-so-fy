package service

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func requireOpenAICodexProbeHeaders(t *testing.T, h http.Header) {
	t.Helper()
	require.Equal(t, codexCLIUserAgent, h.Get("User-Agent"))
	require.Equal(t, "codex_cli_rs", h.Get("Originator"))
	require.Equal(t, codexCLIVersion, h.Get("Version"))
	require.Equal(t, "responses=experimental", h.Get("OpenAI-Beta"))
	require.NotEmpty(t, h.Get("X-Codex-Window-ID"))
}

func TestEnsureCodexIdentityHeaders(t *testing.T) {
	t.Run("补齐缺失身份头", func(t *testing.T) {
		h := make(http.Header)

		ensureCodexIdentityHeaders(h)
		enforceCodexIdentityHeaders(h)

		require.Equal(t, "codex_cli_rs", h.Get("originator"))
		require.Equal(t, codexCLIUserAgent, h.Get("user-agent"))
		require.Equal(t, codexCLIVersion, h.Get("version"))
		require.Equal(t, "responses=experimental", h.Get("OpenAI-Beta"))
	})

	t.Run("保留已有官方UA和合法version并重新配对", func(t *testing.T) {
		const tuiUA = "codex-tui/9.9.9 (Mac OS X 14.0; arm64) iTerm (codex-tui; 9.9.9)"
		h := make(http.Header)
		h.Set("user-agent", tuiUA)
		h.Set("version", "9.9.9")
		h.Set("OpenAI-Beta", "assistants=v2")

		ensureCodexIdentityHeaders(h)
		enforceCodexIdentityHeaders(h)

		require.Equal(t, "codex-tui", h.Get("originator"))
		require.Equal(t, tuiUA, h.Get("user-agent"))
		require.Equal(t, "9.9.9", h.Get("version"))
		require.Equal(t, "responses=experimental", h.Get("OpenAI-Beta"))
	})
}

func TestEnforceCodexIdentityHeaders(t *testing.T) {
	const tuiUA = "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)"

	tests := []struct {
		name           string
		originator     string
		userAgent      string
		version        string
		wantOriginator string
		wantUA         string
		wantVersion    string
	}{
		{
			name:           "错配 originator 按最终 UA 重配",
			originator:     "codex_cli_rs",
			userAgent:      tuiUA,
			wantOriginator: "codex-tui",
			wantUA:         tuiUA,
		},
		{
			name:           "官方配套身份原样保留",
			originator:     "codex-tui",
			userAgent:      tuiUA,
			wantOriginator: "codex-tui",
			wantUA:         tuiUA,
		},
		{
			name:           "第三方 UA 整体回退默认身份",
			originator:     "opencode",
			userAgent:      "luna/1.0.0",
			wantOriginator: "codex_cli_rs",
			wantUA:         codexCLIUserAgent,
		},
		{
			name:           "UA 缺失回退默认身份",
			originator:     "codex_vscode",
			wantOriginator: "codex_cli_rs",
			wantUA:         codexCLIUserAgent,
		},
		{
			name:           "originator override UA 首段被尾部真实身份重写",
			originator:     "cccc",
			userAgent:      "cccc/0.142.0 (Ubuntu 22.4.0; x86_64) screen (codex-tui; 0.142.0)",
			wantOriginator: "codex-tui",
			wantUA:         "codex-tui/0.142.0 (Ubuntu 22.4.0; x86_64) screen (codex-tui; 0.142.0)",
		},
		{
			name:           "低于门槛的 version 提升为内置版本",
			originator:     "codex_cli_rs",
			userAgent:      "codex_cli_rs/0.125.0",
			version:        "0.125.0",
			wantOriginator: "codex_cli_rs",
			wantUA:         "codex_cli_rs/0.125.0",
			wantVersion:    codexCLIVersion,
		},
		{
			name:           "达标 version 原样保留",
			originator:     "codex_cli_rs",
			userAgent:      "codex_cli_rs/0.145.0",
			version:        "0.145.0",
			wantOriginator: "codex_cli_rs",
			wantUA:         "codex_cli_rs/0.145.0",
			wantVersion:    "0.145.0",
		},
		{
			name:           "未携带 version 不注入",
			originator:     "codex_cli_rs",
			userAgent:      "codex_cli_rs/0.98.0",
			wantOriginator: "codex_cli_rs",
			wantUA:         "codex_cli_rs/0.98.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := make(http.Header)
			if tt.originator != "" {
				h.Set("originator", tt.originator)
			}
			if tt.userAgent != "" {
				h.Set("user-agent", tt.userAgent)
			}
			if tt.version != "" {
				h.Set("version", tt.version)
			}

			enforceCodexIdentityHeaders(h)

			require.Equal(t, tt.wantOriginator, h.Get("originator"))
			require.Equal(t, tt.wantUA, h.Get("user-agent"))
			require.Equal(t, tt.wantVersion, h.Get("version"))
		})
	}
}

// enforce 本身仍只负责收口：缺少 originator 时必须保持 no-op，由需要恢复身份的
// 调用方先显式调用 ensureCodexIdentityHeaders。
func TestEnforceCodexIdentityHeaders_NoOriginatorIsNoop(t *testing.T) {
	h := make(http.Header)
	h.Set("user-agent", "third-party-client/1.0.0")

	enforceCodexIdentityHeaders(h)

	require.Empty(t, h.Get("originator"))
	require.Equal(t, "third-party-client/1.0.0", h.Get("user-agent"))
}

func TestOpenAICodexDeviceFingerprintLifecycle(t *testing.T) {
	account := &Account{
		ID:       41,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"chatgpt_account_id": "acct-41",
			"access_token":       "token-before-refresh",
		},
	}

	first := resolveOpenAICodexDeviceFingerprint(account, "")
	require.True(t, first.managed)
	require.NotEmpty(t, first.deviceID)
	_, err := uuid.Parse(first.deviceID)
	require.NoError(t, err)

	account.Credentials["access_token"] = "token-after-refresh"
	proxyID := int64(7001)
	account.ProxyID = &proxyID
	second := resolveOpenAICodexDeviceFingerprint(account, "")
	require.Equal(t, first, second, "token and proxy changes must not rotate a device")

	other := &Account{
		ID:          42,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"chatgpt_account_id": "acct-42"},
	}
	require.NotEqual(t, first.deviceID, resolveOpenAICodexDeviceFingerprint(other, "").deviceID)

	account.Extra = map[string]any{"openai_device_profile_id": "generation-2"}
	require.NotEqual(t, first.deviceID, resolveOpenAICodexDeviceFingerprint(account, "").deviceID)
}

func TestApplyOpenAICodexFingerprintHeadersPreservesOfficialInboundDevice(t *testing.T) {
	account := &Account{ID: 51, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1}}
	h := make(http.Header)
	h.Set("X-Codex-Installation-ID", "official-installation")
	h.Set("X-Codex-Window-ID", "official-window")

	applyOpenAICodexFingerprintHeaders(h, account, 9, "fallback-window", openAICodexDeviceFingerprint{})

	require.Equal(t, "official-installation", h.Get("X-Codex-Installation-ID"))
	require.Equal(t, "official-window", h.Get("X-Codex-Window-ID"))
}

func TestApplyOpenAICodexFingerprintHeadersUsesManagedAccountDevice(t *testing.T) {
	account := &Account{
		ID:       52,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Extra:    map[string]any{"openai_device_id": "managed-installation", openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1},
	}
	h := make(http.Header)
	h.Set("X-Codex-Installation-ID", "downstream-installation")
	h.Set("X-Codex-Window-ID", "downstream-window")

	applyOpenAICodexFingerprintHeaders(h, account, 9, "", openAICodexDeviceFingerprint{})

	require.Equal(t, "managed-installation", h.Get("X-Codex-Installation-ID"))
	require.Equal(t, mapOpenAICodexFingerprintIdentifier(account, 9, "window", "downstream-window"), h.Get("X-Codex-Window-ID"))
}

func TestApplyOpenAICodexFingerprintBodyAlignsSessionAndDevice(t *testing.T) {
	account := &Account{ID: 61, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1}}
	body := []byte(`{"prompt_cache_key":"conversation-1","client_metadata":{"x-codex-window-id":"window-1"}}`)

	mapped, fingerprint := applyOpenAICodexFingerprintBody(body, account, 17, "", true)

	require.True(t, fingerprint.managed)
	require.Equal(t, isolateOpenAIAccountSessionID(account, 17, "conversation-1"), gjson.GetBytes(mapped, "prompt_cache_key").String())
	require.Equal(t, fingerprint.deviceID, gjson.GetBytes(mapped, "client_metadata.x-codex-installation-id").String())
	require.Equal(t, mapOpenAICodexFingerprintIdentifier(account, 17, "window", "window-1"), gjson.GetBytes(mapped, "client_metadata.x-codex-window-id").String())

	mappedAgain, fingerprintAgain := applyOpenAICodexFingerprintBody(body, account, 17, "", true)
	require.Equal(t, string(mapped), string(mappedAgain))
	require.Equal(t, fingerprint, fingerprintAgain)
}

func TestOpenAICodexLegacyFingerprintRemainsUnchanged(t *testing.T) {
	account := &Account{ID: 71, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeLegacy}}
	header := http.Header{}
	header.Set("X-Codex-Installation-ID", "legacy-installation")
	header.Set("X-Codex-Window-ID", "legacy-window")
	applyOpenAICodexFingerprintHeaders(header, account, 9, "new-window-seed", openAICodexDeviceFingerprint{})
	require.Equal(t, "legacy-installation", header.Get("X-Codex-Installation-ID"))
	require.Equal(t, "legacy-window", header.Get("X-Codex-Window-ID"))

	body := []byte(`{"prompt_cache_key":"legacy-cache","client_metadata":{"x-codex-installation-id":"legacy-installation"}}`)
	mapped, fingerprint := applyOpenAICodexFingerprintBody(body, account, 9, "", true)
	require.Equal(t, string(body), string(mapped))
	require.Empty(t, fingerprint.deviceID)
	require.Equal(t, isolateOpenAISessionID(9, "legacy-cache"), isolateOpenAIAccountSessionID(account, 9, "legacy-cache"))
}
