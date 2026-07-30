package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const openAICodexFingerprintProfileVersion = "codex-cli-rs-v1"

const (
	openAICodexFingerprintModeExtraKey = "openai_fingerprint_mode"
	openAICodexFingerprintModeLegacy   = "legacy"
	openAICodexFingerprintModeV1       = "v1"
)

var openAICodexFingerprintNamespace = uuid.NewSHA1(uuid.NameSpaceOID, []byte("sub2api:openai-codex:fingerprint:"+openAICodexFingerprintProfileVersion))

type openAICodexDeviceFingerprint struct {
	deviceID string
	managed  bool
}

// openAICodexFingerprintUsesV1 gates the new fingerprint contract. Accounts
// created before this flag existed have no marker and intentionally retain the
// previous behavior until explicitly migrated.
func openAICodexFingerprintUsesV1(account *Account) bool {
	if account == nil || !account.IsOpenAIOAuth() {
		return false
	}
	mode := strings.TrimSpace(account.GetExtraString(openAICodexFingerprintModeExtraKey))
	if mode == "" {
		return true
	}
	return strings.EqualFold(mode, openAICodexFingerprintModeV1)
}

func openAICodexFingerprintMode(account *Account) string {
	if openAICodexFingerprintUsesV1(account) {
		return openAICodexFingerprintModeV1
	}
	return openAICodexFingerprintModeLegacy
}

// MarkOpenAICodexLegacyFingerprint marks an unversioned persisted account as
// legacy without writing the compatibility marker back to the database.
func MarkOpenAICodexLegacyFingerprint(account *Account) {
	if account == nil || !account.IsOpenAIOAuth() {
		return
	}
	if strings.TrimSpace(account.GetExtraString(openAICodexFingerprintModeExtraKey)) == "" {
		if account.Extra == nil {
			account.Extra = make(map[string]any)
		}
		account.Extra[openAICodexFingerprintModeExtraKey] = openAICodexFingerprintModeLegacy
	}
}

// codexUpstreamMinVersion 上游 /backend-api/codex 接受的最低 version 头：
// 若请求携带 version 且低于该值，上游直接 404（issue #3901，2026-07 实测）。
const codexUpstreamMinVersion = "0.144.0"

// ensureCodexIdentityHeaders 补齐 OAuth（ChatGPT 内部接口）出站请求所需的 Codex 身份头。
// 已有 User-Agent 与 version 保持不变，交给紧随其后的 enforceCodexIdentityHeaders
// 做官方身份配对与最低版本校正。
func ensureCodexIdentityHeaders(h http.Header) {
	if h == nil {
		return
	}
	if strings.TrimSpace(h.Get("user-agent")) == "" {
		h.Set("user-agent", codexCLIUserAgent)
	}
	if strings.TrimSpace(h.Get("originator")) == "" {
		h.Set("originator", "codex_cli_rs")
	}
	if strings.TrimSpace(h.Get("version")) == "" {
		h.Set("version", codexCLIVersion)
	}
	h.Set("OpenAI-Beta", "responses=experimental")
}

// applyOpenAICodexProbeHeaders 为合成探测请求补齐 Codex 身份和引擎指纹。
func applyOpenAICodexProbeHeaders(h http.Header) {
	if h == nil {
		return
	}
	ensureCodexIdentityHeaders(h)
	h.Set("X-Codex-Window-ID", uuid.NewString())
}

// openAICodexFingerprintAccountKey returns a stable upstream-account identity.
// Tokens and proxy addresses are intentionally excluded so refreshes and normal
// network changes do not make a long-lived client look like a new installation.
func openAICodexFingerprintAccountKey(account *Account) string {
	if account == nil || !account.IsOpenAIOAuth() {
		return ""
	}
	accountKey := strings.TrimSpace(account.GetChatGPTAccountID())
	if accountKey == "" && account.ParentAccountID != nil && *account.ParentAccountID > 0 {
		accountKey = "local-parent:" + fmt.Sprint(*account.ParentAccountID)
	}
	if accountKey == "" && account.ID > 0 {
		accountKey = "local:" + fmt.Sprint(account.ID)
	}
	if accountKey == "" {
		return ""
	}
	if profileID := strings.TrimSpace(account.GetExtraString("openai_device_profile_id")); profileID != "" {
		accountKey += "|device-profile:" + profileID
	}
	return accountKey
}

func resolveOpenAICodexDeviceFingerprint(account *Account, inboundDeviceID string) openAICodexDeviceFingerprint {
	if account == nil || !account.IsOpenAIOAuth() {
		return openAICodexDeviceFingerprint{}
	}
	if configured := strings.TrimSpace(account.GetOpenAIDeviceID()); configured != "" {
		return openAICodexDeviceFingerprint{deviceID: configured, managed: true}
	}
	if inbound := strings.TrimSpace(inboundDeviceID); inbound != "" {
		return openAICodexDeviceFingerprint{deviceID: inbound}
	}
	accountKey := openAICodexFingerprintAccountKey(account)
	if accountKey == "" {
		return openAICodexDeviceFingerprint{}
	}
	deviceID := uuid.NewSHA1(openAICodexFingerprintNamespace, []byte("installation:"+accountKey)).String()
	return openAICodexDeviceFingerprint{deviceID: deviceID, managed: true}
}

