CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    name          VARCHAR(100)        NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    role          VARCHAR(20)         NOT NULL,
    points        INT       DEFAULT 0,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);