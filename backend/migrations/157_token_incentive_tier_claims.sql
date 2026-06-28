-- Token incentive tier claims: allow one claim per reached weekly tier.
ALTER TABLE token_incentive_claims
    ADD COLUMN IF NOT EXISTS threshold_tokens BIGINT NOT NULL DEFAULT 0;

CREATE OR REPLACE FUNCTION pg_temp.token_incentive_safe_jsonb(raw text)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN raw::jsonb;
EXCEPTION WHEN OTHERS THEN
    RETURN '[]'::jsonb;
END;
$$;

WITH configured_rules AS (
    SELECT DISTINCT
           (item->>'threshold_tokens')::bigint AS threshold_tokens,
           (item->>'reward_amount')::numeric AS reward_amount
    FROM settings s
    CROSS JOIN LATERAL jsonb_array_elements(pg_temp.token_incentive_safe_jsonb(s.value)) item
    WHERE s.key = 'token_incentive_rules'
      AND item ? 'threshold_tokens'
      AND item ? 'reward_amount'
      AND (item->>'threshold_tokens') ~ '^[0-9]+$'
      AND (item->>'reward_amount') ~ '^[0-9]+(\.[0-9]+)?$'
      AND (item->>'threshold_tokens')::bigint > 0
      AND (item->>'reward_amount')::numeric > 0
),
default_rules(threshold_tokens, reward_amount) AS (
    VALUES
        (50000000::bigint, 2::numeric),
        (100000000::bigint, 5::numeric),
        (500000000::bigint, 10::numeric)
),
rules AS (
    SELECT threshold_tokens, reward_amount FROM configured_rules
    UNION ALL
    SELECT threshold_tokens, reward_amount
    FROM default_rules
    WHERE NOT EXISTS (SELECT 1 FROM configured_rules)
),
resolved AS (
    SELECT c.id,
           COALESCE(
               (
                   SELECT r.threshold_tokens
                   FROM rules r
                   WHERE c.tokens >= r.threshold_tokens
                     AND c.reward_amount = r.reward_amount
                   ORDER BY r.threshold_tokens DESC
                   LIMIT 1
               ),
               (
                   SELECT r.threshold_tokens
                   FROM rules r
                   WHERE c.tokens >= r.threshold_tokens
                   ORDER BY r.threshold_tokens DESC
                   LIMIT 1
               ),
               c.tokens,
               1
           ) AS threshold_tokens
    FROM token_incentive_claims c
    WHERE c.threshold_tokens = 0
)
UPDATE token_incentive_claims c
SET threshold_tokens = GREATEST(resolved.threshold_tokens, 1)
FROM resolved
WHERE c.id = resolved.id;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'token_incentive_claims_threshold_positive_chk'
          AND conrelid = 'token_incentive_claims'::regclass
    ) THEN
        ALTER TABLE token_incentive_claims
            ADD CONSTRAINT token_incentive_claims_threshold_positive_chk
            CHECK (threshold_tokens > 0);
    END IF;
END $$;

ALTER TABLE token_incentive_claims
    DROP CONSTRAINT IF EXISTS token_incentive_claims_user_week_uq;

CREATE UNIQUE INDEX IF NOT EXISTS token_incentive_claims_user_week_tier_uq
    ON token_incentive_claims(user_id, week_start, threshold_tokens);
