CREATE TABLE IF NOT EXISTS moderation_actions
(
    id UUID PRIMARY KEY,
    review_id UUID NOT NULL REFERENCES reviews(id),
    moderator_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);