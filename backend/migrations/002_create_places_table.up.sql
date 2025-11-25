CREATE TABLE places
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(100),
    address    TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted  BOOLEAN  NOT NULL DEFAULT false
);