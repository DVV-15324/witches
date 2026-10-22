package cmd_migrate

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/DVV-15324/witches/cmd/utils"
)

func WitchesMigrateForce(cfg *utils.Config, migrationPath string, VERSION string) error {

	fullDBURL := utils.BuildDatabaseURL(cfg)
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
		return fmt.Errorf("migrate force failed: %w", err)
	}
	return nil
}
