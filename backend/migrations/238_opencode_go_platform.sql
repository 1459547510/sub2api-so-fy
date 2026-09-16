-- Add OpenCode as a first-class platform (account types Zen / GO)
-- without dropping fork media platforms (leo / openai_media).
--
-- Upstream 238 rebuilt the CHECKs without those media values. Databases that
-- already have leo / openai_media quota or composite-route rows then fail
-- ADD CONSTRAINT and abort startup, the same class of failure as 224 / 237.
--
-- Same shape as 237/238_keep: rebuild only when the constraint is missing leo,
-- openai_media, or opencode_go. The rebuilt list is the 12-platform union.
--
-- 1. user_platform_quotas.platform CHECK
-- 2. composite_model_routes.target_platform CHECK
-- 3. channel_monitors / channel_monitor_request_templates provider CHECK
DO $$
DECLARE
    quota_constraint_def TEXT;
    route_constraint_def TEXT;
    monitor_constraint_def TEXT;
    template_constraint_def TEXT;
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
       OR position('opencode_go' IN quota_constraint_def) = 0 THEN
        ALTER TABLE user_platform_quotas
            DROP CONSTRAINT IF EXISTS user_platform_quotas_platform_check;
        ALTER TABLE user_platform_quotas
            ADD CONSTRAINT user_platform_quotas_platform_check
            CHECK (platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go'
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
       OR position('opencode_go' IN route_constraint_def) = 0 THEN
        ALTER TABLE composite_model_routes
            DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check;
        ALTER TABLE composite_model_routes
            ADD CONSTRAINT composite_model_routes_target_platform_check
            CHECK (target_platform IN (
                'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go'
            ));
    END IF;

    SELECT pg_get_constraintdef(c.oid)
      INTO monitor_constraint_def
      FROM pg_constraint c
      JOIN pg_class t ON t.oid = c.conrelid
     WHERE t.relname = 'channel_monitors'
       AND c.conname = 'channel_monitors_provider_check';

    IF monitor_constraint_def IS NULL OR position('opencode_go' IN monitor_constraint_def) = 0 THEN
        ALTER TABLE channel_monitors
            DROP CONSTRAINT IF EXISTS channel_monitors_provider_check;
        ALTER TABLE channel_monitors
            ADD CONSTRAINT channel_monitors_provider_check
            CHECK (provider IN ('openai', 'anthropic', 'gemini', 'grok',
                                'antigravity', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go'));
    END IF;

    SELECT pg_get_constraintdef(c.oid)
      INTO template_constraint_def
      FROM pg_constraint c
      JOIN pg_class t ON t.oid = c.conrelid
     WHERE t.relname = 'channel_monitor_request_templates'
       AND c.conname = 'channel_monitor_request_templates_provider_check';

    IF template_constraint_def IS NULL OR position('opencode_go' IN template_constraint_def) = 0 THEN
        ALTER TABLE channel_monitor_request_templates
            DROP CONSTRAINT IF EXISTS channel_monitor_request_templates_provider_check;
        ALTER TABLE channel_monitor_request_templates
            ADD CONSTRAINT channel_monitor_request_templates_provider_check
            CHECK (provider IN ('openai', 'anthropic', 'gemini', 'grok',
                                'antigravity', 'kimi', 'zhipu', 'deepseek', 'minimax', 'opencode_go'));
    END IF;
END $$;
