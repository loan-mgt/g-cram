// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"
)

type Job struct {
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type Medium struct {
	UserID       string         `json:"user_id"`
	Timestamp    sql.NullInt64  `json:"timestamp"`
	MediaID      string         `json:"media_id"`
	CreationDate sql.NullInt64  `json:"creation_date"`
	Filename     sql.NullString `json:"filename"`
	OldSize      sql.NullInt64  `json:"old_size"`
	NewSize      sql.NullInt64  `json:"new_size"`
	Done         sql.NullInt64  `json:"done"`
}

type User struct {
	ID           string         `json:"id"`
	Token        sql.NullString `json:"token"`
	TokenHash    sql.NullString `json:"token_hash"`
	Subscription sql.NullString `json:"subscription"`
}
