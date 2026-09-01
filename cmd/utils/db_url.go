package utils

import (
	"fmt"
)

func BuildDatabaseURL(driver, dbURL string) string {
	if driver == "" {
		return dbURL
	}

	switch driver {
	case "mysql":
		return fmt.Sprintf("mysql://%s", dbURL)
	case "postgres", "postgresql":
		return fmt.Sprintf("postgres://%s", dbURL)
	case "sqlserver", "mssql":
		return fmt.Sprintf("sqlserver://%s", dbURL)
	default:
		return fmt.Sprintf("%s://%s", driver, dbURL)
	}
}
