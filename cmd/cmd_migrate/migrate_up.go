package cmd_migrate

import (
	"fmt"
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateUp(DB_URL string) {
	// Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsPath()

	fmt.Println("MIGRATION PATH:", migratePath)

	// Print tất cả file trong folder
	files, err := os.ReadDir(migratePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("FILES:")

	for _, file := range files {
		fmt.Println("-", file.Name())
	}

	// Chạy migrate
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", migratePath+":/migrations",
		"--network", "host",
		"migrate/migrate",
		"-path=/migrations",
		"-database", DB_URL,
		"up", "1",
	)
	fmt.Println("FULL DOCKER COMMAND:")
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
