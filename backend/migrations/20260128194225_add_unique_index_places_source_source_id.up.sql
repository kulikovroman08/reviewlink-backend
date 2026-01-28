CREATE UNIQUE INDEX IF NOT EXISTS ux_places_source_source_id
    ON places (source, source_id)
    WHERE source IS NOT NULL AND source_id IS NOT NULL;