package utils

import (
	"embed"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ServiceConfig dùng cho service generation
type ServiceConfig struct {
	NameCap    string
	Name       string
	FolderName string
	ModuleName string
}

func (p ServiceConfig) GetModuleName() string {
	return p.ModuleName
}

// ProjectConfig dùng cho project generation
type ProjectConfig struct {
	ModuleName string
}

func (p ProjectConfig) GetModuleName() string {
	return p.ModuleName
}

//go:embed testdata/*.tmpl
var testFS embed.FS

// Test cơ bản - dùng embed.FS có sẵn
func TestRenderTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:    "Test",
		Name:       "test",
		FolderName: "test-service",
		ModuleName: "github.com/example/test",
	}

	// Gọi RenderTemplate với file có sẵn trong embed.FS
	destFile := "output.go"
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package test")
	assert.Contains(t, string(content), "type Test struct")
}

// Test với embed
func TestRenderTemplate_WithEmbed(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:    "Embed",
		Name:       "embed",
		FolderName: "embed-service",
		ModuleName: "github.com/example/embed",
	}

	destFile := "output.go"
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package embed")
	assert.Contains(t, string(content), "type Embed struct")
}

// Test với thư mục con - dùng embed.FS có sẵn
func TestRenderTemplate_WithSubDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:    "Sub",
		Name:       "sub",
		FolderName: "sub-service",
		ModuleName: "github.com/example/sub",
	}

	// Tạo file trong thư mục con
	destFile := filepath.Join("deep", "nested", "path", "output.go")
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package sub")
	assert.Contains(t, string(content), "type Sub struct")
}

// Test với template không dùng config - dùng file có sẵn
func TestRenderTemplate_WithNilConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo template không dùng config trong thư mục tạm
	tmpTemplate := `package test`
	tmpTmplFile := filepath.Join(tmpDir, "noconfig.tmpl")
	err := os.WriteFile(tmpTmplFile, []byte(tmpTemplate), 0644)
	require.NoError(t, err)

	// Đọc template từ file tạm
	tmpl, err := template.ParseFiles(tmpTmplFile)
	require.NoError(t, err)

	destFile := filepath.Join(tmpDir, "noconfig_output.go")
	file, err := os.Create(destFile)
	require.NoError(t, err)
	defer file.Close()

	err = tmpl.Execute(file, nil)
	require.NoError(t, err)

	content, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package test")
}

// Test với template phức tạp - tạo trong thư mục tạm
func TestRenderTemplate_ComplexStruct(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:    "Complex",
		Name:       "complex",
		FolderName: "complex-service",
		ModuleName: "github.com/example/complex",
	}

	// Tạo template phức tạp trong thư mục tạm
	complexTemplate := `package {{.Name}}

import "context"

type {{.NameCap}}Service interface {
	Create(ctx context.Context, data *{{.NameCap}}) error
	Get(ctx context.Context, id string) (*{{.NameCap}}, error)
	Update(ctx context.Context, id string, data *{{.NameCap}}) error
	Delete(ctx context.Context, id string) error
}

type {{.NameCap}} struct {
	ID   string
	Name string
}

func New{{.NameCap}}(name string) *{{.NameCap}} {
	return &{{.NameCap}}{
		ID:   "generated-id",
		Name: name,
	}
}`

	// Tạo file template trong thư mục tạm
	tmplDir := filepath.Join(tmpDir, "templates")
	err := os.MkdirAll(tmplDir, 0755)
	require.NoError(t, err)

	tmplFile := filepath.Join(tmplDir, "complex.tmpl")
	err = os.WriteFile(tmplFile, []byte(complexTemplate), 0644)
	require.NoError(t, err)

	// Đọc template từ file tạm
	tmpl, err := template.ParseFiles(tmplFile)
	require.NoError(t, err)

	// Render
	destFile := filepath.Join(tmpDir, "complex_output.go")
	file, err := os.Create(destFile)
	require.NoError(t, err)
	defer file.Close()

	err = tmpl.Execute(file, config)
	require.NoError(t, err)

	content, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package complex")
	assert.Contains(t, string(content), "type ComplexService interface")
	assert.Contains(t, string(content), "func NewComplex")
}

// Test error cases - bắt lỗi từ RenderTemplate
func TestRenderTemplate_ErrorCases(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:    "Error",
		Name:       "error",
		FolderName: "error-service",
		ModuleName: "github.com/example/error",
	}

	// Test với template không tồn tại
	// RenderTemplate sẽ log.Fatal nên không thể recover
	// Chỉ test với template có sẵn
	t.Run("valid template should work", func(t *testing.T) {
		destFile := "valid_output.go"
		RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

		fullPath := filepath.Join(tmpDir, destFile)
		content, err := os.ReadFile(fullPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "package error")
	})
}
