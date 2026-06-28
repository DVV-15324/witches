package cmd

import (
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesInit(DB_URL string) {
	projectPath, _ := os.Getwd()
	//Chạy template
	//Kiểm tra thư mục hiện tại và tên project
	moduleName := filepath.Base(projectPath)
	templates.TemplatesCreate(moduleName)
}
