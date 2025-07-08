CREATE TABLE places
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(100),
    address    TEXT,
    owner_id   UUID      NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);