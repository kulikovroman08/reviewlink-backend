ALTER TABLE places
    ADD COLUMN source TEXT,
  ADD COLUMN source_id TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS places_source_source_id_uq
    ON places (source, source_id)
    WHERE source IS NOT NULL AND source_id IS NOT NULL;
