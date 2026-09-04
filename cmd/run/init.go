package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

var statWitchesInit = os.Stat
var getwdWitchesInit = os.Getwd
var initWitchesInit = templates.CreateTemplateGoArc

func WitchesInit(DBdriver string) error {
	_, err := statWitchesInit("witches.env")

	if os.IsNotExist(err) {
		return fmt.Errorf("witches.env not found")
	}

	if err != nil {
		return fmt.Errorf("stat witches.env: %w", err)
	}
	projectPath, err := getwdWitchesInit()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName := filepath.Base(projectPath)
	if err := initWitchesInit(moduleName, DBdriver); err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}
