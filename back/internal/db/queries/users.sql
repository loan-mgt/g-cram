-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByTokenHash :one
SELECT * FROM users WHERE token_hash = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, token, token_hash, subscription) VALUES (?, ?, ?, "");

-- name: UpdateUserSubscription :exec
UPDATE users SET subscription = ? WHERE id = ?;

-- name: UpdateUserToken :exec
UPDATE users SET token = ?, token_hash = ? WHERE id = ?;