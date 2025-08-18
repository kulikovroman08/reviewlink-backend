CREATE TABLE review_tokens
(
    id          UUID PRIMARY KEY,
    place_id    UUID        NOT NULL REFERENCES places (id),
    token_value VARCHAR(20) NOT NULL UNIQUE,
    is_used     BOOLEAN     NOT NULL DEFAULT false,
    expires_at  TIMESTAMP   NOT NULL
);