package easyjson

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratorEasyJsonRequest(t *testing.T) {
	rootDir := findProjectRoot(t)

	inputDir := filepath.Join(rootDir, "test", "easyjson", "dto", "request")

	//  Output vẫn là input, nhưng xóa file gen cũ trước
	outputDir := inputDir

	// Xóa file gen cũ nếu có
	oldGenFile := filepath.Join(outputDir, "request_easyjson.go")
	os.Remove(oldGenFile)

	// Chạy generator
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)

	// Kiểm tra
	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}

func TestGeneratorEasyJsonResponse(t *testing.T) {
	rootDir := findProjectRoot(t)

	inputDir := filepath.Join(rootDir, "test", "easyjson", "dto", "response")

	//  Output vẫn là input, nhưng xóa file gen cũ trước
	outputDir := inputDir

	// Xóa file gen cũ nếu có
	oldGenFile := filepath.Join(outputDir, "response_easyjson.go")
	os.Remove(oldGenFile)

	// Chạy generator
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)

	// Kiểm tra
	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}
func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Cannot find project root")
		}
		dir = parent
	}
}
