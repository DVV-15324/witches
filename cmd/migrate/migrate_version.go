package cmd_migrate

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/DVV-15324/witches/cmd/utils"
)

func WitchesMigrateVersion(DB_URL string, DB_DRIVER string, migrationPath string) error {
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	cmd := exec.Command(
		"migrate",
		"-path", migrationPath,
		"-database", fullDBURL,
		"version",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("migrate version failed: %w", err)
	}
	return nil
}
