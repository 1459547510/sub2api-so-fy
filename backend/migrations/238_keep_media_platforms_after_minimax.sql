-- 237 added MiniMax but dropped the fork media platforms (leo / openai_media).
-- That CHECK rejects existing quota / composite-route rows and aborts startup,
-- the same class of failure as upstream 224.
--
-- Restore the full 11-platform set and keep MiniMax. Idempotent when the
-- constraint already contains leo, openai_media, and minimax.
DO $$
DECLARE
    quota_constraint_def TEXT;
    route_constraint_def TEXT;
BEGIN
    SELECT pg_get_constraintdef(c.oid)
      INTO quota_constraint_def
      FROM pg_constraint c
      JOIN pg_class t ON t.oid = c.conrelid
     WHERE t.relname = 'user_platform_quotas'
       AND c.conname = 'user_platform_quotas_platform_check';

    IF quota_constraint_def IS NULL
       OR position('leo' IN quota_constraint_def) = 0
       OR position('openai_media' IN quota_constraint_def) = 0
       OR position('minimax' IN quota_constraint_def) = 0 THEN
        ALTER TABLE user_platform_quotas
            DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check;
        ALTER TABLE user_platform_quotas
            ADD CONSTRAINT user_platform_quotas_platform_check
            CHECK (platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax'
            ));
    END IF;

    SELECT pg_get_constraintdef(c.oid)
      INTO route_constraint_def
      FROM pg_constraint c
      JOIN pg_class t ON t.oid = c.conrelid
     WHERE t.relname = 'composite_model_routes'
       AND c.conname = 'composite_model_routes_target_platform_check';

    IF route_constraint_def IS NULL
       OR position('leo' IN route_constraint_def) = 0
       OR position('openai_media' IN route_constraint_def) = 0
       OR position('minimax' IN route_constraint_def) = 0 THEN
        ALTER TABLE composite_model_routes
            DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check;
        ALTER TABLE composite_model_routes
            ADD CONSTRAINT composite_model_routes_target_platform_check
            CHECK (target_platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax'
            ));
    END IF;
END $$;
