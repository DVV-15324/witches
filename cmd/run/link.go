package cmd

import (
	"fmt"
	"os"

	cmd_utils "github.com/DVV-15324/witches/cmd/utils"
	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

var statWitchesLink = os.Stat
var getwdWitchesLink = os.Getwd
var linkWitchesLink = templates.AddGoDomainFromLink
var getModuleNameLink = cmd_utils.GetModuleName

func WitchesLink(domainName, repoURL string) error {
	_, err := statWitchesLink("go.mod")
	if os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found")
	}
	if err != nil {
		return fmt.Errorf("stat go.mod: %w", err)
	}

	projectPath, err := getwdWitchesLink()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	moduleName, err := getModuleNameLink()
	if err != nil {
		return fmt.Errorf("get module name: %w", err)
	}

	if err := linkWitchesLink(projectPath, moduleName, domainName, repoURL); err != nil {
		return fmt.Errorf("link domain: %w", err)
	}
	return nil
}
