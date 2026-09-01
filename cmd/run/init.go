package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesInit(DBdriver string) error {
	if _, err := os.Stat("witches.env"); os.IsNotExist(err) {
		return fmt.Errorf("witches.env not found")
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName := filepath.Base(projectPath)
	if err := templates.CreateTemplateGoArc(moduleName, DBdriver); err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}
