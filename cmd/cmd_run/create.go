package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// En: This function creates a new project with environment configuration
// Vi: Hàm tạo dự án mới với file cấu hình môi trường
func WitchesCreate(project string, DB_DRIVER string) {
	projectPath := filepath.Join(".", project)
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//En: Create witches.env
	//Vi: Tạo witches.env
	envPath := filepath.Join(projectPath, "witches.env")

	//En: O_CREATE:Permission to create a file when it doesn't exist
	//En: O_WRONLY:Write permissions only -> avoids making random edits during development.
	//Vi: O_CREATE: Quyền tạo khi file không tồn tại
	//Vi: O_WRONLY: Chỉ quyền ghi -> tránh trường hợp trong lúc phát triến sửa lung tung
	file, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//En: Make sure the file is closed
	//Vi: Chắc chắn file đóng
	defer file.Close()
	_, err = file.WriteString(
		"APP_PORT=YOUR_PORT_APP\n" +
			"DB_PASSWORD=YOUR_PASSWORD\n" +
			"DB_NAME=YOUR_DATABASE\n" +
			"DB_HOST=DB_HOST\n" +
			"DB_PORT=YOUR_PORT_APP\n" +
			"REDIS_PORT=REDIS_PORT\n" +
			fmt.Sprintf("DB_DRIVER=%s\n", DB_DRIVER),
	)
	if err != nil {
		log.Fatal(err)
	}
}
