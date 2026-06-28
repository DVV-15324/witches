package auth

import (
	"database/sql"
)

type RepositoryAuth struct {
	db *sql.DB
}

func NewRepositoryAuth(db *sql.DB) *RepositoryAuth {
	return &RepositoryAuth{db: db}
}
