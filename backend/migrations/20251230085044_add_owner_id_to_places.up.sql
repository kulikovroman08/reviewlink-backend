ALTER TABLE places
    ADD COLUMN owner_id UUID;

ALTER TABLE places
    ADD CONSTRAINT fk_places_owner
        FOREIGN KEY (owner_id) REFERENCES users(id);

CREATE INDEX idx_places_owner_id ON places(owner_id);