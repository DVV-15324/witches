package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesLink(domainName, repoURL string) error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName := filepath.Base(projectPath)
	if err := templates.AddGoDomainFromLink(projectPath, moduleName, domainName, repoURL); err != nil {
		return fmt.Errorf("link domain: %w", err)
	}
	return nil
}
