// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: job.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createJob = `-- name: CreateJob :exec
INSERT INTO jobs (user_id, timestamp) VALUES (?, CURRENT_TIMESTAMP)
`

func (q *Queries) CreateJob(ctx context.Context, userID string) error {
	_, err := q.exec(ctx, q.createJobStmt, createJob, userID)
	return err
}

const getUserJobDetails = `-- name: GetUserJobDetails :many
SELECT jobs.timestamp, SUM(media.new_size) as sum_new_size, SUM(media.old_size) as sum_old_size, COUNT(media.done) as count_done FROM jobs
LEFT JOIN media ON jobs.user_id = media.user_id AND jobs.timestamp = media.timestamp
 WHERE jobs.user_id = ?
 ORDER BY jobs.timestamp DESC
`

type GetUserJobDetailsRow struct {
	Timestamp  int64           `json:"timestamp"`
	SumNewSize sql.NullFloat64 `json:"sum_new_size"`
	SumOldSize sql.NullFloat64 `json:"sum_old_size"`
	CountDone  int64           `json:"count_done"`
}

func (q *Queries) GetUserJobDetails(ctx context.Context, userID string) ([]GetUserJobDetailsRow, error) {
	rows, err := q.query(ctx, q.getUserJobDetailsStmt, getUserJobDetails, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserJobDetailsRow
	for rows.Next() {
		var i GetUserJobDetailsRow
		if err := rows.Scan(
			&i.Timestamp,
			&i.SumNewSize,
			&i.SumOldSize,
			&i.CountDone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
