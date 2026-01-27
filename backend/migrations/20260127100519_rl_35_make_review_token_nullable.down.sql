DELETE FROM reviews
WHERE token_id IS NULL;

ALTER TABLE reviews
    ALTER COLUMN token_id SET NOT NULL;