package cmd_utils

import (
	"log"
	"os"
	"path/filepath"
)

// Lấy đường dẫn của migrate/migrations/ của bên người dùng
func GetMigrationsPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(pwd, "migrate", "migrations")
}

// Lấy đường dẫn đến framework (nơi chứa binary witches.exe)
// Lấy thông tin và đường dẫn đến tài nguyên của framework để tái sử dụng
func GetFrameworkPath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exe)
}