func mapOpenAICodexFingerprintIdentifier(account *Account, apiKeyID int64, kind, raw string) string {
	raw = strings.TrimSpace(raw)
	accountKey := openAICodexFingerprintAccountKey(account)
	if raw == "" || accountKey == "" {
		return raw
	}
	name := fmt.Sprintf("%s:account:%s:api-key:%d:value:%s", kind, accountKey, apiKeyID, raw)
	return uuid.NewSHA1(openAICodexFingerprintNamespace, []byte(name)).String()
}

// applyOpenAICodexFingerprintHeaders aligns installation/window headers at the
// final outbound boundary. Official inbound device identities are preserved;
// managed or synthetic identities are scoped to the selected upstream account.
func applyOpenAICodexFingerprintHeaders(h http.Header, account *Account, apiKeyID int64, fallbackWindowSeed string, resolved openAICodexDeviceFingerprint) {
	if h == nil || !openAICodexFingerprintUsesV1(account) {
		return
	}
	fingerprint := resolved
	if fingerprint.deviceID == "" {
		fingerprint = resolveOpenAICodexDeviceFingerprint(account, h.Get("X-Codex-Installation-ID"))
	}
	if fingerprint.deviceID != "" {
		h.Set("X-Codex-Installation-ID", fingerprint.deviceID)
	}
	if !fingerprint.managed {
		return
	}
	windowSeed := strings.TrimSpace(h.Get("X-Codex-Window-ID"))
	if windowSeed == "" {
		windowSeed = strings.TrimSpace(fallbackWindowSeed)
	}
	if windowSeed != "" {
		h.Set("X-Codex-Window-ID", mapOpenAICodexFingerprintIdentifier(account, apiKeyID, "window", windowSeed))
	}
}

// applyOpenAICodexFingerprintBody keeps the request body in the same device
// domain as the headers while preserving the original JSON representation.
func applyOpenAICodexFingerprintBody(body []byte, account *Account, apiKeyID int64, inboundDeviceID string, includeClientMetadata bool) ([]byte, openAICodexDeviceFingerprint) {
	if len(body) == 0 || !openAICodexFingerprintUsesV1(account) || !gjson.ValidBytes(body) {
		return body, openAICodexDeviceFingerprint{}
	}
	promptCacheKey := strings.TrimSpace(gjson.GetBytes(body, "prompt_cache_key").String())
	if promptCacheKey != "" && openAICodexFingerprintAccountKey(account) != "" {
		mapped := isolateOpenAIAccountSessionID(account, apiKeyID, promptCacheKey)
		if updated, err := sjson.SetBytes(body, "prompt_cache_key", mapped); err == nil {
			body = updated
		}
	}
	if !includeClientMetadata {
		return body, openAICodexDeviceFingerprint{}
	}

	bodyDeviceID := strings.TrimSpace(gjson.GetBytes(body, "client_metadata.x-codex-installation-id").String())
	if bodyDeviceID == "" {
		bodyDeviceID = strings.TrimSpace(inboundDeviceID)
	}
	fingerprint := resolveOpenAICodexDeviceFingerprint(account, bodyDeviceID)
	metadata := gjson.GetBytes(body, "client_metadata")
	if fingerprint.deviceID != "" && (!metadata.Exists() || metadata.IsObject()) {
		if updated, err := sjson.SetBytes(body, "client_metadata.x-codex-installation-id", fingerprint.deviceID); err == nil {
			body = updated
		}
	}
	if fingerprint.managed {
		windowID := strings.TrimSpace(gjson.GetBytes(body, "client_metadata.x-codex-window-id").String())
		if windowID != "" {
			mapped := mapOpenAICodexFingerprintIdentifier(account, apiKeyID, "window", windowID)
			if updated, err := sjson.SetBytes(body, "client_metadata.x-codex-window-id", mapped); err == nil {
				body = updated
			}
		}
	}
	return body, fingerprint
}

// enforceCodexIdentityHeaders 收口 OAuth（ChatGPT 内部接口）出站请求的客户端身份头。
// 上游要求 originator 与 User-Agent 首段配套且为官方客户端标识，version 头（若携带）
// 不低于 0.144.0，任一不满足即 404（issue #3901）。以最终 User-Agent 为准推导配套
// originator；推导不出官方身份（第三方 UA / UA 缺失）时整体回退为默认 Codex CLI 身份。
//
// 仅对携带 originator 的请求生效；需要从缺失身份头恢复的调用方应先调用
// ensureCodexIdentityHeaders。
// 必须在所有 User-Agent 改写（自定义 UA / ForceCodexCLI / 浏览器 UA 兜底）之后调用。
func enforceCodexIdentityHeaders(h http.Header) {
	if h == nil || h.Get("originator") == "" {
		return
	}
	originator, pairedUA, ok := openai.PairCodexClientIdentity(h.Get("user-agent"))
	if !ok {
		originator, pairedUA = "codex_cli_rs", codexCLIUserAgent
	}
	h.Set("user-agent", pairedUA)
	h.Set("originator", originator)
	if v := strings.TrimSpace(h.Get("version")); v != "" && CompareVersions(v, codexUpstreamMinVersion) < 0 {
		h.Set("version", codexCLIVersion)
	}
}
