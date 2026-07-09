# Grok / xAI support

This fork supports xAI Grok OAuth accounts through the OpenAI-compatible gateway paths.

## Text/API models

Default supported Grok text/API models include:

- `grok-4.5`
- `grok-4.3`
- `grok-build-0.1`
- `grok-composer-2.5-fast`
- `grok-4.20-0309-reasoning`
- `grok-4.20-0309-non-reasoning`
- `grok-4.20-multi-agent-0309`

Aliases added by this fork:

- `grok-4.5-latest` -> `grok-4.5`
- `grok-build-latest` -> `grok-4.5`
- `grok-build` -> `grok-build-0.1`
- `grok-composer` -> `grok-composer-2.5-fast`

For compatibility, existing aliases `grok` and `grok-latest` still map to `grok-4.3` until the administrator changes account/group mappings explicitly.

## Pricing fallback

Fallback billing for `grok-4.5` follows the current xAI Pricing page: `$2.00 / 1M input tokens`, `$0.50 / 1M cached input tokens`, and `$6.00 / 1M output tokens`.
