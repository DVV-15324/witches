package utils

import (
	"fmt"

	"strings"
)

// En: Build database URL for different drivers
// Vi: Xây dựng URL database cho các driver khác nhau
func BuildDatabaseURL(driver, dbURL string) string {
	// En: Remove any existing protocol prefix
	// Vi: Xóa bất kỳ prefix protocol nào có sẵn
	cleanURL := strings.TrimPrefix(dbURL, driver+"://")
	cleanURL = strings.TrimPrefix(cleanURL, "mysql://")
	cleanURL = strings.TrimPrefix(cleanURL, "postgres://")
	cleanURL = strings.TrimPrefix(cleanURL, "postgresql://")
	cleanURL = strings.TrimPrefix(cleanURL, "sqlserver://")
	cleanURL = strings.TrimPrefix(cleanURL, "mssql://")

	switch driver {
	case "mysql":
		// MySQL format: mysql://user:pass@tcp(host:port)/dbname?params
		// Ví dụ: mysql://root:password@tcp(localhost:3306)/mydb?charset=utf8
		return fmt.Sprintf("mysql://%s", cleanURL)

	case "postgres", "postgresql":
		// PostgreSQL format: postgres://user:pass@host:port/dbname?params
		// Ví dụ: postgres://user:password@localhost:5432/mydb?sslmode=disable
		return fmt.Sprintf("postgres://%s", cleanURL)

	case "sqlserver", "mssql":
		// SQL Server format: sqlserver://user:pass@host:port?database=dbname&params
		// Ví dụ: sqlserver://sa:password@localhost:1433?database=mydb&encrypt=disable
		return fmt.Sprintf("sqlserver://%s", cleanURL)

	default:
		// En: Default format for other drivers
		// Vi: Định dạng mặc định cho các driver khác
		return fmt.Sprintf("%s://%s", driver, cleanURL)
	}
}
