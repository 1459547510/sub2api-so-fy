package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (r *opsRepository) FindCyberPolicyBanCandidate(ctx context.Context, accountID int64, since, until time.Time) (*service.CyberPolicyBanCandidate, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil ops repository")
	}
	if accountID <= 0 || since.IsZero() || until.IsZero() || until.Before(since) {
		return nil, nil
	}

	var candidate service.CyberPolicyBanCandidate
	var apiKeyID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
WITH ranked_users AS (
  SELECT
    e.user_id,
    COUNT(*) AS hit_count,
    MAX(e.created_at) AS last_hit_at
  FROM ops_error_logs e
  JOIN users u ON u.id = e.user_id
  WHERE e.account_id = $1
    AND e.error_type = 'cyber_policy'
    AND e.created_at >= $2
    AND e.created_at <= $3
    AND u.deleted_at IS NULL
    AND u.role <> $4
  GROUP BY e.user_id
  ORDER BY hit_count DESC, last_hit_at DESC, e.user_id ASC
  LIMIT 1
), latest_hit AS (
  SELECT e.*
  FROM ops_error_logs e
  JOIN ranked_users ranked ON ranked.user_id = e.user_id
  WHERE e.account_id = $1
    AND e.error_type = 'cyber_policy'
    AND e.created_at >= $2
    AND e.created_at <= $3
  ORDER BY e.created_at DESC, e.id DESC
  LIMIT 1
)
SELECT
  ranked.user_id,
  ranked.hit_count,
  ranked.last_hit_at,
  hit.id,
  COALESCE(hit.request_id, ''),
  COALESCE(hit.client_request_id, ''),
  hit.api_key_id,
  COALESCE(hit.api_key_prefix, ''),
  COALESCE(hit.client_ip::text, ''),
  COALESCE(hit.user_agent, ''),
  COALESCE(NULLIF(TRIM(hit.requested_model), ''), hit.model, ''),
  COALESCE(NULLIF(TRIM(hit.inbound_endpoint), ''), hit.request_path, ''),
  COALESCE((
    SELECT moderation.input_excerpt
    FROM content_moderation_logs moderation
    WHERE moderation.request_id = hit.request_id
      AND moderation.user_id = hit.user_id
      AND moderation.action = 'cyber_policy'
    ORDER BY moderation.created_at DESC, moderation.id DESC
    LIMIT 1
  ), '')
FROM ranked_users ranked
JOIN latest_hit hit ON hit.user_id = ranked.user_id
`, accountID, since.UTC(), until.UTC(), service.RoleAdmin).Scan(
		&candidate.UserID,
		&candidate.HitCount,
		&candidate.LastHitAt,
		&candidate.OpsErrorLogID,
		&candidate.RequestID,
		&candidate.ClientRequestID,
		&apiKeyID,
		&candidate.APIKeyPrefix,
		&candidate.ClientIP,
		&candidate.UserAgent,
		&candidate.Model,
		&candidate.RequestPath,
		&candidate.InputExcerpt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if apiKeyID.Valid {
		candidate.APIKeyID = &apiKeyID.Int64
	}
	return &candidate, nil
}
