CREATE TABLE teams (
    id         SERIAL PRIMARY KEY,
    team_name  TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);