-- Re-assert the full fork platform set for databases that may already have
-- applied the original upstream 227 migration before switching to this build.
ALTER TABLE composite_model_routes
    DROP CONSTRAINT IF EXISTS composite_model_routes_target_platform_check;

ALTER TABLE composite_model_routes
    ADD CONSTRAINT composite_model_routes_target_platform_check
    CHECK (target_platform IN ('anthropic', 'openai', 'gemini', 'antigravity', 'grok',
                               'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek'));
