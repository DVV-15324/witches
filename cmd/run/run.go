package cmd

import (
	"fmt"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DVV-15324/witches/pkg/core/easyjson"
)

var getwdProjectRoot = os.Getwd
var generateEasyJSON = easyjson.GeneratorEasyJson
var generateAllDTOs = generateEasyJSONForAllDTOs

func WitchesRun() error {
	if err := generateAllDTOs(); err != nil {
		return fmt.Errorf("generate easyjson: %w", err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run project: %w", err)
	}
	return nil
}

func generateEasyJSONForAllDTOs() error {
	rootDir := findProjectRoot()
	if rootDir == "" {
		return fmt.Errorf("cannot find project root")
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
		return fmt.Errorf("walk internal: %w", err)
	}

	if len(dtoDirs) == 0 {
		return nil
	}

	for _, dtoDir := range dtoDirs {
		requestDir := filepath.Join(dtoDir, "request")
		responseDir := filepath.Join(dtoDir, "response")

		if _, err := os.Stat(requestDir); err == nil {
			if err := generateEasyJSONForDir(requestDir, "request"); err != nil {
				return err
			}
		}
		if _, err := os.Stat(responseDir); err == nil {
			if err := generateEasyJSONForDir(responseDir, "response"); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateEasyJSONForDir(dir, name string) error {
	removeEasyJSONFiles(dir)
	fset := token.NewFileSet()
	err := generateEasyJSON(fset, dir, dir)
	if err != nil {
		return fmt.Errorf("%s error: %w", name, err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if len(files) > 0 {
		fmt.Printf("Generated %s easyjson: %d file(s)\n", name, len(files))
	}
	return nil
}

func findProjectRoot() string {
	dir, err := getwdProjectRoot()
	if err != nil {
		return ""
	}

	// Chỉ kiểm tra thư mục hiện tại, không tìm lên cha
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return dir
	}
	return ""
}

func removeEasyJSONFiles(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if err != nil {
		return
	}
	for _, f := range files {
		_ = os.Remove(f)
	}
}
