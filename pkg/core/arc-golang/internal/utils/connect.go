package utils

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5"
)

func ConnectPostgreSQL(DB_URL string) (*sql.DB, error) {
	db, err := sql.Open(
		"pgx",
		DB_URL,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
