package cmd

import (
	"log"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesRollback(domainName string) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	moduleName := filepath.Base(projectPath)
	_ = templates.RollbackDomain(projectPath, moduleName, domainName)
}
