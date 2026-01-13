ALTER TABLE reviews
DROP COLUMN IF EXISTS helpful_count,
DROP COLUMN IF EXISTS unhelpful_count;
