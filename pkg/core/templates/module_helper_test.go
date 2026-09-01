package template

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed module/shared/domain/domain.go.tmpl
var testFS embed.FS

// ==================== TEST generateSharedDomain ====================

func TestGenerateSharedDomain_Success(t *testing.T) {
	tmpDir := t.TempDir()

	config := ModuleConfig{
		Name:    "book",
		NameCap: "Book",
	}

	err := generateSharedDomain(tmpDir, config, testFS)
	assert.NoError(t, err)

	destFile := filepath.Join(tmpDir, "internal", "shared", "domain", "book.go")
	content, err := os.ReadFile(destFile)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestGenerateSharedDomain_TemplateNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	config := ModuleConfig{
		Name:    "book",
		NameCap: "Book",
	}

	// Dùng embed.FS không có template
	var emptyFS embed.FS
	err := generateSharedDomain(tmpDir, config, emptyFS)
	assert.Error(t, err)
	//  Sửa: error message thực tế
	assert.Contains(t, err.Error(), "template module/shared/domain/domain.go.tmpl not found")
}

// ==================== TEST updateKeyObject ====================

func TestUpdateKeyObject_Success(t *testing.T) {
	tmpDir := t.TempDir()
	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	err := os.MkdirAll(utilsDir, 0755)
	require.NoError(t, err)

	keyFile := filepath.Join(utilsDir, "key_object.go")
	content := `package utils

type KeyObject int64

const (
	ObjectUser KeyObject = 1
	ObjectPost KeyObject = 2
)
`
	err = os.WriteFile(keyFile, []byte(content), 0644)
	require.NoError(t, err)

	config := ModuleConfig{
		Name:    "book",
		NameCap: "Book",
	}

	err = updateKeyObject(tmpDir, config)
	assert.NoError(t, err)

	newContent, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), "ObjectBook KeyObject = 3")
}

func TestUpdateKeyObject_NoObjectConstant(t *testing.T) {
	tmpDir := t.TempDir()
	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	err := os.MkdirAll(utilsDir, 0755)
	require.NoError(t, err)

	keyFile := filepath.Join(utilsDir, "key_object.go")
	content := `package utils

type KeyObject int64
`
	err = os.WriteFile(keyFile, []byte(content), 0644)
	require.NoError(t, err)

	config := ModuleConfig{
		Name:    "book",
		NameCap: "Book",
	}

	err = updateKeyObject(tmpDir, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Object constant found")
}

func TestUpdateKeyObject_FileNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	config := ModuleConfig{
		Name:    "book",
		NameCap: "Book",
	}

	err := updateKeyObject(tmpDir, config)
	assert.Error(t, err)
}

// ==================== TEST AddModuleField ====================

func TestAddModuleField_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")

	content := `package routers

type Modules struct {
	User *user.UserModule
}

func (m *Modules) InitModules() *Modules {
	return &Modules{
		User: userModule,
	}
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddModuleField(filePath, "book", "Book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/example/project/internal/book"`)
	assert.Contains(t, string(newContent), `Book *book.BookModule`)
}

func TestAddModuleField_FileNotExist(t *testing.T) {
	err := AddModuleField("/nonexistent/file.go", "book", "Book", "github.com/example/project")
	assert.Error(t, err)
}

func TestAddModuleField_StructNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")

	content := `package routers

type OtherStruct struct {}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddModuleField(filePath, "book", "Book", "github.com/example/project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "struct Modules not found")
}

func TestAddModuleField_FieldExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")

	content := `package routers

type Modules struct {
	Book *book.BookModule
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddModuleField(filePath, "book", "Book", "github.com/example/project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "field Book already exists")
}

// ==================== TEST AddModuleInit ====================

func TestAddModuleInit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")

	content := `package routers

type Modules struct {
	User *user.UserModule
}

func (m *Modules) InitModules() *Modules {
	return &Modules{
		User: userModule,
	}
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddModuleInit(filePath, "book", "Book")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `Book: bookModule`)
	assert.Contains(t, string(newContent), `bookModule := book.NewBookModule(core)`)
}

func TestAddModuleInit_FunctionNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")

	content := `package routers

type Modules struct{}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddModuleInit(filePath, "book", "Book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "function InitModules not found")
}

// ==================== TEST AddRouteRegistration ====================

func TestAddRouteRegistration_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "routers.go")

	content := `package routers

func initModule(gen *SwaggerGenerator, modules *Modules) error {
	return nil
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddRouteRegistration(filePath, "book", "Book")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `gen.AddTag("book", "Book endpoints")`)
	assert.Contains(t, string(newContent), `modules.Book.RegisterPublicRoutes(gen)`)
	assert.Contains(t, string(newContent), `modules.Book.RegisterProtectedRoutes(gen)`)
}

func TestAddRouteRegistration_FunctionNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "routers.go")

	content := `package routers
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddRouteRegistration(filePath, "book", "Book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "function initModule not found")
}

func TestAddRouteRegistration_DomainExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "routers.go")

	content := `package routers

func initModule(gen *SwaggerGenerator, modules *Modules) error {
	gen.AddTag("book", "Book endpoints")
	modules.Book.RegisterPublicRoutes(gen)
	modules.Book.RegisterProtectedRoutes(gen)
	return nil
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddRouteRegistration(filePath, "book", "Book")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestAddRouteRegistration_EmptyBody(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "routers.go")

	content := `package routers

func initModule(gen *SwaggerGenerator, modules *Modules) error {
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddRouteRegistration(filePath, "book", "Book")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `gen.AddTag("book", "Book endpoints")`)
}

// ==================== BENCHMARK ====================

func BenchmarkAddModuleField(b *testing.B) {
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "modules.go")
	content := `package routers

type Modules struct {
	User *user.UserModule
}

func (m *Modules) InitModules() *Modules {
	return &Modules{
		User: userModule,
	}
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		_ = AddModuleField(filePath, "book", "Book", "github.com/example/project")
	}
}
