-- Harden token incentive tier claims for databases upgraded from the one-claim-per-week schema.
ALTER TABLE token_incentive_claims
    ADD COLUMN IF NOT EXISTS threshold_tokens BIGINT NOT NULL DEFAULT 0;

ALTER TABLE token_incentive_claims
    ALTER COLUMN threshold_tokens SET DEFAULT 0;

UPDATE token_incentive_claims
SET threshold_tokens = GREATEST(COALESCE(NULLIF(threshold_tokens, 0), tokens, 1), 1)
WHERE threshold_tokens IS NULL OR threshold_tokens <= 0;

ALTER TABLE token_incentive_claims
    ALTER COLUMN threshold_tokens SET NOT NULL;

DO $$
DECLARE
    constraint_record record;
    index_record record;
BEGIN
    FOR constraint_record IN
        SELECT con.conname
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_namespace ns ON ns.oid = rel.relnamespace
        WHERE ns.nspname = 'public'
          AND rel.relname = 'token_incentive_claims'
          AND con.contype = 'u'
          AND (
              SELECT array_agg(att.attname::text ORDER BY att.attname::text)
              FROM unnest(con.conkey) AS cols(attnum)
              JOIN pg_attribute att ON att.attrelid = rel.oid AND att.attnum = cols.attnum
          ) = ARRAY['user_id', 'week_start']::text[]
    LOOP
        EXECUTE format('ALTER TABLE public.token_incentive_claims DROP CONSTRAINT IF EXISTS %I', constraint_record.conname);
    END LOOP;

    FOR index_record IN
        SELECT idx.relname AS index_name
        FROM pg_index i
        JOIN pg_class idx ON idx.oid = i.indexrelid
        JOIN pg_class rel ON rel.oid = i.indrelid
        JOIN pg_namespace ns ON ns.oid = rel.relnamespace
        WHERE ns.nspname = 'public'
          AND rel.relname = 'token_incentive_claims'
          AND i.indisunique
          AND NOT EXISTS (
              SELECT 1
              FROM pg_constraint con
              WHERE con.conindid = i.indexrelid
          )
          AND (
              SELECT array_agg(att.attname::text ORDER BY att.attname::text)
              FROM unnest(string_to_array(i.indkey::text, ' ')::smallint[]) AS cols(attnum)
              JOIN pg_attribute att ON att.attrelid = rel.oid AND att.attnum = cols.attnum
          ) = ARRAY['user_id', 'week_start']::text[]
    LOOP
        EXECUTE format('DROP INDEX IF EXISTS public.%I', index_record.index_name);
    END LOOP;

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

CREATE UNIQUE INDEX IF NOT EXISTS token_incentive_claims_user_week_tier_uq
    ON token_incentive_claims(user_id, week_start, threshold_tokens);

ALTER TABLE redeem_codes
    ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT NULL;
