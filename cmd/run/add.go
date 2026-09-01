package cmd

import (
	"log"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesAdd(moduleName string, DBdriver string) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	projectName := filepath.Base(projectPath)
	templates.AddModule(projectPath, projectName, moduleName, DBdriver)
}
