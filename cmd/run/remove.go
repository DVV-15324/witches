package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cmd_utils "github.com/DVV-15324/witches/cmd/utils"
	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

var statWitchesRemove = os.Stat
var getwdWitchesRemove = os.Getwd
var linkWitchesRemove = templates.RemoveModule
var stdinWitchesRemove io.Reader = os.Stdin
var skipConfirmWitchesRemove = false
var getModuleNameRemove = cmd_utils.GetModuleGoName

func WitchesRemove(moduleName string) error {
	_, err := statWitchesRemove("go.mod")
	if os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}
	if err != nil {
		return fmt.Errorf("stat go.mod: %w", err)
	}

	modulePath := filepath.Join("internal", moduleName)
	if _, err := statWitchesRemove(modulePath); os.IsNotExist(err) {
		return fmt.Errorf("module '%s' not found", moduleName)
	}

	// Confirmation prompt (skippable in tests)
	if !skipConfirmWitchesRemove && !confirmRemove(moduleName) {
		return fmt.Errorf("remove cancelled")
	}

	projectPath, err := getwdWitchesRemove()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	projectName, err := getModuleNameRemove()

	if err != nil {
		return fmt.Errorf("get module name: %w", err)
	}

	if err := linkWitchesRemove(projectPath, projectName, moduleName); err != nil {
		return fmt.Errorf("remove module: %w", err)
	}
	return nil
}

func confirmRemove(moduleName string) bool {
	fmt.Printf("Are you sure you want to remove module '%s'? This cannot be undone. [y/N]: ", moduleName)

	reader := bufio.NewReader(stdinWitchesRemove)
	response, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
