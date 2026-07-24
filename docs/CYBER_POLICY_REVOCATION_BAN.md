# Cyber Policy Revocation Ban

When an OpenAI Pro or ProLite OAuth account is permanently marked as `Token revoked (401)`, the Ops error pipeline checks whether the same `account_id` received upstream `cyber_policy` errors during the preceding 30 days. The tier is read from `accounts.credentials.plan_type`; account names and billing multipliers are not used.

The user with the most matching hits is disabled. If multiple users have the same hit count, the user with the latest hit is selected; a remaining tie is resolved by the lowest user ID. Administrator accounts are excluded.

Only `ops_error_logs.error_type = 'cyber_policy'` is counted. Local `cyber_policy_session_blocked` rejections never reached the upstream provider and are not included. Generic 401 responses, `Unauthorized` responses without a revoked-token result, OpenAI Plus/Team/Free accounts, non-OpenAI accounts, and accounts that were not persisted in `error` status do not trigger this rule.

For each new upstream `cyber_policy` hit, the gateway stores the current turn's extracted text in `content_moderation_logs.input_excerpt` for the Responses, Chat Completions, Messages, and Responses WebSocket paths. Images and base64 payloads are omitted, credential-like values are redacted, and the retained text is limited to 4,000 Unicode characters. Existing request metadata remains available with the record. Historical rows that were stored without an excerpt cannot be reconstructed and remain empty.

The action runs after the triggering Ops error is persisted. It therefore requires Ops monitoring and the corresponding historical error logs to be available. Repeated observations are idempotent while the selected user remains disabled.

After a user is disabled, the rule synchronously appends an audit entry with action `security.cyber_policy_revocation_ban` to the existing `audit_logs` store. The entry records the revoked account, credential account, plan type, attributed user, hit count, latest matching Ops error ID, request IDs, API key ID and masked prefix, client IP, user agent, model, request path, revocation request IDs, and timestamps.

When the matching content-moderation record is available, its already-redacted input excerpt is copied into the audit entry's redacted request body. Full raw prompts and credentials are not added. Administrators can review the rule event and the retained request summary from **Security Audit > Audit Logs**.
