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
		"github.com/mailru/easyjson/jwriter",
		"github.com/mailru/easyjson/jlexer",
	}
	drivers := map[string]string{
		"mysql":      "github.com/golang-migrate/migrate/v4/database/mysql@latest",
		"postgres":   "github.comgolang-migrate/migrate/v4/database/postgres@latest",
		"postgresql": "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"mssql":      "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		"sqlserver":  "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
	}
	runCmd("go", "mod", "tidy")
	if driverPath, ok := drivers[DB_DRIVER]; ok {
		runCmd("go", "get", driverPath)
	}
	for _, tool := range tools {
		runCmd("go", "install", tool)
	}
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
