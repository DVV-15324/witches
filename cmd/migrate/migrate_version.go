package cmd_migrate

import (
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateVersion(DB_URL string, DB_DRIVER string) {
	migratePath := utils.GetMigrationsURL()
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	cmd := exec.Command(
		"migrate",
		"-path", migratePath,
		"-database", fullDBURL,
		"version",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
