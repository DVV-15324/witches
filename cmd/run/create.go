package cmd

import (
	"log"
	"os"
	"path/filepath"

	content "github.com/DVV-15324/witches/cmd/cmd_utils"
)

// En: This function creates a new project with environment configuration
// Vi: Hàm tạo dự án mới với file cấu hình môi trường
func WitchesCreate(project string, projectType string, DB_DRIVER string) {
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
	if projectType == "refresh" {
		content := content.CreateContentRefresh(DB_DRIVER)
		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		content := content.CreateContentAccess(DB_DRIVER)
		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}
	}

}
