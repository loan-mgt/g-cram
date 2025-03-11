-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, token, websocket) VALUES (?, ?, ?);

-- name: UpdateUserTokenAndWebsocket :exec
UPDATE users SET token = ?, websocket = ? WHERE id = ?;