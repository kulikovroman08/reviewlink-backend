CREATE TABLE IF NOT EXISTS user_restrictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    restriction_type VARCHAR(50) NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL
);

ALTER TABLE user_restrictions
    ADD CONSTRAINT uniq_user_restriction UNIQUE (user_id, restriction_type);

CREATE INDEX IF NOT EXISTS idx_user_restrictions_user_type
    ON user_restrictions(user_id, restriction_type);