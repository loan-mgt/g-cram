// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"
)

type Job struct {
	SessionID sql.NullString `json:"session_id"`
	UserID    sql.NullString `json:"user_id"`
	Timestamp sql.NullInt64  `json:"timestamp"`
}

type Medium struct {
	SessionID    sql.NullString `json:"session_id"`
	MediaID      sql.NullString `json:"media_id"`
	CreationDate sql.NullInt64  `json:"creation_date"`
	Filename     sql.NullString `json:"filename"`
	OldSize      sql.NullInt64  `json:"old_size"`
	NewSize      sql.NullInt64  `json:"new_size"`
	Done         sql.NullInt64  `json:"done"`
}

type User struct {
	ID    string         `json:"id"`
	Token sql.NullString `json:"token"`
}
