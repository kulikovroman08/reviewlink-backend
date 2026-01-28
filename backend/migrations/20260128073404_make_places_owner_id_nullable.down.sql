DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM places WHERE owner_id IS NULL) THEN
    RAISE EXCEPTION 'Cannot set NOT NULL: some places.owner_id are NULL';
END IF;
END $$;

ALTER TABLE places
    ALTER COLUMN owner_id SET NOT NULL;
