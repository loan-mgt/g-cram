package db

import (
	"database/sql"
	"log"

	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*sqlc.Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(cfg *config.Config) *Store {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &Store{
		db:      db,
		Queries: sqlc.New(db),
	}
}

// Close closes the database connection
func (s *Store) Close() {
	s.db.Close()
}

// ExecTx executes a function within a database transaction
func (s *Store) ExecTx(fn func(*sqlc.Queries) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
