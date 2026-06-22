-- Token incentive plan: weekly self-claim reward for high token usage.
CREATE TABLE IF NOT EXISTS token_incentive_claims (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start    TIMESTAMPTZ NOT NULL,
    week_end      TIMESTAMPTZ NOT NULL,
    tokens        BIGINT NOT NULL DEFAULT 0,
    reward_amount DECIMAL(20, 8) NOT NULL DEFAULT 10,
    status        VARCHAR(20) NOT NULL DEFAULT 'claimed',
    claimed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT token_incentive_claims_user_week_uq UNIQUE (user_id, week_start)
);

CREATE INDEX IF NOT EXISTS idx_token_incentive_claims_user_claimed_at
    ON token_incentive_claims(user_id, claimed_at DESC);

INSERT INTO settings (key, value, updated_at)
VALUES ('token_incentive_enabled', 'false', NOW())
ON CONFLICT (key) DO NOTHING;
