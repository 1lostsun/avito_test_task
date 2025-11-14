CREATE TABLE pull_requests (
    id                SERIAL PRIMARY KEY,
    pull_request_id   TEXT NOT NULL UNIQUE,
    pull_request_name TEXT NOT NULL,
    author_id         INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status            pr_status NOT NULL DEFAULT 'OPEN',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at         TIMESTAMPTZ
);

