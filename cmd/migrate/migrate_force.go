package cmd_migrate

import (
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateForce(DB_URL string, DB_DRIVER string, migrationPath string, VERSION string) {

	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	cmd := exec.Command(
		"migrate",
		"-path", migrationPath,
		"-database", fullDBURL,
		"force", VERSION,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
