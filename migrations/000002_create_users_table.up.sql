CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    user_id    TEXT NOT NULL UNIQUE,
    username   TEXT NOT NULL,
    team_id    INT NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);