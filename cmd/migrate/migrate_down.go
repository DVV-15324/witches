package cmd_migrate

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"os"
	"os/exec"
)

func WitchesMigrateDown(DB_URL string, DB_DRIVER string, migrationPath string) error {
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	cmd := exec.Command(
		"migrate",
		"-path", migrationPath,
		"-database", fullDBURL,
		"down", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("migrate down failed: %w", err)
	}
	return nil
}
