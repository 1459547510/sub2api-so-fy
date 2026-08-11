# Leo Model Sync

Sub2 administrators can use **Sync latest supported models** on a Leo API-key account. Sub2 requests `GET {base_url}/models` with the account Bearer API key and imports model IDs from the OpenAI-compatible `data` array returned by LeoStudio.

The account `base_url` must remain a valid Leo `/v1` URL. Model-list requests use the existing upstream URL policy, account proxy, TLS fingerprint profile, timeout, response-size limit, header overrides, and error redaction. Upstream response bodies and credentials are not returned to the administrator UI.

## Verification

For an existing Leo account, syncing must cause both of these requests:

- `POST /api/v1/admin/accounts/:id/models/sync-upstream`
- `GET {leo_base_url}/models`

The response must contain the model IDs exposed by LeoStudio `/v1/models`, and the account model mapping is updated only after the administrator saves the form.

## Rollback

Revert the Leo changes in `backend/internal/service/upstream_models.go` and `frontend/src/components/account/ModelWhitelistSelector.vue`, then rebuild and redeploy Sub2API. Existing account mappings and credentials require no data rollback.
