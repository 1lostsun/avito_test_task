CREATE TABLE pull_requests (
    id                SERIAL PRIMARY KEY,
    pull_request_id   TEXT NOT NULL UNIQUE,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    status            TEXT NOT NULL DEFAULT 'OPEN',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at         TIMESTAMPTZ,
    CONSTRAINT chk_status CHECK (status IN ('OPEN', 'MERGED'))
);

