package cmd

import (
	"fmt"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DVV-15324/witches/pkg/core/easyjson"
)

func WitchesRun() {
	generateEasyJSONForAllDTOs()
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func generateEasyJSONForAllDTOs() {
	rootDir := findProjectRoot()
	if rootDir == "" {
		log.Fatalf("Error: Cannot find project %s", rootDir)
		return
	}
	var dtoDirs []string
	err := filepath.Walk(filepath.Join(rootDir, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && filepath.Base(path) == "dto" {
			dtoDirs = append(dtoDirs, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error: walking %v", rootDir)
		return
	}
	if len(dtoDirs) == 0 {
		return
	}
	for _, dtoDir := range dtoDirs {
		requestDir := filepath.Join(dtoDir, "request")
		responseDir := filepath.Join(dtoDir, "response")

		if _, err := os.Stat(requestDir); err == nil {
			generateEasyJSONForDir(requestDir, "request")
		}
		if _, err := os.Stat(responseDir); err == nil {
			generateEasyJSONForDir(responseDir, "response")
		}
	}
}

func generateEasyJSONForDir(dir, name string) {
	removeEasyJSONFiles(dir)
	fset := token.NewFileSet()
	err := easyjson.GeneratorEasyJson(fset, dir, dir)
	if err != nil {
		log.Fatalf("Error: %s error: %v\n", name, err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if len(files) > 0 {
		fmt.Printf("Generated %s easyjson: %d file(s)\n", name, len(files))
	}
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: Cannot get current directory: %v\n", err)
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Println("Error: Cannot find project root (go.mod not found)")
			return ""
		}
		dir = parent
	}
}
func removeEasyJSONFiles(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if err != nil {
		return
	}
	for _, f := range files {
		os.Remove(f)
	}
}
