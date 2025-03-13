// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createJobStmt, err = db.PrepareContext(ctx, createJob); err != nil {
		return nil, fmt.Errorf("error preparing query CreateJob: %w", err)
	}
	if q.createMediaStmt, err = db.PrepareContext(ctx, createMedia); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMedia: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserJobDetailsStmt, err = db.PrepareContext(ctx, getUserJobDetails); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserJobDetails: %w", err)
	}
	if q.setMediaDoneStmt, err = db.PrepareContext(ctx, setMediaDone); err != nil {
		return nil, fmt.Errorf("error preparing query SetMediaDone: %w", err)
	}
	if q.updateUserTokenStmt, err = db.PrepareContext(ctx, updateUserToken); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserToken: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createJobStmt != nil {
		if cerr := q.createJobStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createJobStmt: %w", cerr)
		}
	}
	if q.createMediaStmt != nil {
		if cerr := q.createMediaStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMediaStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserJobDetailsStmt != nil {
		if cerr := q.getUserJobDetailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserJobDetailsStmt: %w", cerr)
		}
	}
	if q.setMediaDoneStmt != nil {
		if cerr := q.setMediaDoneStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setMediaDoneStmt: %w", cerr)
		}
	}
	if q.updateUserTokenStmt != nil {
		if cerr := q.updateUserTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserTokenStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                    DBTX
	tx                    *sql.Tx
	createJobStmt         *sql.Stmt
	createMediaStmt       *sql.Stmt
	createUserStmt        *sql.Stmt
	getUserStmt           *sql.Stmt
	getUserJobDetailsStmt *sql.Stmt
	setMediaDoneStmt      *sql.Stmt
	updateUserTokenStmt   *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                    tx,
		tx:                    tx,
		createJobStmt:         q.createJobStmt,
		createMediaStmt:       q.createMediaStmt,
		createUserStmt:        q.createUserStmt,
		getUserStmt:           q.getUserStmt,
		getUserJobDetailsStmt: q.getUserJobDetailsStmt,
		setMediaDoneStmt:      q.setMediaDoneStmt,
		updateUserTokenStmt:   q.updateUserTokenStmt,
	}
}
