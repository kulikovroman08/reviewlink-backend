CREATE TABLE reviews
(
    id          UUID PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users (id),
    place_id    UUID        NOT NULL REFERENCES places (id),
    token_id    UUID        NOT NULL REFERENCES review_tokens (id),
    content     TEXT        NOT NULL,
    rating      INTEGER     NOT NULL CHECK (rating >= 1 AND rating <= 5),
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    is_verified BOOLEAN     NOT NULL DEFAULT false,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);