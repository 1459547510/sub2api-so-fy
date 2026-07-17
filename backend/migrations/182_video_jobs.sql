CREATE TABLE IF NOT EXISTS video_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    api_key_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    upstream_job_id BIGINT,
    status VARCHAR(32) NOT NULL,
    requested_model VARCHAR(128) NOT NULL,
    upstream_model VARCHAR(128) NOT NULL,
    prompt TEXT NOT NULL,
    resolution VARCHAR(32) NOT NULL,
    duration_seconds INT NOT NULL,
    aspect_ratio VARCHAR(32) NOT NULL,
    audio BOOLEAN NOT NULL DEFAULT FALSE,
    image_source VARCHAR(16) NOT NULL DEFAULT 'none',
    image_url TEXT,
    local_input_name VARCHAR(255),
    result JSONB,
    error_message TEXT,
    hold_amount DECIMAL(20,10),
    actual_cost DECIMAL(20,10),
    billing_snapshot JSONB,
    request_hash VARCHAR(128) NOT NULL,
    settled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_video_jobs_job_id
    ON video_jobs (job_id);
CREATE INDEX IF NOT EXISTS idx_video_jobs_user_created
    ON video_jobs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_video_jobs_api_key_created
    ON video_jobs (api_key_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_video_jobs_status_updated
    ON video_jobs (status, updated_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_video_jobs_account_upstream_job
    ON video_jobs (account_id, upstream_job_id)
    WHERE upstream_job_id IS NOT NULL;
