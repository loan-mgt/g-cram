-- name: CreateJob :exec
INSERT INTO jobs (session_id, user_id, timestamp) VALUES (?, ?, CURRENT_TIMESTAMP);

-- name: GetUserJobDetails :many
SELECT * FROM jobs
LEFT JOIN media ON jobs.session_id = media.session_id
 WHERE user_id = ?
