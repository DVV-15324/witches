package cmd

import (
	"log"
	"os"
	"os/exec"
)

// En: This function installs dependencies
// Vi: Hàm cài đặt các thư viện cần dùng
func WitchesInstall(DB_DRIVER string) {
	tools := []string{
		"github.com/mailru/easyjson/...@latest",
	}
	dbTools := map[string][]string{
		"mysql": {
			"github.com/golang-migrate/migrate/v4/database/mysql@latest",
		},
		"postgres": {
			"github.com/golang-migrate/migrate/v4/database/postgres@latest",
		},
		"postgresql": {
			"github.com/golang-migrate/migrate/v4/database/postgres@latest",
		},
		"mssql": {
			"github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		},
		"sqlserver": {
			"github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		},
	}
	if driverTools, ok := dbTools[DB_DRIVER]; ok {
		tools = append(tools, driverTools...)
	}
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error: failed to install  %v", err)
	}
	for _, tool := range tools {
		cmd := exec.Command("go", "install", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Error: failed to install %s: %v", tool, err)
		}
	}
}
