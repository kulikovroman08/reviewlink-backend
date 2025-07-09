CREATE TABLE bonus_rewards
(
    id              UUID PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users (id),
    place_id        UUID        NOT NULL REFERENCES places (id),
    points_spent    INTEGER     NOT NULL,
    required_points INTEGER     NOT NULL,
    reward_type     VARCHAR(50) NOT NULL,
    qr_token        VARCHAR(20) NOT NULL UNIQUE,
    is_redeemed     BOOLEAN   DEFAULT false,
    generated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    redeemed_at     TIMESTAMP
);