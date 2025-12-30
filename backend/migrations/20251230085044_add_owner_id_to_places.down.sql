DROP INDEX IF EXISTS idx_places_owner_id;

ALTER TABLE places
DROP CONSTRAINT IF EXISTS fk_places_owner;

ALTER TABLE places
DROP COLUMN IF EXISTS owner_id;