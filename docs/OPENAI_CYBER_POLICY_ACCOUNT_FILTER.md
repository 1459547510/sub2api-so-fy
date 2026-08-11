# OpenAI CY User Account Filter

This feature lets an administrator protect selected OpenAI accounts from users
who have ever triggered an upstream `cyber_policy` (CY) event.

## Account setting

Set the account `extra` key below to `true`:

```json
{
  "openai_block_cyber_policy_users": true
}
```

The default is disabled. The admin account editor exposes this as an OpenAI
account-level toggle and preserves all other `extra` keys.

## Marking and scheduling

- An upstream `cyber_policy` event recorded by `RecordCyberPolicyEvent` marks a
  regular user in `cyber_policy_user_marks`.
- Migration `221_cyber_policy_user_marks.sql` backfills historical flagged
  `cyber_policy` rows. Local `cyber_policy_session_blocked` events are not used.
- Administrator users are excluded from both the backfill and future marks.
- User ID `55` (`hjt13845049131@163.com`) is explicitly exempt from account
  filtering even if a historical marker already exists.
- A marked user skips protected accounts before scoring and during final account
  rechecks. Sticky sessions, weighted sticky fallback, previous-response IDs,
  WebSocket turns, model routing, and retry failover therefore fall through to
  another eligible account.
- If every eligible account is protected, the existing no-available-account
  error is returned.

The marker is stored in the database permanently. Redis stores positive markers
without expiry and negative results for five minutes to reduce repeated reads;
cache and database failures fail open and preserve existing account availability.

## Known limitation

New markers are currently written through `RecordCyberPolicyEvent`, which returns
without recording when the global risk-control switch is disabled. Consequently,
an upstream `cyber_policy` event that occurs while global risk control is disabled
will not mark that user for this account filter. The current operating assumption
is that global risk control remains enabled. Before disabling it, either decouple
marker persistence from that switch or accept and separately reconcile the gap.

This feature does not disable users, change account status or scheduler scores,
alter account identity or reverse-proxy fingerprints, or modify billing.

## Rollback

Before deployment, revert the code and documentation changes and remove the
`221_cyber_policy_user_marks.sql` migration from the release. After migration
application, disable the account toggles first; removing the table requires a
separate reviewed migration because the marker data is intentionally durable.
