package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

var statWitchesRollback = os.Stat
var getwdWitchesRollback = os.Getwd
var linkWitchesRollback = templates.RollbackDomain

func WitchesRollback(domainName string) error {
	_, err := statWitchesRollback("go.mod")
	if os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}
	if err != nil {
		return fmt.Errorf("stat go.mod: %w", err)
	}

	domainPath := filepath.Join("internal", domainName)
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		return fmt.Errorf("domain '%s' not found", domainName)
	}

	projectPath, err := getwdWitchesRollback()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName := filepath.Base(projectPath)
	if err := linkWitchesRollback(projectPath, moduleName, domainName); err != nil {
		return fmt.Errorf("rollback domain: %w", err)
	}
	return nil
}
