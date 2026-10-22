package utils

import (
	"fmt"
)

func BuildDatabaseURL(cfg *Config) string {
	switch cfg.DBDriver {
	case "postgres", "postgresql":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	case "sqlserver", "mssql":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	default:
		return ""
	}
}
