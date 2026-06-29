package cmd_utils

import (
	"log"
	"os"
	"path/filepath"
)

// En: Get the current path of the user's created project.
// Vi: Lấy đường dẫn hiện tại của dự án đã tạo của người dùng
func GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// En: Get the path to migrate/migrations/ from the user side.
// Vi: Lấy đường dẫn của migrate/migrations/ của bên người dùng
func GetMigrationsPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(pwd, "migrate", "migrations")
}

// En: Get the path to the framework (where the binary witches.exe is located)
// Vi: Lấy đường dẫn đến framework (nơi chứa binary witches.exe)
func GetFrameworkPath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exe)
}
