-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, token) VALUES (?, ?);

-- name: UpdateUserToken :exec
UPDATE users SET token = ? WHERE id = ?;