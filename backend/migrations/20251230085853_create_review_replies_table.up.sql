CREATE TABLE review_replies (
                                id         UUID PRIMARY KEY,
                                review_id  UUID NOT NULL UNIQUE,
                                admin_id   UUID NOT NULL,
                                content    VARCHAR(1000) NOT NULL,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP,

                                CONSTRAINT fk_reply_review FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
                                CONSTRAINT fk_reply_admin  FOREIGN KEY (admin_id)  REFERENCES users(id)
);

CREATE INDEX idx_review_replies_admin_id ON review_replies(admin_id);
