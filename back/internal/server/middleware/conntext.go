package middleware

import (
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
)

type MiddleWareContext struct {
	db  *db.Store
	cfg *config.Config
}

func NewMiddleWareContextr(db *db.Store, cfg *config.Config) *MiddleWareContext {
	return &MiddleWareContext{db: db, cfg: cfg}
}
