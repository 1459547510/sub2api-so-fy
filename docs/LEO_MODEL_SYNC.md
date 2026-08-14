# Leo Model Sync

Sub2 administrators can use **Sync latest supported models** on an API-key
account or **Sync upstream models** in channel pricing. Sub2 requests
`GET {base_url}/models` with the account Bearer API key and imports the public
model IDs from the OpenAI-compatible `data` array.

The pricing sync classifies registered video IDs as video billing entries and
all other live catalog entries as image billing entries. Each model gets an
independent pricing row. After saving, an administrator only needs to fill the
video resolution prices or the per-image price. A newly published image model
therefore does not require a code release. A newly published video model still
requires a local capability specification before it can be exposed safely.

The account `base_url` must remain a valid Leo `/v1` URL. Model-list requests use the existing upstream URL policy, account proxy, TLS fingerprint profile, timeout, response-size limit, header overrides, and error redaction. Upstream response bodies and credentials are not returned to the administrator UI.

## Verification

For an existing Leo account, syncing must cause both of these requests:

- `POST /api/v1/admin/accounts/:id/models/sync-upstream`
- `GET {leo_base_url}/models`

The response must contain the public model IDs exposed by `/v1/models`. Raw
asset identifiers, account identifiers, credentials, and upstream ownership
metadata must not be returned to the client.

## Rollback

Revert the catalog changes in `backend/internal/service/upstream_models.go`,
`backend/internal/handler/admin/channel_handler.go`, and the channel pricing
form, then rebuild and redeploy Sub2API. Existing account credentials and
saved prices require no data rollback.
