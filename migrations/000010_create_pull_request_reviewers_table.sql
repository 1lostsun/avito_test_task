CREATE TABLE pull_request_reviewers (
    pr_id    INT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id  INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (pr_id, user_id)
);