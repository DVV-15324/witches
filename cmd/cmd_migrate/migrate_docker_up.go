package cmd_migrate

import (
	utils "core-v/cmd/cmd_utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func WitchesMigrateDockerUp(DB_PROFILE string) {
	// Lấy đường dẫn tài nguyên của framework đang lưu trong máy của người dùng
	// Tái sử dụng ở máy người dùng
	frameworkPath := utils.GetFrameworkPath()
	dockerPath := filepath.Join(frameworkPath, "internal", "core", "database")

	// Kiểm tra thư mục tồn tại
	if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
		log.Fatalf("Docker config not found at: %s\nPlease ensure witches is installed correctly.", dockerPath)
	}

	// Chạy docker khởi tạo docker migrate
	cmd := exec.Command("docker", "compose",
		"--profile", "migrate",
		"up", "--build", "-d")
	cmd.Dir = dockerPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
