package service

import (
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
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
		const vscodeUA = "codex_vscode/9.9.9 (Mac OS X 14.0; arm64) vscode (codex_vscode; 9.9.9)"
		h := make(http.Header)
		h.Set("user-agent", vscodeUA)
		h.Set("version", "9.9.9")
		h.Set("OpenAI-Beta", "assistants=v2")

		ensureCodexIdentityHeaders(h)
		enforceCodexIdentityHeaders(h)

		require.Equal(t, "codex_vscode", h.Get("originator"))
		require.Equal(t, vscodeUA, h.Get("user-agent"))
		require.Equal(t, "9.9.9", h.Get("version"))
		require.Equal(t, "responses=experimental", h.Get("OpenAI-Beta"))
	})

	t.Run("降载身份归一化后保留版本与终端指纹", func(t *testing.T) {
		h := make(http.Header)
		h.Set("user-agent", "codex-tui/9.9.9 (Mac OS X 14.0; arm64) iTerm (codex-tui; 9.9.9)")
		h.Set("version", "9.9.9")

		ensureCodexIdentityHeaders(h)
		enforceCodexIdentityHeaders(h)

		require.Equal(t, "codex_cli_rs", h.Get("originator"))
		require.Equal(t, "codex_cli_rs/9.9.9 (Mac OS X 14.0; arm64) iTerm", h.Get("user-agent"))
		require.Equal(t, "9.9.9", h.Get("version"))
	})
}

func TestEnforceCodexIdentityHeaders(t *testing.T) {
	const tuiUA = "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)"
	// codex-tui 落在上游降载桶，收口时统一改写为 CLI 身份（保留版本/OS/架构/终端指纹）。
	const tuiNormalizedUA = "codex_cli_rs/0.140.2 (Mac OS X 14.0; arm64) iTerm"

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
			name:           "错配 originator 按最终 UA 重配后归一化",
			originator:     "codex_cli_rs",
			userAgent:      tuiUA,
			wantOriginator: "codex_cli_rs",
			wantUA:         tuiNormalizedUA,
		},
		{
			name:           "降载身份改写为 CLI 身份",
			originator:     "codex-tui",
			userAgent:      tuiUA,
			wantOriginator: "codex_cli_rs",
			wantUA:         tuiNormalizedUA,
		},
		{
			name:           "非降载官方身份原样保留",
			originator:     "codex_vscode",
			userAgent:      "codex_vscode/1.2.3 (Ubuntu 22.4.0; x86_64) vscode (codex_vscode; 1.2.3)",
			wantOriginator: "codex_vscode",
			wantUA:         "codex_vscode/1.2.3 (Ubuntu 22.4.0; x86_64) vscode (codex_vscode; 1.2.3)",
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
			name:           "originator override UA 首段被尾部真实身份重写后归一化",
			originator:     "cccc",
			userAgent:      "cccc/0.142.0 (Ubuntu 22.4.0; x86_64) screen (codex-tui; 0.142.0)",
			wantOriginator: "codex_cli_rs",
			wantUA:         "codex_cli_rs/0.142.0 (Ubuntu 22.4.0; x86_64) screen",
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

// 开关是进程级快照，零值 Config（测试 / 工具手工构造，不经 viper）必须落在「归一化开启」
// 一侧，否则任意一处零值构造都会静默关掉全局保护。
//
// 不得给本文件的开关类用例加 t.Parallel()：它们改写进程级状态。
func TestCodexOriginatorNormalizationZeroValueConfigKeepsItEnabled(t *testing.T) {
	var cfg config.Config
	require.False(t, cfg.Gateway.DisableCodexOriginatorNormalization,
		"零值必须表示归一化开启；若改为正向命名的 NormalizeCodexOriginator，零值会静默关闭保护")

	SetCodexOriginatorNormalizationEnabled(!cfg.Gateway.DisableCodexOriginatorNormalization)
	t.Cleanup(func() { SetCodexOriginatorNormalizationEnabled(true) })

	h := make(http.Header)
	h.Set("originator", "codex-tui")
	h.Set("user-agent", "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)")

	enforceCodexIdentityHeaders(h)

	require.Equal(t, "codex_cli_rs", h.Get("originator"))
}

// 关闭归一化后必须完整退回配对语义：降载身份逐字保留，供上游调整分桶后回滚使用。
func TestEnforceCodexIdentityHeaders_NormalizationDisabled(t *testing.T) {
	const tuiUA = "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)"

	SetCodexOriginatorNormalizationEnabled(false)
	t.Cleanup(func() { SetCodexOriginatorNormalizationEnabled(true) })

	h := make(http.Header)
	h.Set("originator", "codex-tui")
	h.Set("user-agent", tuiUA)

	enforceCodexIdentityHeaders(h)

	require.Equal(t, "codex-tui", h.Get("originator"))
	require.Equal(t, tuiUA, h.Get("user-agent"))
}

// 归一化必须是幂等的：重复收口（如透传路径先后经过多次改写）不得反复裁剪 UA。
func TestEnforceCodexIdentityHeaders_NormalizationIsIdempotent(t *testing.T) {
	h := make(http.Header)
	h.Set("originator", "codex-tui")
	h.Set("user-agent", "codex-tui/0.140.2 (Mac OS X 14.0; arm64) iTerm (codex-tui; 0.140.2)")

	enforceCodexIdentityHeaders(h)
	first := h.Get("user-agent")
	enforceCodexIdentityHeaders(h)

	require.Equal(t, first, h.Get("user-agent"))
	require.Equal(t, "codex_cli_rs", h.Get("originator"))
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

func TestApplyOpenAICodexFingerprintHeadersScopesInboundDeviceToAccount(t *testing.T) {
	account := &Account{
		ID:       51,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Credentials: map[string]any{
			"chatgpt_account_id": "acct-51",
		},
		Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1},
	}
	h := make(http.Header)
	h.Set("X-Codex-Installation-ID", "official-installation")
	h.Set("X-Codex-Window-ID", "official-window")

	applyOpenAICodexFingerprintHeaders(h, account, 9, "fallback-window", openAICodexDeviceFingerprint{})

	require.Equal(t, mapOpenAICodexInstallationID(account, "official-installation"), h.Get("X-Codex-Installation-ID"))
	require.Equal(t, mapOpenAICodexFingerprintIdentifier(account, 9, "window", "official-window"), h.Get("X-Codex-Window-ID"))
	require.NotEqual(t, "official-installation", h.Get("X-Codex-Installation-ID"))

	other := &Account{
		ID:          52,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"chatgpt_account_id": "acct-52"},
		Extra:       map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1},
	}
	require.NotEqual(t,
		mapOpenAICodexInstallationID(account, "official-installation"),
		mapOpenAICodexInstallationID(other, "official-installation"),
	)
}

