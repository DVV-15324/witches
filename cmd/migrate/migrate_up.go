package cmd_migrate

import (
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
)

// En: Apply 1 pending funciton
// Vi: Chức năng triển khai Migrate
func WitchesMigrateUp(DB_URL string) {
	//En: Get the path to the migrate/migrations/ folder
	//Vi: Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsPath()

	//En: Start executing
	//Vi: Bắt đầu thực thi
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", migratePath+":/migrations",
		"--network", "host",
		"migrate/migrate",
		"-path=/migrations",
		"-database", DB_URL,
		"up", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
