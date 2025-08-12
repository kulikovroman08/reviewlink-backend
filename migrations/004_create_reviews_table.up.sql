CREATE TABLE reviews
(
    id         UUID PRIMARY KEY,
    user_id    UUID         NOT NULL REFERENCES users (id),
    place_id   UUID         NOT NULL REFERENCES places (id),
    token_id   UUID         NOT NULL REFERENCES review_tokens (id),
    content    VARCHAR(500) NOT NULL,
    rating     INTEGER      NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP    NOT NULL DEFAULT now()
);