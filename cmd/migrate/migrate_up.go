package cmd_migrate

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	utils "github.com/DVV-15324/witches/cmd/utils"
)

// En: Apply 1 pending function with database driver support
// Vi: Chức năng triển khai Migrate với hỗ trợ nhiều loại database
func WitchesMigrateUp(DB_URL string, DB_DRIVER string) {
	// En: Get the path to the migrate/migrations/ folder
	// Vi: Lấy đường dẫn đến folder migrate/migrations/
	migratePath := utils.GetMigrationsURL()

	// En: Build full database URL based on driver
	// Vi: Xây dựng URL đầy đủ dựa trên driver
	fullDBURL := utils.BuildDatabaseURL(DB_DRIVER, DB_URL)
	fmt.Println(migratePath, fullDBURL)
	// En: Start executing using local migrate binary
	// Vi: Bắt đầu thực thi sử dụng binary migrate local
	cmd := exec.Command(
		"migrate", // Sử dụng migrate binary đã cài đặt trên hệ thống
		"-path", migratePath,
		"-database", fullDBURL,
		"up", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
