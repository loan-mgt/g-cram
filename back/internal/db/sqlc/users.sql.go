// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, token, websocket) VALUES (?, ?, ?)
`

type CreateUserParams struct {
	ID        string         `json:"id"`
	Token     sql.NullString `json:"token"`
	Websocket sql.NullString `json:"websocket"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.exec(ctx, q.createUserStmt, createUser, arg.ID, arg.Token, arg.Websocket)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, token, websocket FROM users WHERE id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Token, &i.Websocket)
	return i, err
}

const updateUserToken = `-- name: UpdateUserToken :exec
UPDATE users SET token = ? WHERE id = ?
`

type UpdateUserTokenParams struct {
	Token sql.NullString `json:"token"`
	ID    string         `json:"id"`
}

func (q *Queries) UpdateUserToken(ctx context.Context, arg UpdateUserTokenParams) error {
	_, err := q.exec(ctx, q.updateUserTokenStmt, updateUserToken, arg.Token, arg.ID)
	return err
}

const updateUserWebsocket = `-- name: UpdateUserWebsocket :exec
UPDATE users SET websocket = ? WHERE id = ?
`

type UpdateUserWebsocketParams struct {
	Websocket sql.NullString `json:"websocket"`
	ID        string         `json:"id"`
}

func (q *Queries) UpdateUserWebsocket(ctx context.Context, arg UpdateUserWebsocketParams) error {
	_, err := q.exec(ctx, q.updateUserWebsocketStmt, updateUserWebsocket, arg.Websocket, arg.ID)
	return err
}
