package utils

import (
	"fmt"
	"net"
	"net/url"
)

func BuildDatabaseURL(cfg *Config) string {
	switch cfg.DBDriver {
	case "postgres", "postgresql":
		u := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
			Host:   net.JoinHostPort(cfg.DBHost, fmt.Sprint(cfg.DBPort)),
			Path:   "/" + cfg.DBName,
		}
		q := u.Query()
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()
		return u.String()

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
