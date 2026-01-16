ALTER TABLE reviews
ADD COLUMN helpful_count int4 NOT NULL DEFAULT 0,
ADD COLUMN unhelpful_count int4 NOT NULL DEFAULT 0;
