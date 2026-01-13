CREATE TABLE review_votes (
                              review_id  uuid      NOT NULL,
                              user_id    uuid      NOT NULL,
                              value      int2      NOT NULL,
                              created_at timestamp NOT NULL DEFAULT now(),
                              updated_at timestamp NOT NULL DEFAULT now(),
                              CONSTRAINT review_votes_value_check CHECK (value IN (1, -1)),
                              CONSTRAINT review_votes_review_fk FOREIGN KEY (review_id)
                                  REFERENCES reviews(id) ON DELETE CASCADE,
                              CONSTRAINT review_votes_user_fk FOREIGN KEY (user_id)
                                  REFERENCES users(id) ON DELETE CASCADE,
                              CONSTRAINT review_votes_unique UNIQUE (review_id, user_id)
);

CREATE INDEX idx_review_votes_review_id ON review_votes(review_id);
CREATE INDEX idx_review_votes_user_id   ON review_votes(user_id);
