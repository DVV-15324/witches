package easyjson

import (
	"os"

	"path/filepath"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/token"
)

// ==================== Test GeneratorEasyJson ====================

func TestGeneratorEasyJsonRequest(t *testing.T) {
	rootDir := findProjectRoot(t)
	inputDir := filepath.Join(rootDir, "pkg", "core", "easyjson", "test", "request")

	outputDir := inputDir
	oldGenFile := filepath.Join(outputDir, "request_easyjson.go")
	os.Remove(oldGenFile)

	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)
	assert.NoError(t, err)

	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}

func TestGeneratorEasyJsonResponse(t *testing.T) {
	rootDir := findProjectRoot(t)
	inputDir := filepath.Join(rootDir, "pkg", "core", "easyjson", "test", "response")

	outputDir := inputDir
	oldGenFile := filepath.Join(outputDir, "response_easyjson.go")
	os.Remove(oldGenFile)

	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, inputDir, outputDir)
	assert.NoError(t, err)

	if _, err = os.Stat(oldGenFile); err == nil {
		t.Logf("Generated file found: %s", oldGenFile)
	} else {
		t.Log("No generated file")
	}
}

// -------------------- Test GeneratorEasyJson with file input --------------------
func TestGeneratorEasyJson_InputIsFile(t *testing.T) {
	// Skip: easyjson requires full GOPATH/module context that temp dir can't satisfy
	t.Skip("Skipping: easyjson requires proper go module context with dependencies resolved")
}

func TestGeneratorEasyJson_InputDirNoMarker(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test
type T struct{}`
	err := os.WriteFile(goFile, []byte(content), 0644)
	require.NoError(t, err)

	fset := token.NewFileSet()
	err = GeneratorEasyJson(fset, tmpDir, "")
	assert.NoError(t, err)
}

func TestGeneratorEasyJson_InputNotExist(t *testing.T) {
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, "/non/existent/path", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input path does not exist")
}

func TestGeneratorEasyJson_InputEmpty(t *testing.T) {
	fset := token.NewFileSet()
	err := GeneratorEasyJson(fset, "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "input path is empty")
}

// ==================== Test generateEasyJSON ====================

func TestGenerateEasyJSON_FileNotExist(t *testing.T) {
	err := generateEasyJSON("non_existent_file.go")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}

func TestGenerateEasyJSON_FileExists(t *testing.T) {
	// Skip: easyjson requires full GOPATH/module context with dependencies resolved
	t.Skip("Skipping: easyjson requires proper go module context with dependencies resolved")
}

func TestGenerateEasyJSON_NoEasyJSON(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", "")

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(tmpFile, []byte("package test\ntype T struct{}"), 0644)
	require.NoError(t, err)

	err = generateEasyJSON(tmpFile)
	assert.Error(t, err)
	// Lỗi có thể là "executable file not found" hoặc "no such file or directory"
	assert.Contains(t, err.Error(), "executable")
}

// ==================== Helper ====================

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
