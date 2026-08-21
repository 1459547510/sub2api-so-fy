# Composite Groups

Composite groups are an admin routing layer for API keys that should choose a
concrete provider from the requested model instead of binding the key to a
single provider group. They support both built-in model detection and an
admin-configured model route registry for public model aliases.

## Supported Providers

Composite groups can route to these concrete account platforms:

- Anthropic
- Gemini
- OpenAI
- Antigravity
- Grok
- Leo (video and image generation)
- OpenAI Media (self-hosted OpenAI-compatible image/video pools)
- Kimi
- Zhipu GLM
- DeepSeek

The selected concrete platform is used for account selection, user platform
quota checks, post-usage billing, ops error platform attribution, channel
mapping/pricing lookup, and platform usage reporting.

## Route Registry

Admins can configure routes on a composite group from the group list's
`Routes` action or through the admin API:

- `GET /api/v1/admin/groups/:id/composite-routes`
- `POST /api/v1/admin/groups/:id/composite-routes`
- `PUT /api/v1/admin/groups/:id/composite-routes/:route_id`
- `DELETE /api/v1/admin/groups/:id/composite-routes/:route_id`
- `POST /api/v1/admin/groups/:id/composite-routes/preview`

Each route belongs to one composite group and contains:

- `public_model`: model identifier the client sends.
- `match_type`: `exact` or `prefix`.
- `target_platform`: concrete provider platform.
- `upstream_model`: model identifier sent upstream. If omitted, the public
  model is reused.
- `endpoint`: `any`, `messages`, `count_tokens`, `responses`,
  `chat_completions`, `embeddings`, `images`, or `gemini`. Video requests use
  the `any` route because their legacy endpoint is `/v1/videos/generations`.
- `priority`: lower values win after match specificity.
- `enabled`: disabled routes are ignored by runtime resolution but remain
  visible to admins.

Resolution order is explicit route first, then built-in detection. When more
than one explicit route matches, exact matches beat prefix matches,
endpoint-specific routes beat `any`, longer prefixes beat shorter prefixes,
then lower `priority`, then lower route id.

For JSON-body endpoints, the gateway rewrites the request `model` field to the
route's `upstream_model` before dispatch. For Gemini native paths such as
`/v1beta/models/{model}:generateContent`, the gateway resolves `{model}` and
the handler forwards the resolved upstream model.

Codex Alpha Search and Live requests use the `responses` route domain. Live
requests resolve the model from `session.model`, including multipart `session`
payloads, and apply the configured `upstream_model` before dispatch.
Codex model manifest requests reuse the existing OpenAI account selection and
failover path within the Composite group.

## Built-In Detection

Composite routing detects common public model IDs and provider-prefixed IDs:

- `claude-*` and `anthropic/claude-*` route to Anthropic.
- `gemini-*` and `google/gemini-*` route to Gemini.
- `gpt-*`, `o*`, `codex-*`, `text-embedding-*`, `dall-e-*`, and
  `openai/*` route to OpenAI.
- `grok-*` and `xai/grok-*` route to Grok.

Unknown or ambiguous model names fail closed with a client error instead of
guessing a provider.

## Admin Workflows

- Admins can create a group with platform `composite`.
- Admins can add, edit, delete, and preview composite model routes.
- Composite groups can copy accounts from concrete provider groups.
- Concrete provider accounts can be assigned directly to composite groups from
  account create/edit and bulk account workflows.
- Subscription payment plans can bind to a composite group when that group's
  `subscription_type` is `subscription`. The plan grants access to the
  composite group; each request is still billed and quota-checked against the
  resolved concrete provider platform.
- Channel configuration exposes composite groups in concrete provider sections.
  The channel `group_ids` payload is still flat; provider-specific model
  mapping and pricing remain keyed by concrete platform.

## Bucket 2 Setup: OpenAI + Claude + Gemini + Grok

Use one composite subscription group when one customer-facing plan should expose
model aliases across OpenAI, Claude, Gemini, and Grok without issuing separate
keys per provider.

1. Create concrete provider groups for the upstream account pools, for example
   `OpenAI Paid`, `Claude Paid`, `Gemini Paid`, and `Grok Paid`.
