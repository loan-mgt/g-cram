-- +migrate Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    token TEXT,
    token_hash TEXT,
    subscription TEXT
);
CREATE INDEX idx_users_id ON users(id);
CREATE INDEX idx_users_token_hash ON users(token_hash);

CREATE TABLE jobs (
    session_id TEXT,
    user_id TEXT,
    timestamp INTEGER
);
CREATE INDEX idx_jobs_session_id ON jobs(session_id);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);

CREATE TABLE media (
    session_id TEXT,
    media_id TEXT,
    creation_date INTEGER,
    filename TEXT,
    old_size INTEGER,
    new_size INTEGER,
    done INTEGER
);
CREATE INDEX idx_media_session_id ON media(session_id);

-- +migrate Down
DROP INDEX idx_media_session_id;
DROP TABLE media;
DROP INDEX idx_jobs_session_id;
DROP INDEX idx_jobs_user_id;
DROP TABLE jobs;
DROP INDEX idx_users_id;
DROP TABLE users;

