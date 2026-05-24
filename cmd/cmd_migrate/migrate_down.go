package cmd_migrate

import (
	utils "core-v/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
)

func WitchesMigrateDown(DB_URL string) {
	// Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsPath()
	// copy file khi init của người dùng migrate/migrations -> /migrations của docker
	// --rm: Chạy xong xóa
	// ảnh migrate/migrate có sẵn trên docker hub
	// up 1
	cmd := exec.Command("docker", "run", "--rm",
		"-v", migratePath+":/migrations",
		"--network", "host",
		"migrate/migrate",
		"-path=/migrations",
		"-database", DB_URL,
		"down", "1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
