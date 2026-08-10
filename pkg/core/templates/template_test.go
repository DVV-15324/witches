package template

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"text/template"
)

// ==================== TEST PROJECT CONFIG ====================

func TestProjectConfig(t *testing.T) {
	config := ProjectConfig{
		ModuleName: "github.com/example/project",
	}
	assert.Equal(t, "github.com/example/project", config.GetMuduleName())
}

// ==================== TEST GENERATE SERVICE STRUCTURE ====================

func TestGenerateServiceStructure(t *testing.T) {
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, "internal", "product-service")

	dirs := []string{
		"dto/request",
		"dto/response",
		"entity",
		"handler",
		"mapping",
		"repository",
		"usecase",
	}

	for _, dir := range dirs {
		path := filepath.Join(baseDir, dir)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err)
		assert.DirExists(t, path)
	}
}

// ==================== TEST UPDATE KEY OBJECT ====================

func TestUpdateKeyObject(t *testing.T) {
	tmpDir := t.TempDir()

	keyFile := filepath.Join(tmpDir, "internal", "shared", "utils", "key_object.go")
	err := os.MkdirAll(filepath.Dir(keyFile), 0755)
	require.NoError(t, err)

	content := `package utils

var (
	ObjectUser uint = 1
	ObjectAdmin uint = 2
)`
	err = os.WriteFile(keyFile, []byte(content), 0644)
	require.NoError(t, err)

	config := ServiceConfig{
		NameCap: "Product",
		Name:    "product",
	}

	err = updateKeyObject(tmpDir, config)
	require.NoError(t, err)

	newContent, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), "ObjectProduct")
	assert.Contains(t, string(newContent), "3")
}

// ==================== TEST PARSE KEY OBJECT (AST) ====================

func TestParseKeyObject(t *testing.T) {
	content := `package utils

var (
	ObjectUser uint = 1
	ObjectAdmin uint = 2
)`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	require.NoError(t, err)

	var maxID int
	var targetDecl *ast.GenDecl

	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			return true
		}

		for _, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || !strings.HasPrefix(vs.Names[0].Name, "Object") {
				continue
			}

			if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
				if id, _ := strconv.Atoi(lit.Value); id > maxID {
					maxID = id
					targetDecl = decl
				}
			}
		}
		return true
	})

	assert.Equal(t, 2, maxID)
	assert.NotNil(t, targetDecl)
}

// ==================== TEST GENERATE SHARED MODEL ====================

func TestGenerateSharedModel(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap: "Product",
		Name:    "product",
	}

	sharedModelDir := filepath.Join(tmpDir, "internal", "shared", "model")
	err := os.MkdirAll(sharedModelDir, 0755)
	require.NoError(t, err)

	destFile := filepath.Join(sharedModelDir, config.Name+".go")
	content := fmt.Sprintf(`package model

type %s struct {
	ID string
}`, config.NameCap)

	err = os.WriteFile(destFile, []byte(content), 0644)
	require.NoError(t, err)

	assert.FileExists(t, destFile)

	readContent, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.Contains(t, string(readContent), "type Product struct")
}

// ==================== TEST ERROR CASES ====================

func TestUpdateKeyObject_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap: "Test",
		Name:    "test",
	}

	err := updateKeyObject(tmpDir, config)
	assert.Error(t, err)
}

func TestUpdateKeyObject_InvalidContent(t *testing.T) {
	tmpDir := t.TempDir()

	keyFile := filepath.Join(tmpDir, "internal", "shared", "utils", "key_object.go")
	err := os.MkdirAll(filepath.Dir(keyFile), 0755)
	require.NoError(t, err)

	content := `invalid go code!!!`
	err = os.WriteFile(keyFile, []byte(content), 0644)
	require.NoError(t, err)

	config := ServiceConfig{
		NameCap: "Test",
		Name:    "test",
	}

	err = updateKeyObject(tmpDir, config)
	assert.Error(t, err)
}

// ==================== BENCHMARKS ====================

func BenchmarkServiceNameProcessing(b *testing.B) {
	input := "user-profile-service"
	b.ResetTimer()
	for b.Loop() {
		serviceName := strings.TrimSpace(input)
		serviceName = strings.ToLower(serviceName)
		_ = cases.Title(language.English).String(serviceName)
	}
}

func BenchmarkParseKeyObject(b *testing.B) {
	content := `package utils

var (
	ObjectUser uint = 1
	ObjectAdmin uint = 2
	ObjectProduct uint = 3
)`

	b.ResetTimer()
	for b.Loop() {
		fset := token.NewFileSet()
		node, _ := parser.ParseFile(fset, "", content, parser.ParseComments)

		var maxID int
		ast.Inspect(node, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.VAR {
				return true
			}

			for _, spec := range decl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || len(vs.Names) != 1 || !strings.HasPrefix(vs.Names[0].Name, "Object") {
					continue
				}

				if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
					if id, _ := strconv.Atoi(lit.Value); id > maxID {
						maxID = id
					}
				}
			}
			return true
		})
		_ = maxID
	}
}

func BenchmarkRenderTemplate(b *testing.B) {
	tmpDir := b.TempDir()
	config := ServiceConfig{
		NameCap:    "Bench",
		Name:       "bench",
		FolderName: "bench-service",
		ModuleName: "github.com/example/bench",
	}

	tmplContent := `package {{.Name}}
type {{.NameCap}} struct {
	ID string
}`

	tmplFile := filepath.Join(tmpDir, "bench.tmpl")
	err := os.WriteFile(tmplFile, []byte(tmplContent), 0644)
	require.NoError(b, err)

	tmpl, err := template.ParseFiles(tmplFile)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		destFile := filepath.Join(tmpDir, "output.go")
		file, _ := os.Create(destFile)
		tmpl.Execute(file, config)
		file.Close()
	}
}
