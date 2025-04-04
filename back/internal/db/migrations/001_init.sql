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
    user_id TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    PRIMARY KEY (user_id, timestamp),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_user_id_timestamp ON jobs(user_id, timestamp);

CREATE TABLE media (
    user_id TEXT NOT NULL,
    timestamp INTEGER,
    media_id TEXT NOT NULL,
    creation_date INTEGER,
    filename TEXT,
    old_size INTEGER,
    new_size INTEGER,
    done INTEGER,
    PRIMARY KEY (user_id, timestamp, media_id),
    FOREIGN KEY (user_id, timestamp) REFERENCES jobs(user_id, timestamp)
);
CREATE INDEX idx_media_user_id ON media(user_id);

-- +migrate Down
DROP INDEX idx_media_user_id;
DROP TABLE media;
DROP INDEX idx_jobs_user_id_timestamp;
DROP INDEX idx_jobs_user_id;
DROP TABLE jobs;
DROP INDEX idx_users_token_hash;
DROP INDEX idx_users_id;
DROP TABLE users;

