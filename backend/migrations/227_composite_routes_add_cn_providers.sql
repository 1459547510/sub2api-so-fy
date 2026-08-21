-- Allow Composite model routes to target both upstream CN providers and the
-- fork media providers. Keeping the union here prevents existing media route
-- rows from blocking this migration before the follow-up compatibility pass.
ALTER TABLE composite_model_routes
    DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check;

ALTER TABLE composite_model_routes
    ADD CONSTRAINT composite_model_routes_target_platform_check
    CHECK (target_platform IN ('anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                               'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek'));
