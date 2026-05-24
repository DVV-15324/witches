package cmd_migrate

import (
	utils "core-v/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateVersion(DB_URL string) {
	migratePath := utils.GetMigrationsPath()
	cmd := exec.Command("docker", "run", "--rm",
		"-v", migratePath+":/migrations",
		"--network", "host",
		"migrate/migrate",
		"-path=/migrations",
		"-database", DB_URL,
		"version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