func TestApplyOpenAICodexFingerprintHeadersPreservesLegacyInboundDevice(t *testing.T) {
	account := &Account{ID: 53, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeLegacy}}
	h := make(http.Header)
	h.Set("X-Codex-Installation-ID", "legacy-installation")
	h.Set("X-Codex-Window-ID", "legacy-window")

	applyOpenAICodexFingerprintHeaders(h, account, 9, "fallback-window", openAICodexDeviceFingerprint{})

	require.Equal(t, "legacy-installation", h.Get("X-Codex-Installation-ID"))
	require.Equal(t, "legacy-window", h.Get("X-Codex-Window-ID"))
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

func TestApplyOpenAICodexFingerprintBodyUsesSameInboundInstallationMapping(t *testing.T) {
	account := &Account{
		ID:          62,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Credentials: map[string]any{"chatgpt_account_id": "acct-62"},
		Extra:       map[string]any{openAICodexFingerprintModeExtraKey: openAICodexFingerprintModeV1},
	}
	body := []byte(`{"client_metadata":{"x-codex-installation-id":"inbound-installation"}}`)
	mapped, fingerprint := applyOpenAICodexFingerprintBody(body, account, 17, "", true)
	require.True(t, fingerprint.managed)
	require.Equal(t, mapOpenAICodexInstallationID(account, "inbound-installation"), fingerprint.deviceID)
	require.Equal(t, fingerprint.deviceID, gjson.GetBytes(mapped, "client_metadata.x-codex-installation-id").String())

	header := make(http.Header)
	header.Set("X-Codex-Installation-ID", "inbound-installation")
	applyOpenAICodexFingerprintHeaders(header, account, 17, "", fingerprint)
	require.Equal(t, fingerprint.deviceID, header.Get("X-Codex-Installation-ID"))
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
