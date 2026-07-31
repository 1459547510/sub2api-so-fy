# OpenAI Codex reverse-proxy fingerprint

## Purpose

The gateway keeps OpenAI OAuth traffic consistent with a long-lived Codex client without fabricating account subscription, quota, payment, or server-side trust state.

The implementation combines two useful properties:

- CPA-style deterministic identity isolation for sessions and synthetic client identifiers.
- Sub2API's existing upstream-account, proxy, HTTP/2, and connection-pool ownership.

## Existing and new accounts

The rollout is deliberately split by account generation:

- Accounts created before this mode was introduced and loaded without an explicit mode remain `legacy` at runtime. They keep the previous device, window, session, conversation, and prompt-cache behavior; the compatibility marker is not written back during a read.
- New OpenAI OAuth accounts created through the normal account creation services receive `accounts.extra.openai_fingerprint_mode = "v1"` and use the fingerprint rules in this document.
- An explicit `"legacy"` or `"v1"` value in `accounts.extra.openai_fingerprint_mode` overrides the default. This allows a staged per-account migration without changing the database schema.

The `legacy`/`v1` split applies only to OpenAI OAuth accounts. OpenAI API-key accounts and other platforms are unchanged.

Proxy addresses, proxy credentials, access tokens, and refresh tokens are never part of a device seed. Token refreshes and normal proxy address rotation therefore do not create a new device.

## Device identity precedence

For OpenAI OAuth accounts, the outbound installation identity is resolved in this order:

1. `accounts.extra.openai_device_id`: an operator-managed installation ID.
2. An inbound `X-Codex-Installation-ID` or body `client_metadata.x-codex-installation-id`.
3. A deterministic installation ID derived from the upstream ChatGPT account identity.

In `v1`, an inbound identity is deterministically mapped with the selected
upstream ChatGPT account identity before it is sent upstream. This keeps one
downstream installation stable for that account while preventing the same raw
installation value from being presented to multiple upstream accounts. Its
window ID is mapped into the same account namespace. In `legacy`, inbound
installation and window IDs retain the previous pass-through behavior.

Set `accounts.extra.openai_device_profile_id` to a new stable value when an intentional device rotation is required. Changing this value rotates derived installation, window, session, conversation, and prompt-cache namespaces as one generation.

`openai_device_profile_id` is evaluated only for accounts using fingerprint mode `v1`.

Example account `extra` values:

```json
{
  "openai_device_id": "operator-captured-installation-id",
  "openai_device_profile_id": "workstation-2"
}
```

`openai_device_id` takes precedence. `openai_device_profile_id` affects deterministic mapping and allows a controlled generation change without using proxy identity.

## Identifier lifetimes

| Identifier | Lifetime and scope |
| --- | --- |
| Installation ID | Stable for the account/device profile; `v1` maps inbound values per upstream account, while `legacy` preserves pass-through |
| Window ID | `v1` maps the downstream window per upstream account and API key; `legacy` preserves pass-through |
| Session and conversation ID | Stable for the conversation and isolated by upstream account plus downstream API key |
| Prompt cache key | Uses the same account and downstream-tenant namespace as the session |
| Turn and request IDs | Retain their request/turn lifecycle; they are not converted into long-lived device identifiers |

The account namespace prefers `chatgpt_account_id`, falls back to the parent account for shadow accounts, and finally uses the local account ID. Access and refresh tokens are deliberately ignored.

## Covered outbound paths

The same fingerprint helpers are used by:

- Responses HTTP forwarding and passthrough.
- Chat Completions and Messages bridges that target Codex Responses.
- Responses WebSocket v2 and WebSocket ingress.
- Alpha/search and Live.
- Account connection, compact capability, and usage probes.

OpenAI API-key upstreams are unchanged. The mapping applies only to OpenAI OAuth traffic sent to ChatGPT internal APIs.

## Transport boundary

OpenAI requests continue to use the project's OpenAI transport profile, which prefers HTTP/2 and has proxy compatibility fallback. The existing configurable custom TLS fingerprint is a Node.js/Claude profile and remains limited to Anthropic accounts; applying it to a `codex_cli_rs` identity would create an inconsistent fingerprint.

A future Codex-specific TLS profile must be based on a packet capture from the exact official Codex client version and must align User-Agent, Originator, version, ALPN, and HTTP/2 behavior as one bundle.

## Operational checks

After changing device settings, verify:

1. Repeated requests for one account and conversation keep the same mapped identifiers.
2. Token refresh and proxy address rotation do not change the installation ID.
3. Different upstream accounts or downstream API keys do not share session/cache identifiers.
4. HTTP and WebSocket requests expose the same installation/window namespace.
5. Changing `openai_device_profile_id` rotates the complete managed namespace once.

## Account failover boundary

Codex can report a plan/model incompatibility as either an HTTP 400 response or
an HTTP 200 `response.failed` event. The message
`model is not supported when using Codex with a ChatGPT account` is treated as
a deterministic request error: the current account receives the error and the
same request is not fanned out across the OAuth pool. Capacity and transient
processing errors remain eligible for their existing bounded retry/failover
policy.
