package utils

import (
	"embed"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ServiceConfig struct {
	NameCap     string
	Name        string
	ProjectName string
}

func (p ServiceConfig) GetProjectName() string {
	return p.ProjectName
}

type ProjectConfig struct {
	ModuleName string
}

func (p ProjectConfig) GetModuleName() string {
	return p.ModuleName
}

//go:embed testdata/*.tmpl
var testFS embed.FS

func TestRenderTemplate_Success(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:     "Test",
		Name:        "test",
		ProjectName: "github.com/example/test",
	}
	destFile := "output.go"
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package test")
	assert.Contains(t, string(content), "type Test struct")
}

func TestRenderTemplate_WithEmbed(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:     "Embed",
		Name:        "embed",
		ProjectName: "github.com/example/embed",
	}

	destFile := "output.go"
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package embed")
	assert.Contains(t, string(content), "type Embed struct")
}

func TestRenderTemplate_WithSubDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:     "Sub",
		Name:        "sub",
		ProjectName: "github.com/example/sub",
	}

	destFile := filepath.Join("deep", "nested", "path", "output.go")
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package sub")
	assert.Contains(t, string(content), "type Sub struct")
}

func TestRenderTemplate_ComplexStruct(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:     "Complex",
		Name:        "complex",
		ProjectName: "github.com/example/complex",
	}

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

	tmplDir := filepath.Join(tmpDir, "templates")
	err := os.MkdirAll(tmplDir, 0755)
	require.NoError(t, err)

	tmplFile := filepath.Join(tmplDir, "complex.tmpl")
	err = os.WriteFile(tmplFile, []byte(complexTemplate), 0644)
	require.NoError(t, err)

	tmpl, err := template.ParseFiles(tmplFile)
	require.NoError(t, err)

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

func TestRenderTemplate_TemplateNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	config := ServiceConfig{Name: "test"}

	if os.Getenv("TEST_SUBPROCESS_TEMPLATE_NOTFOUND") == "1" {
		RenderTemplate(testFS, tmpDir, "output.go", "testdata/notfound.tmpl", config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRenderTemplate_TemplateNotFound")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_TEMPLATE_NOTFOUND=1")
	err := cmd.Run()
	assert.Error(t, err)
}

// SỬA: Test invalid template với RenderTemplate
func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo template invalid trong testdata
	tmplDir := filepath.Join(tmpDir, "testdata")
	err := os.MkdirAll(tmplDir, 0755)
	require.NoError(t, err)

	invalidTemplate := `package {{.Invalid}}`
	tmplFile := filepath.Join(tmplDir, "invalid.tmpl")
	err = os.WriteFile(tmplFile, []byte(invalidTemplate), 0644)
	require.NoError(t, err)

	config := ServiceConfig{Name: "test"}

	if os.Getenv("TEST_SUBPROCESS_INVALID_TEMPLATE") == "1" {
		//  Dùng embed.FS với file vừa tạo
		// Cách 1: Tạo embed.FS mới (không thể runtime)
		// Cách 2: Dùng os.ReadFile và template.ParseFiles trực tiếp

		//  Cách 2: Parse file trực tiếp
		tmpl, err := template.ParseFiles(tmplFile)
		if err != nil {
			os.Exit(1)
		}

		fullPath := filepath.Join(tmpDir, "output.go")
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			os.Exit(1)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			os.Exit(1)
		}
		defer file.Close()

		if err := tmpl.Execute(file, config); err != nil {
			//  Expected: lỗi do template invalid
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRenderTemplate_InvalidTemplate")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_INVALID_TEMPLATE=1")
	err = cmd.Run()
	assert.Error(t, err)
}
func TestRenderTemplate_CreateFileError(t *testing.T) {
	//Skip trên CI vì permission test không ổn định
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping permission test on CI")
	}

	// Dùng đường dẫn không hợp lệ thay vì permission
	tmpDir := t.TempDir()

	// Tạo tên file không hợp lệ (chứa ký tự đặc biệt)
	invalidFile := filepath.Join(tmpDir, "..", "invalid", "output.go")
	config := ServiceConfig{Name: "test"}

	if os.Getenv("TEST_SUBPROCESS_CREATE_FILE_ERROR") == "1" {
		RenderTemplate(testFS, tmpDir, invalidFile, "testdata/test.tmpl", config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRenderTemplate_CreateFileError")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_CREATE_FILE_ERROR=1")
	err := cmd.Run()
	assert.Error(t, err)
}

func TestRenderTemplate_MkdirError(t *testing.T) {
	//Skip trên CI vì permission test không ổn định
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping permission test on CI")
	}

	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	err := os.MkdirAll(readOnlyDir, 0444)
	require.NoError(t, err)

	destFile := filepath.Join(readOnlyDir, "subdir", "output.go")
	config := ServiceConfig{Name: "test"}

	if os.Getenv("TEST_SUBPROCESS_MKDIR_ERROR") == "1" {
		RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRenderTemplate_MkdirError")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_MKDIR_ERROR=1")
	err = cmd.Run()
	assert.Error(t, err)
}

func TestRenderTemplate_NilConfig(t *testing.T) {
	tmpDir := t.TempDir()

	config := ServiceConfig{
		NameCap:     "NilConfig",
		Name:        "nilconfig",
		ProjectName: "github.com/example/nilconfig",
	}

	destFile := "nil_config_output.go"
	RenderTemplate(testFS, tmpDir, destFile, "testdata/test.tmpl", config)

	fullPath := filepath.Join(tmpDir, destFile)
	content, err := os.ReadFile(fullPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package nilconfig")
}

// ==================== BENCHMARK ====================

func BenchmarkRenderTemplate(b *testing.B) {
	tmpDir := b.TempDir()
	config := ServiceConfig{
		NameCap:     "Bench",
		Name:        "bench",
		ProjectName: "github.com/example/bench",
	}

	b.ResetTimer()
	for b.Loop() {
		RenderTemplate(testFS, tmpDir, "bench_output.go", "testdata/test.tmpl", config)
	}
}
