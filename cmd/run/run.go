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
	// Gen easyjson cho request và response
	GeneratorEasyJsonRequest()
	GeneratorEasyJsonResponse()

	// Start executing
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// GeneratorEasyJsonRequest - Gen easyjson cho DTO request
func GeneratorEasyJsonRequest() {
	rootDir := findProjectRoot()
	if rootDir == "" {
		fmt.Println("Cannot find project root")
		return
	}

	inputDir := filepath.Join(rootDir, "test", "easyjson", "dto", "request")
	outputDir := inputDir

	// Xóa tất cả file *_easyjson.go cũ
	removeEasyJSONFiles(outputDir)

	fset := token.NewFileSet()
	err := easyjson.GeneratorEasyJson(fset, inputDir, outputDir)
	if err != nil {
		log.Fatalln("No file go request generated")
	}
	// Kiểm tra file gen
	files, _ := filepath.Glob(filepath.Join(outputDir, "*_easyjson.go")) //Glob => Chỉ tìm 1 cấp (*.go) trong khi filepath.Walk Tìm tất cả cấp (**/*.go) => Glob trường hợp này tìm trong 1 folder nên Glop tạm
	if len(files) > 0 {
		fmt.Printf("Generated request easyjson: %d file(s)\n", len(files))
	} else {
		fmt.Println("No file go request generated")
	}
}

// GeneratorEasyJsonResponse - Gen easyjson cho DTO response
func GeneratorEasyJsonResponse() {
	rootDir := findProjectRoot()
	if rootDir == "" {
		fmt.Println("Cannot find project root")
		return
	}

	inputDir := filepath.Join(rootDir, "test", "easyjson", "dto", "response")
	outputDir := inputDir

	// Xóa tất cả file *_easyjson.go cũ
	removeEasyJSONFiles(outputDir)

	fset := token.NewFileSet()
	err := easyjson.GeneratorEasyJson(fset, inputDir, outputDir)

	if err != nil {
		log.Fatalln("No file go response generated")
	}

	// Kiểm tra file gen
	files, _ := filepath.Glob(filepath.Join(outputDir, "*_easyjson.go")) //Glob => Chỉ tìm 1 cấp (*.go) trong khi filepath.Walk Tìm tất cả cấp (**/*.go) => Glob trường hợp này tìm trong 1 folder nên Glop tạm
	if len(files) > 0 {
		fmt.Printf("Generated response easyjson: %d file(s)\n", len(files))
	} else {
		fmt.Println("No file go response generated")
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
	files, err := filepath.Glob(filepath.Join(dir, "*_easyjson.go")) //Glob => Chỉ tìm 1 cấp (*.go) trong khi filepath.Walk Tìm tất cả cấp (**/*.go) => Glob trường hợp này tìm trong 1 folder nên Glop tạm
	if err != nil {
		return
	}
	for _, f := range files {
		if err := os.Remove(f); err == nil {
			fmt.Printf("Removed old: %s\n", filepath.Base(f))
		}
	}
}
