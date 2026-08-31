package cmd

import (
	"log"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesLink(domainName, repoURL string) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	moduleName := filepath.Base(projectPath)
	templates.AddGoDomainFromLink(projectPath, moduleName, domainName, repoURL)
}
