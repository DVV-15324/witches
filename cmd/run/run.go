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

// WitchesRun - Runs the application with auto-generate easyjson
func WitchesRun() {
	// Gen easyjson cho tất cả DTO trong project
	generateEasyJSONForAllDTOs()

	// Start executing
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// generateEasyJSONForAllDTOs - Gen easyjson cho tất cả DTO trong internal/*/dto
func generateEasyJSONForAllDTOs() {
	rootDir := findProjectRoot()
	if rootDir == "" {
		fmt.Println("Cannot find project root")
		return
	}

	// Tìm tất cả thư mục dto
	var dtoDirs []string
	err := filepath.Walk(filepath.Join(rootDir, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && filepath.Base(path) == "dto" {
			dtoDirs = append(dtoDirs, path)
			//fmt.Printf("  Found DTO folder: %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking: %v\n", err)
		return
	}

	if len(dtoDirs) == 0 {
		//fmt.Println("No DTO folders found")
		return
	}

	// Gen cho từng DTO folder
	for _, dtoDir := range dtoDirs {
		// Gen cho request và response
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

// generateEasyJSONForDir - Gen easyjson cho 1 thư mục
func generateEasyJSONForDir(dir, name string) {
	// Xóa file gen cũ
	removeEasyJSONFiles(dir)

	fset := token.NewFileSet()
	err := easyjson.GeneratorEasyJson(fset, dir, dir)
	if err != nil {
		fmt.Printf("%s error: %v\n", name, err)
	}

	// Kiểm tra file gen
	files, _ := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if len(files) > 0 {
		fmt.Printf("Generated %s easyjson: %d file(s)\n", name, len(files))
	} else {
		fmt.Printf("No %s_easyjson.go generated\n", name)
	}
}

// findProjectRoot - Tìm thư mục chứa go.mod
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Cannot get current directory: %v\n", err)
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Println("Cannot find project root (go.mod not found)")
			return ""
		}
		dir = parent
	}
}

// removeEasyJSONFiles - Xóa tất cả file *_easyjson.go trong thư mục
func removeEasyJSONFiles(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*_easyjson.go"))
	if err != nil {
		return
	}
	for _, f := range files {
		os.Remove(f)
	}
}
