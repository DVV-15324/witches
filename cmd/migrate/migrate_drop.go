package cmd_migrate

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/utils"
	"os"
	"os/exec"
)

func WitchesMigrateDrop(cfg *utils.Config, migrationPath string) error {
	fullDBURL := utils.BuildDatabaseURL(cfg)
	cmd := exec.Command(
		"migrate",
		"-path", migrationPath,
		"-database", fullDBURL,
		"drop", "-f",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("migrate drop failed: %w", err)
	}
	return nil
}
