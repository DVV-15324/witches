package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesRollback(domainName string) error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}

	domainPath := filepath.Join("internal", domainName)
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		return fmt.Errorf("domain '%s' not found", domainName)
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName := filepath.Base(projectPath)
	if err := templates.RollbackDomain(projectPath, moduleName, domainName); err != nil {
		return fmt.Errorf("rollback domain: %w", err)
	}
	return nil
}
