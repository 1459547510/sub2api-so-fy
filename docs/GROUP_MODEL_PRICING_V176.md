# v0.1.176 Group Model Pricing

## Deployment impact

Release `v0.1.176-fy.1` adds two columns to the PostgreSQL `groups` table through
`backend/migrations/221_group_model_pricing.sql`:

- `long_context_pricing_enabled BOOLEAN NOT NULL DEFAULT TRUE`
- `model_pricing JSONB`

The application applies pending SQL migrations during normal startup. Deploy the
new binary and restart the service once; no separate manual SQL command is
required. A failed migration prevents startup and must be resolved before the
service can run the new version.

## Compatibility

Existing groups keep long-context pricing enabled because the migration defaults
and backfills `long_context_pricing_enabled` to `TRUE`. An empty `model_pricing`
value keeps the existing resolution order: channel pricing, then built-in model
pricing. Configured group model pricing takes priority over both.

Leo video group pricing retains the fork's model-specific resolution and
per-second validation. Other platforms can use the upstream generic video
default price or tier pricing.

The fork also contains `221_cyber_policy_user_marks.sql`. Migration identity is
the complete filename, so both files are intentional and neither existing
migration should be renamed.

## Verification

After restart, verify the service is active and confirm both columns exist:

```sql
SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_name = 'groups'
  AND column_name IN ('long_context_pricing_enabled', 'model_pricing');
```

## Rollback

Reverting the application binary does not require dropping these nullable or
backward-compatible columns. Leave them in place for a normal rollback. Dropping
them would destroy configured group pricing and requires a separate, explicitly
approved database change.
