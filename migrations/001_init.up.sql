CREATE TABLE pull_requests (
    pull_request_id    TEXT PRIMARY KEY,
    author_id          TEXT NOT NULL,
    pull_request_name  TEXT NOT NULL,
    status             TEXT NOT NULL,
    assigned_reviewers TEXT[] NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at          TIMESTAMPTZ
);

CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id   TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_name TEXT NOT NULL REFERENCES teams(team_name)
);