-- name: CreateMedia :exec
INSERT INTO media (user_id, timestamp, media_id, creation_date, filename, base_url, old_size, new_size, step, done) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?);

-- name: SetMediaDone :exec
UPDATE media SET done = ? WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: RemoveMedia :exec
DELETE FROM media WHERE media_id = ?;

-- name: ClearUserTmpMedia :exec
DELETE FROM media WHERE user_id = ? and timestamp = NULL;

-- name: GetMedias :many
SELECT * FROM media WHERE user_id = ? and step = 0;

-- name: GetMedia :one
SELECT * FROM media WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: SetMediaTimestamp :exec
UPDATE media SET timestamp = ? WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: CountUserMedia :one
SELECT COUNT(*) FROM media WHERE user_id = ? and step = 0;

-- name: GetMediaCurrentStep :one
SELECT step FROM media WHERE media_id = ? and user_id = ?;

-- name: SetMediaStep :exec
UPDATE media SET step = ? WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: SetMediaNewSize :exec
UPDATE media SET new_size = ? WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: SetMediaOldSize :exec
UPDATE media SET old_size = ? WHERE media_id = ? and user_id = ? and timestamp = ?;

-- name: CountUserMediaInJob :one
SELECT COUNT(*) FROM media WHERE user_id = ? and timestamp = ?;

-- name: CountUserMediaInJobAtStep :one
SELECT COUNT(*) FROM media WHERE user_id = ? and timestamp = ? and step >= ?;

-- name: GetJobSpace :one
SELECT SUM(new_size) as sum_new_size, SUM(old_size) as sum_old_size FROM media WHERE user_id = ? and timestamp = ?;

-- name: GetUserJob :many
SELECT COUNT(*) as nb_media, SUM(done) as nb_media_done, timestamp, SUM(old_size) as old_size, SUM(new_size) as new_size FROM media WHERE user_id = ? GROUP BY timestamp ORDER BY timestamp DESC LIMIT 5;
