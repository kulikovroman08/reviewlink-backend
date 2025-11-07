ALTER TABLE bonus_rewards
    ADD COLUMN is_used BOOLEAN DEFAULT false,
    ADD COLUMN used_at TIMESTAMP NULL;
