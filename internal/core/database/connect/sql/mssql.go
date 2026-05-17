package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectMSSQL(user, pass, host string, port int, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		user, pass, host, port, dbname,
	)

	return sql.Open("sqlserver", dsn)
}
