UPDATE places
SET owner_id = (
    SELECT id FROM users
    WHERE role = 'admin'
    ORDER BY created_at
    LIMIT 1
    )
WHERE owner_id IS NULL;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM places WHERE owner_id IS NULL) THEN
    RAISE EXCEPTION 'Cannot set NOT NULL: some places.owner_id are NULL (no admin to backfill?)';
END IF;
END $$;

ALTER TABLE places
    ALTER COLUMN owner_id SET NOT NULL;
