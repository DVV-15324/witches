package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return path
}

func GetMigrationsPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return filepath.Join(pwd, "migrate", "migrations")
}

func GetFrameworkPath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return filepath.Dir(exe)
}

func GetMigrationsURL(pathM string) string {
	pwd, _ := os.Getwd()

	// Xử lý các trường hợp
	switch {
	case pathM == "":
		// Mặc định: migrate/migrations
		path := filepath.Join(pwd, "migrate", "migrations")
		os.MkdirAll(path, 0755)
		return filepath.ToSlash(path)

	case !strings.ContainsAny(pathM, "/\\."):
		// Domain: internal/{domain}/migrate/migrations
		path := filepath.Join(pwd, "internal", pathM, "migrate", "migrations")
		os.MkdirAll(path, 0755)
		return filepath.ToSlash(path)

	default:
		// Đường dẫn tùy chỉnh
		path := filepath.Join(pwd, pathM)
		os.MkdirAll(path, 0755)
		return filepath.ToSlash(path)
	}
}
