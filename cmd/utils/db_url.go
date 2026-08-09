package utils

import (
	"fmt"

	"strings"
)

func BuildDatabaseURL(driver, dbURL string) string {
	cleanURL := strings.TrimPrefix(dbURL, driver+"://")
	cleanURL = strings.TrimPrefix(cleanURL, "mysql://")
	cleanURL = strings.TrimPrefix(cleanURL, "postgres://")
	cleanURL = strings.TrimPrefix(cleanURL, "postgresql://")
	cleanURL = strings.TrimPrefix(cleanURL, "sqlserver://")
	cleanURL = strings.TrimPrefix(cleanURL, "mssql://")

	switch driver {
	case "mysql":
		return fmt.Sprintf("mysql://%s", cleanURL)

	case "postgres", "postgresql":
		return fmt.Sprintf("postgres://%s", cleanURL)

	case "sqlserver", "mssql":
		return fmt.Sprintf("sqlserver://%s", cleanURL)

	default:
		return fmt.Sprintf("%s://%s", driver, cleanURL)
	}
}
