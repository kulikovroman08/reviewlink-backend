DROP INDEX IF EXISTS ux_places_source_source_id;

CREATE UNIQUE INDEX ux_places_source_source_id
    ON places (source, source_id);
