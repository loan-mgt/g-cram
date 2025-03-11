-- name: CreateMedia :exec
INSERT INTO media (session_id, media_id, creation_date, filename, old_size, new_size, done) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: SetMediaDone :exec
UPDATE media SET done = ?, new_size = ? WHERE media_id = ?;
