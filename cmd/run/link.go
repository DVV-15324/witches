package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

var statWitchesLink = os.Stat
var getwdWitchesLink = os.Getwd
var linkWitchesLink = templates.AddGoDomainFromLink

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

	moduleName := filepath.Base(projectPath)
	if err := linkWitchesLink(projectPath, moduleName, domainName, repoURL); err != nil {
		return fmt.Errorf("link domain: %w", err)
	}
	return nil
}
