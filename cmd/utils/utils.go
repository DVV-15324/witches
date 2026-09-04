package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var executableUtils = os.Executable
var getwdUtils = os.Getwd

func GetCurrentPath() string {
	path, err := getwdUtils()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return path
}

func GetMigrationsPath() string {
	pwd, err := getwdUtils()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return filepath.Join(pwd, "migrate", "migrations")
}

func GetFrameworkPath() string {
	exe, err := executableUtils()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return filepath.Dir(exe)
}

func GetMigrationsURL(pathM string) string {
	pwd, _ := getwdUtils()

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
