package cmd_migrate

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/DVV-15324/witches/cmd/utils"
)

func WitchesMigrateUp(cfg *utils.Config, migrationPath string) error {

	fullDBURL := utils.BuildDatabaseURL(cfg)
	cmd := exec.Command(
		"migrate",
		"-path", migrationPath,
		"-database", fullDBURL,
		"up", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}
	return nil
}
