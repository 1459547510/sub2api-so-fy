-- Add TypeSafe (Jev System One) as a first-class platform
-- without dropping fork media platforms (leo / openai_media).
--
-- Upstream 241 rebuilt the CHECKs without those media values. Databases that
-- already have leo / openai_media quota or composite-route rows then fail
-- ADD CONSTRAINT and abort startup, the same class of failure as 224 / 237 / 238.
--
-- Same shape as 238_opencode_go: rebuild only when the constraint is missing
-- leo, openai_media, or typesafe. The rebuilt list is the 13-platform union.
--
-- TypeSafe is not a chat model and does not enter channel-monitor providers,
-- so channel_monitors / channel_monitor_request_templates stay unchanged.
--
-- 1. user_platform_quotas.platform CHECK
-- 2. composite_model_routes.target_platform CHECK
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
       OR position('typesafe' IN quota_constraint_def) = 0 THEN
        ALTER TABLE user_platform_quotas
            DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check;
        ALTER TABLE user_platform_quotas
            ADD CONSTRAINT user_platform_quotas_platform_check
            CHECK (platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go', 'typesafe'
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
       OR position('typesafe' IN route_constraint_def) = 0 THEN
        ALTER TABLE composite_model_routes
            DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check;
        ALTER TABLE composite_model_routes
            ADD CONSTRAINT composite_model_routes_target_platform_check
            CHECK (target_platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go', 'typesafe'
            ));
    END IF;
END $$;
