DROP INDEX IF EXISTS places_source_source_id_uq;

ALTER TABLE places
DROP COLUMN IF EXISTS source_id,
DROP COLUMN IF EXISTS source;
