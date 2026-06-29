package cmd_migrate

import (
	utils "github.com/DVV-15324/witches/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
)

// En: Force set migration version function
// Vi: Hàm thiết lập phiên bản Migrate
func WitchesMigrateForce(DB_URL, VERSION string) {
	//En: Get the path to the migrate/migrations/ folder
	//Vi: Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsPath()
	//En: Start executing
	//Vi: Bắt đầu thực thi
	cmd := exec.Command("docker", "run", "--rm",
		"-v", migratePath+":/migrations",
		"--network", "host",
		"migrate/migrate",
		"-path=/migrations",
		"-database", DB_URL,
		"force", VERSION)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
