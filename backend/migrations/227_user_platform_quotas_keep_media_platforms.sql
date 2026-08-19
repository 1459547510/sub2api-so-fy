-- Keep user_platform_quotas.platform aligned with AllowedQuotaPlatforms.
--
-- 224 added kimi/zhipu/deepseek but dropped the fork media platforms
-- (leo / openai_media). That CHECK rejects existing quota rows and aborts
-- startup. This migration restores the full 10-platform list and is
-- idempotent when the constraint already contains both media and CN names.
DO $$
DECLARE
    constraint_def TEXT;
BEGIN
    SELECT pg_get_constraintdef(c.oid)
      INTO constraint_def
      FROM pg_constraint c
      JOIN pg_class t ON t.oid = c.conrelid
     WHERE t.relname = 'user_platform_quotas'
       AND c.conname = 'user_platform_quotas_platform_check';

    IF constraint_def IS NULL
       OR position('leo' IN constraint_def) = 0
       OR position('openai_media' IN constraint_def) = 0
       OR position('kimi' IN constraint_def) = 0 THEN
        ALTER TABLE user_platform_quotas
            DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check;
        ALTER TABLE user_platform_quotas
            ADD CONSTRAINT user_platform_quotas_platform_check
            CHECK (platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek'
            ));
    END IF;
END $$;
