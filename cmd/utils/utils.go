package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func GetModuleGoName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("cannot read go.mod: %w", err)
	}

	re := regexp.MustCompile(`(?m)^module\s+([^\s]+)`)
	matches := re.FindStringSubmatch(string(data))

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1]), nil
	}

	return "", fmt.Errorf("module name not found in go.mod")
}
