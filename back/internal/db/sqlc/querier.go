// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateJob(ctx context.Context, arg CreateJobParams) error
	CreateMedia(ctx context.Context, arg CreateMediaParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	GetUser(ctx context.Context, id string) (User, error)
	GetUserJobDetails(ctx context.Context, userID sql.NullString) ([]GetUserJobDetailsRow, error)
	SetMediaDone(ctx context.Context, arg SetMediaDoneParams) error
	UpdateUserToken(ctx context.Context, arg UpdateUserTokenParams) error
	UpdateUserWebsocket(ctx context.Context, arg UpdateUserWebsocketParams) error
}

var _ Querier = (*Queries)(nil)
