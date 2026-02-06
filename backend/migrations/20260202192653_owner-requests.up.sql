CREATE TABLE owner_requests
(
    id          uuid        NOT NULL PRIMARY KEY,
    user_id     uuid        NOT NULL REFERENCES users (id),

    source      text        NOT NULL,
    source_id   text        NOT NULL,

    status      varchar(20) NOT NULL,
    created_at  timestamp   NOT NULL DEFAULT now(),

    reviewed_at timestamp NULL,
    reviewed_by uuid NULL REFERENCES users(id),
    comment     text NULL,

    CONSTRAINT owner_requests_status_check
        CHECK (status IN ('pending', 'approved', 'rejected'))
);

CREATE INDEX idx_owner_requests_status_created_at
    ON owner_requests (status, created_at DESC);

CREATE INDEX idx_owner_requests_user_id_created_at
    ON owner_requests (user_id, created_at DESC);

CREATE UNIQUE INDEX ux_owner_requests_pending_user_source
    ON owner_requests (user_id, source, source_id) WHERE status = 'pending';
