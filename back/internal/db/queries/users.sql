-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, token, subscription) VALUES (?, ?, "");

-- name: UpdateUserSubscription :exec
UPDATE users SET subscription = ? WHERE id = ?;

-- name: UpdateUserToken :exec
UPDATE users SET token = ? WHERE id = ?;