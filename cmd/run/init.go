package cmd

import (
	"log"
	"os"
	"path/filepath"

	templates "github.com/DVV-15324/witches/pkg/core/templates"
)

func WitchesInitCaptain(DBdriver string) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	moduleName := filepath.Base(projectPath)
	templates.CreateCaptainGoArc(moduleName, DBdriver)
}

func WitchesInitMember(DBdriver string) {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	moduleName := filepath.Base(projectPath)
	templates.CreateMemberGoArc(moduleName, DBdriver)
}
