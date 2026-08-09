package utils

import (
	"log"
	"os"
	"path/filepath"
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

func GetMigrationsURL() string {
	pwd, _ := os.Getwd()
	path := filepath.ToSlash(filepath.Join(pwd, "migrate", "migrations"))
	return path
}
