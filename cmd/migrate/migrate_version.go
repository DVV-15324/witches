package cmd_migrate

import (
	utils "github.com/DVV-15324/witches/cmd/utils"
	"log"
	"os"
	"os/exec"
)

// En: Show current migration version
// Vi: Hiển thị phiên bản migration hiện tại
func WitchesMigrateVersion(DB_URL string, DB_DRIVER string) {
	// En: Get the path to the migrate/migrations/ folder
	// Vi: Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsURL()

	// En: Build full database URL based on driver
	// Vi: Xây dựng URL đầy đủ dựa trên driver
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)

	// En: Start executing using local migrate binary
	// Vi: Bắt đầu thực thi sử dụng binary migrate local
	cmd := exec.Command(
		"migrate", // Sử dụng migrate binary đã cài đặt trên hệ thống
		"-path", migratePath,
		"-database", fullDBURL,
		"version",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
