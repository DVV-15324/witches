package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesAdd(moduleName string, DBdriver string) error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}

	domainPath := filepath.Join("internal", moduleName)
	if _, err := os.Stat(domainPath); err == nil {
		return fmt.Errorf("domain '%s' already exists", moduleName)
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %e", err)
	}

	projectName := filepath.Base(projectPath)
	if err := templates.AddModule(projectPath, projectName, moduleName, DBdriver); err != nil {
		return fmt.Errorf("add module: %e", err)
	}
	return nil
}
