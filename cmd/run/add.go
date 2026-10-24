package cmd

import (
	"fmt"
	cmd_utils "github.com/DVV-15324/witches/cmd/utils"
	templates "github.com/DVV-15324/witches/pkg/core/templates"
	"os"
	"path/filepath"
)

var getwdWitchesAdd = os.Getwd
var addModule = templates.AddModule
var getModuleNameAdd = cmd_utils.GetModuleName

func WitchesAdd(moduleName string, DBdriver string) error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}

	domainPath := filepath.Join("internal", moduleName)
	if _, err := os.Stat(domainPath); err == nil {
		return fmt.Errorf("domain '%s' already exists", moduleName)
	}

	projectPath, err := getwdWitchesAdd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	projectName, err := getModuleNameAdd()
	if err != nil {
		return fmt.Errorf("get module name: %w", err)
	}

	if err := addModule(projectPath, projectName, moduleName, DBdriver); err != nil {
		return fmt.Errorf("add module: %w", err)
	}
	return nil
}
