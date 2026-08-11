package cmd

import (
	"log"
	"os"
	"os/exec"
)

func WitchesInstall(DB_DRIVER string) {
	tools := []string{
		"github.com/mailru/easyjson/easyjson@latest",
		"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
	}
	drivers := map[string]string{
		"mysql":      "github.com/golang-migrate/migrate/v4/database/mysql@latest",
		"postgres":   "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"postgresql": "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"mssql":      "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		"sqlserver":  "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
	}

	for _, tool := range tools {
		log.Printf("  Installing %s...", tool)
		runCmd("go", "install", tool)
	}

	if driverPath, ok := drivers[DB_DRIVER]; ok {
		log.Printf("  Adding driver: %s", DB_DRIVER)
		runCmd("go", "get", driverPath)
	}

	runCmd("go", "get", "github.com/mailru/easyjson@latest")

	runCmd("go", "mod", "tidy")

	log.Println("Installation complete!")
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("%s %v", name, args)
	}
}
