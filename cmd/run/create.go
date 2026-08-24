package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	content "github.com/DVV-15324/witches/cmd/utils"
)

func WitchesCreate(project string) {
	projectPath := filepath.Join(".", project)
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		log.Fatalf("Error: Project '%s' already exists!", project)
		return
	}
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	envPath := filepath.Join(projectPath, "witches.env")
	if _, err := os.Stat(envPath); err == nil {
		log.Fatalf("Error: File 'witches.env' already exists in '%s'", projectPath)
		return
	}
	file, err := os.OpenFile(envPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer func() { _ = file.Close() }()
	contentData := content.CreateContentRefresh()
	_, err = file.WriteString(contentData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Project '%s' created successfully!\n", project)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", project)
	fmt.Printf("  Edit witches.env\n")
	fmt.Printf("  witches database generate\n")
}