2. Create a `composite` group with `subscription_type` set to `subscription`.
3. Assign provider accounts directly to the composite group, or copy accounts
   from the concrete provider groups during group creation.
4. Add explicit routes for public aliases that should not rely on built-in
   model detection:

   | Public model | Endpoint | Target platform | Upstream model |
   | --- | --- | --- | --- |
   | `all/gpt-5` | `responses` | `openai` | `gpt-5` |
   | `all/claude-sonnet` | `messages` | `anthropic` | `claude-sonnet-4-6` |
   | `all/gemini-pro` | `gemini` | `gemini` | `gemini-2.5-pro` |
   | `all/grok` | `responses` | `grok` | `grok-4.3` |

5. Configure channel pricing and model mapping under the concrete platforms
   named in each route. Composite routing does not create pricing records. Leo
   media routes must explicitly target `leo` because names such as
   `gemini-omni-flash` and `grok-imagine-1.5` can also refer to other provider
   families.
6. Create a subscription payment plan for the composite group.

The same composite group can also rely on built-in detection for standard model
names such as `gpt-*`, `claude-*`, `gemini-*`, and `grok-*`. Explicit routes are
recommended for bundled plan aliases because they make endpoint, provider, and
upstream model attribution reviewable in the admin UI.

## Leo + OpenAI Media In One Group

The professional media group can contain both account types at the same time:

| Account pool | Stored platform | Authentication | Upstream contract |
| --- | --- | --- | --- |
| Existing LeoStudio accounts | `leo` | API key + `/v1` Base URL | Existing Leo video contract |
| New self-hosted media pools | `openai_media` | API key + `/v1` Base URL | OpenAI-compatible images/videos |

Example account credentials for a new pool:

```json
{
  "platform": "openai_media",
  "type": "apikey",
  "credentials": {
    "base_url": "https://pool-a.example.com/v1",
    "api_key": "YOUR_UPSTREAM_API_KEY",
    "model_mapping": {
      "seedance-2.0": "seedance-2.0"
    }
  }
}
```

Add explicit routes to the same composite group, for example:

| Public model | Endpoint | Target platform |
| --- | --- | --- |
| `seedance-2.0` | `any` | `leo` |
| `happy-horse-1.1` | `any` | `openai_media` |

The client keeps the existing OpenAI-compatible endpoint and request body. The
target platform controls account selection and billing; it does not alter the
public API key or require a client migration.

## Limits

Composite routes choose a concrete provider and upstream model; they do not
create synthetic model metadata, pricing, or upstream capability records by
themselves. Keep channel pricing/model mapping configured for the concrete
provider platforms that the routes target.

Legacy Leo clients remain compatible when their API key is moved to a
composite group. The existing `/v1/videos/generations`, `/v1/videos/jobs/*`,
and `/v1/videos/uploads` request and response contracts remain unchanged for
Leo-routed models. Leo async job IDs are generated with the `vidjob_` prefix;
when a client polls a generic legacy status or content path, that prefix sends
the lookup back to the Leo job service. Unprefixed external request IDs retain
the existing external media lookup path (currently Grok); a new asynchronous
provider must register its task-ID binding before being enabled in this group.

For a professional video/image composite group, keep the customer's old
`Authorization` header and endpoint paths. Add explicit routes for every
ambiguous media model, for example `seedance-2.0` -> `leo` and
`grok-imagine-1.5` -> `grok`; the public model name and request body do not need
to change on the client side.

LeoStudio accounts remain in this same composite group with
`platform=leo`. New self-hosted media pools use `platform=openai_media` and
the same API-key fields (`base_url`, `api_key`, and `model_mapping`). Add both
route targets to the group; do not migrate or rename the existing LeoStudio
accounts. A route targeting `leo` selects only the original LeoStudio pool,
while a route targeting `openai_media` selects only the new OpenAI-compatible
pool.

This PR intentionally does not implement:

- AUTO smart-routing among multiple providers for the same abstract task.
- Direct API-key binding to several existing groups without a composite group.
- Protocol-agnostic provider decoupling or a LiteLLM-style adapter rewrite.
