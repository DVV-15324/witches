package cmd

import (
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

// En: This function creates template for project
// Vi : Hàm tạo templates cho dự án
func WitchesInit(DB_URL string) {
	projectPath, _ := os.Getwd()
	//Vi: Kiểm tra thư mục hiện tại và tên dự án
	//En: Check the current directory and project name
	moduleName := filepath.Base(projectPath)
	// n: Execute template
	//Vi: Thực thi template
	templates.TemplatesCreate(moduleName)
}
