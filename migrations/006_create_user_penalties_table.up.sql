CREATE TABLE IF NOT EXISTS user_penalties
(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    penalty_until TIMESTAMP NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);