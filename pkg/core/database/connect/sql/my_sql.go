package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL(user, pass, host string, port int, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, pass, host, port, dbname,
	)

	return sql.Open("mysql", dsn)
}
