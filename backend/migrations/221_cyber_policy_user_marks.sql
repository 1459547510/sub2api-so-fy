CREATE TABLE IF NOT EXISTS cyber_policy_user_marks (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO cyber_policy_user_marks (user_id, first_triggered_at, last_triggered_at)
SELECT l.user_id, MIN(l.created_at), MAX(l.created_at)
FROM content_moderation_logs l
JOIN users u ON u.id = l.user_id
WHERE l.user_id IS NOT NULL
  AND l.flagged = TRUE
  AND l.action = 'cyber_policy'
  AND u.role <> 'admin'
GROUP BY l.user_id
ON CONFLICT (user_id) DO UPDATE SET
    first_triggered_at = LEAST(cyber_policy_user_marks.first_triggered_at, EXCLUDED.first_triggered_at),
    last_triggered_at = GREATEST(cyber_policy_user_marks.last_triggered_at, EXCLUDED.last_triggered_at);
