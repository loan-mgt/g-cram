-- name: CreateMedia :exec
INSERT INTO media (user_id, timestamp, media_id, creation_date, filename, old_size, new_size, done) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: SetMediaDone :exec
UPDATE media SET done = ?, new_size = ? WHERE media_id = ?;

-- name: RemoveMedia :exec
DELETE FROM media WHERE media_id = ?;

-- name: ClearUserTmpMedia :exec
DELETE FROM media WHERE user_id = ? and timestamp = NULL;

-- name: GetMedias :many
SELECT * FROM media WHERE user_id = ? and timestamp = ?;

-- name: SetMediaTimestamp :exec
UPDATE media SET timestamp = ? WHERE media_id = ? and user_id = ?;