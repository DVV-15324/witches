package cmd_migrate

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	utils "github.com/DVV-15324/witches/cmd/utils"
)

func WitchesMigrateUp(DB_URL string, DB_DRIVER string) {
	migratePath := utils.GetMigrationsURL()
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	fmt.Println(migratePath, fullDBURL)
	cmd := exec.Command(
		"migrate",
		"-path", migratePath,
		"-database", fullDBURL,
		"up", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
