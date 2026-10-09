import type { GroupPlatform } from '@/types'

export const OPENAI_CC_SWITCH_CODEX_MODEL = 'gpt-5.5'
export const GROK_CC_SWITCH_MODEL = 'grok-4.6'

export type CcSwitchClientType = 'claude' | 'gemini' | 'codex' | 'grok'

export interface CcSwitchImportConfig {
  app: string
  endpoint: string
  model?: string
  haikuModel?: string
  sonnetModel?: string
  opusModel?: string
  configJson?: string
}

export interface CcSwitchImportDeeplinkInput {
  baseUrl: string
  platform?: GroupPlatform | null
  clientType: CcSwitchClientType
  providerName: string
  apiKey: string
  usageScript: string
}

/**
 * Balance query CC Switch runs against the imported provider. CC Switch fills
 * `{{baseUrl}}` with the provider's base URL as stored — Codex and Grok imports
 * carry a trailing `/v1` (see `withV1Endpoint`), Claude ones do not, and users
 * may edit it either way afterwards — then evaluates the script, so the URL
 * strips an existing `/v1` instead of blindly appending one (`/v1/v1/usage`
 * is a 404 and CC Switch shows "query failed").
 */
export const CC_SWITCH_USAGE_SCRIPT = `({
    request: {
      url: "{{baseUrl}}".replace(/\\/+$/, "").replace(/\\/v1$/, "") + "/v1/usage",
      method: "GET",
      headers: { "Authorization": "Bearer {{apiKey}}" }
    },
    extractor: function(response) {
      const remaining = response?.remaining ?? response?.quota?.remaining ?? response?.balance;
      const unit = response?.unit ?? response?.quota?.unit ?? "USD";
      return {
        isValid: response?.is_active ?? response?.isValid ?? true,
        remaining,
        unit
      };
    }
  })`

export function withV1Endpoint(baseUrl: string): string {
  const normalizedBaseUrl = baseUrl.replace(/\/+$/, '')
  return normalizedBaseUrl.endsWith('/v1') ? normalizedBaseUrl : `${normalizedBaseUrl}/v1`
}

function withoutTrailingSlashes(baseUrl: string): string {
  return baseUrl.replace(/\/+$/, '')
}

export function stripV1Endpoint(baseUrl: string): string {
  return baseUrl.replace(/\/+$/, '').replace(/\/v1$/i, '')
}

export function grokClaudeCodeEnv(model: string = GROK_CC_SWITCH_MODEL): Record<string, string> {
  return {
    ANTHROPIC_MODEL: model,
    ANTHROPIC_DEFAULT_OPUS_MODEL: model,
    ANTHROPIC_DEFAULT_SONNET_MODEL: model,
    ANTHROPIC_DEFAULT_HAIKU_MODEL: model,
    ANTHROPIC_DEFAULT_FABLE_MODEL: model,
    CLAUDE_CODE_SUBAGENT_MODEL: model,
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: '1'
  }
}

export function encodeCcSwitchConfig(value: string): string {
  const bytes = new TextEncoder().encode(value)
  let binary = ''
  for (const byte of bytes) {
    binary += String.fromCharCode(byte)
  }
  return btoa(binary)
}

export function decodeCcSwitchConfig(value: string): string {
  const binary = atob(value)
  const bytes = Uint8Array.from(binary, (ch) => ch.charCodeAt(0))
  return new TextDecoder().decode(bytes)
}

function resolveGrokImportConfig(
  clientType: CcSwitchClientType,
  baseUrl: string
): CcSwitchImportConfig {
  const model = GROK_CC_SWITCH_MODEL

  if (clientType === 'claude') {
    return {
      app: 'claude',
      endpoint: stripV1Endpoint(baseUrl),
      model,
      haikuModel: model,
      sonnetModel: model,
      opusModel: model,
      configJson: JSON.stringify({ env: grokClaudeCodeEnv(model) })
    }
  }

  if (clientType === 'codex') {
    return {
      app: 'codex',
      endpoint: withV1Endpoint(baseUrl),
      model
    }
  }

  return {
    app: 'grokbuild',
    endpoint: withV1Endpoint(baseUrl),
    model
  }
}

export function resolveCcSwitchImportConfig(
  platform: GroupPlatform | undefined | null,
  clientType: CcSwitchClientType,
  baseUrl: string
): CcSwitchImportConfig {
  switch (platform || 'anthropic') {
    case 'antigravity':
      return {
        app: clientType === 'gemini' ? 'gemini' : 'claude',
        endpoint: `${baseUrl.replace(/\/+$/, '')}/antigravity`
      }
    case 'openai':
      return {
        app: 'codex',
        // CC Switch's Codex provider appends the OpenAI-compatible path itself.
        // Passing /v1 here can make the client request /v1/v1/....
        endpoint: withoutTrailingSlashes(baseUrl),
        model: OPENAI_CC_SWITCH_CODEX_MODEL
      }
    case 'gemini':
      return {
        app: 'gemini',
        endpoint: baseUrl
      }
    case 'grok':
      return resolveGrokImportConfig(clientType, baseUrl)
    default:
      return {
        app: 'claude',
        endpoint: baseUrl
      }
  }
}

export function buildCcSwitchImportDeeplink(input: CcSwitchImportDeeplinkInput): string {
  const config = resolveCcSwitchImportConfig(input.platform, input.clientType, input.baseUrl)
  const entries: [string, string][] = [
    ['resource', 'provider'],
    ['app', config.app],
    ['name', input.providerName],
    ['homepage', input.baseUrl],
    ['endpoint', config.endpoint],
    ['apiKey', input.apiKey],
    ['configFormat', 'json'],
    ['usageEnabled', 'true'],
    ['usageScript', btoa(input.usageScript)],
    ['usageAutoInterval', '30']
  ]

  if (config.model) {
    entries.splice(2, 0, ['model', config.model])
  }
  if (config.haikuModel) {
    entries.push(['haikuModel', config.haikuModel])
  }
  if (config.sonnetModel) {
    entries.push(['sonnetModel', config.sonnetModel])
  }
  if (config.opusModel) {
    entries.push(['opusModel', config.opusModel])
  }
  if (config.configJson) {
    entries.push(['config', encodeCcSwitchConfig(config.configJson)])
  }

  return `ccswitch://v1/import?${new URLSearchParams(entries).toString()}`
}
