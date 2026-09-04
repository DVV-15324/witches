package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	content "github.com/DVV-15324/witches/cmd/utils"
)

var mkdirAllWitchesCreate = os.MkdirAll
var openFileWitchesCreate = os.OpenFile

func WitchesCreate(project string) error {
	projectPath := filepath.Join(".", project)
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("project '%s' already exists", project)
	}

	if err := mkdirAllWitchesCreate(projectPath, os.ModePerm); err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}

	envPath := filepath.Join(projectPath, "witches.env")
	_, err := os.Stat(envPath)
	if err == nil {
		return fmt.Errorf("file 'witches.env' already exists in '%s'", projectPath)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("stat witches.env: %w", err)
	}
	file, err := openFileWitchesCreate(envPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("create witches.env: %w", err)
	}
	defer file.Close()

	contentData := content.CreateContent()
	if _, err := file.WriteString(contentData); err != nil {
		return fmt.Errorf("write witches.env: %w", err)
	}

	fmt.Printf("Project '%s' created successfully!\n", project)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", project)
	fmt.Printf("  Edit witches.env\n")
	fmt.Printf("  witches database generate\n")
	return nil
}
