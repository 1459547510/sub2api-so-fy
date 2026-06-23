-- Token incentive plan configurable weekly tiers.
INSERT INTO settings (key, value, updated_at)
VALUES (
    'token_incentive_rules',
    '[{"threshold_tokens":50000000,"reward_amount":2},{"threshold_tokens":100000000,"reward_amount":5},{"threshold_tokens":500000000,"reward_amount":10}]',
    NOW()
)
ON CONFLICT (key) DO NOTHING;
