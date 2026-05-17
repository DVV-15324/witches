package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
)

func ConnectPostgreSQL(userName string, passWord string, url string, post int, dataBase string) (*sql.DB, error) {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", userName, passWord, url, post, dataBase),
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
