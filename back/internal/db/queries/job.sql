-- name: CreateJob :exec
INSERT INTO jobs (user_id, timestamp) VALUES (?, CURRENT_TIMESTAMP);

-- name: GetUserJobDetails :many
SELECT jobs.timestamp, SUM(media.new_size) as sum_new_size, SUM(media.old_size) as sum_old_size, COUNT(media.done) as count_done FROM jobs
LEFT JOIN media ON jobs.user_id = media.user_id AND jobs.timestamp = media.timestamp
 WHERE jobs.user_id = ?
 ORDER BY jobs.timestamp DESC;
