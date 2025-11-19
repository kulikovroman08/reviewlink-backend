CREATE TABLE bonus_rewards
(
    id              UUID PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users (id),
    place_id        UUID        NOT NULL REFERENCES places (id),
    required_points INTEGER     NOT NULL,
    reward_type     VARCHAR(50) NOT NULL,
    qr_token        VARCHAR(20) NOT NULL UNIQUE
);