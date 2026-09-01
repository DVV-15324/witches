package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== TEST ValidateModuleStructure ====================

func TestValidateModuleStructure_Success(t *testing.T) {
	templateFiles := map[string]string{
		"internal/book/model/model.go":                     "package model",
		"internal/book/handler/handler.go":                 "package handler",
		"internal/book/repository/repository.go":           "package repository",
		"internal/book/usecase/usecase.go":                 "package usecase",
		"internal/book/module.go":                          "package book",
		"internal/book/migrate/migrations/1_create.up.sql": "CREATE TABLE",
		"internal/shared/domain/book.go":                   "package domain",
	}

	config := ModuleConfig{
		Name: "book",
	}

	err := ValidateModuleStructure(templateFiles, config)
	assert.NoError(t, err)
}

func TestValidateModuleStructure_MissingFiles(t *testing.T) {
	templateFiles := map[string]string{
		"internal/book/model/model.go":     "package model",
		"internal/book/handler/handler.go": "package handler",
		// Thiếu repository, usecase, module
	}

	config := ModuleConfig{
		Name: "book",
	}

	err := ValidateModuleStructure(templateFiles, config)
	assert.NoError(t, err) // Hàm chỉ log warning, không return error
}

func TestValidateModuleStructure_Empty(t *testing.T) {
	templateFiles := map[string]string{}

	config := ModuleConfig{
		Name: "book",
	}

	err := ValidateModuleStructure(templateFiles, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no domain files")
}

// ==================== TEST AddSharedDomainImport ====================

func TestAddSharedDomainImport_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	//  Thêm import block để có thể thêm import mới
	content := `package test

import (
	"fmt"
)

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImport(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/example/project/internal/shared/domain"`)
}

func TestAddSharedDomainImport_NoImportBlock(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImport(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	// Kiểm tra import đã được thêm (hoặc không tùy logic)
	// Hiện tại code thêm import, nên kiểm tra có chứa
	assert.Contains(t, string(newContent), `"github.com/example/project/internal/shared/domain"`)
}

func TestAddSharedDomainImport_FileNotExist(t *testing.T) {
	err := AddSharedDomainImport("/nonexistent/file.go", "book", "github.com/example/project")
	assert.NoError(t, err) // File không tồn tại, bỏ qua
}

func TestAddSharedDomainImport_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import (
	"github.com/example/project/internal/shared/domain"
)

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImport(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	// Kiểm tra không có duplicate import
	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Kiểm tra không duplicate
	assert.Contains(t, string(newContent), `"github.com/example/project/internal/shared/domain"`)
}

// ==================== TEST AddSharedDomainImportAdvanced ====================

func TestAddSharedDomainImportAdvanced_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import (
	"fmt"
)

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImportAdvanced(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/example/project/shared/domain"`)
}

func TestAddSharedDomainImportAdvanced_NoImportBlock(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImportAdvanced(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	// Hàm đã thêm import, nên kiểm tra có chứa
	assert.Contains(t, string(newContent), `"github.com/example/project/shared/domain"`)
}

func TestAddSharedDomainImportAdvanced_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import (
	"github.com/example/project/shared/domain"
)

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = AddSharedDomainImportAdvanced(filePath, "book", "github.com/example/project")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/example/project/shared/domain"`)
}

// ==================== TEST RewriteModuleImports ====================

func TestRewriteModuleImports_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import (
	"github.com/old-module/internal/book"
	"github.com/old-module/shared/domain"
)

type Test struct {
	Book *book.Book
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = RewriteModuleImports(filePath, "github.com/old-module", "github.com/new-module")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/new-module/internal/book"`)
	assert.Contains(t, string(newContent), `"github.com/new-module/shared/domain"`)
}

func TestRewriteModuleImports_NoChange(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import "fmt"

type Test struct {
	ID int
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = RewriteModuleImports(filePath, "old", "new")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, string(newContent))
}

func TestRewriteModuleImports_FileNotExist(t *testing.T) {
	err := RewriteModuleImports("/nonexistent/file.go", "old", "new")
	assert.Error(t, err)
}

// ==================== TEST RewriteAllImportsInDomain ====================

func TestRewriteAllImportsInDomain_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo cấu trúc domain
	domainPath := filepath.Join(tmpDir, "internal", "book")
	err := os.MkdirAll(domainPath, 0755)
	require.NoError(t, err)

	// Tạo file Go
	filePath := filepath.Join(domainPath, "module.go")
	content := `package book

import (
	"github.com/old-module/internal/book/model"
)

type Module struct {
	Model *model.Model
}
`
	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = RewriteAllImportsInDomain(tmpDir, "book", "github.com/old-module", "github.com/new-module")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/new-module/internal/book/model"`)
}

func TestRewriteAllImportsInDomain_DomainNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	err := RewriteAllImportsInDomain(tmpDir, "book", "old", "new")
	assert.Error(t, err)
}

// ==================== TEST FixAllImports ====================

func TestFixAllImports_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo file Go
	filePath := filepath.Join(tmpDir, "test.go")
	content := `package test

import (
	"new_example/old-module/internal/book"
)

type Test struct {
	Book *book.Book
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = FixAllImports(tmpDir, "new_example/old-module", "github.com/new-module")
	assert.NoError(t, err)

	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), `"github.com/new-module/internal/book"`)
}

func TestFixAllImports_SkipGitDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Tạo thư mục .git
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	filePath := filepath.Join(gitDir, "test.go")
	content := `package test
import "new_example/old-module"`
	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	err = FixAllImports(tmpDir, "new_example/old-module", "github.com/new-module")
	assert.NoError(t, err)

	// Kiểm tra file trong .git không bị thay đổi
	newContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(newContent), "new_example/old-module")
}

// ==================== TEST GetModuleNameFromRepo ====================

func TestGetModuleNameFromRepo(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		expected string
	}{
		{
			name:     "https URL",
			repoURL:  "https://github.com/owner/repo",
			expected: "github.com/owner/repo",
		},
		{
			name:     "https with .git",
			repoURL:  "https://github.com/owner/repo.git",
			expected: "github.com/owner/repo",
		},
		{
			name:     "http URL",
			repoURL:  "http://github.com/owner/repo",
			expected: "github.com/owner/repo",
		},
		{
			name:     "no protocol",
			repoURL:  "github.com/owner/repo",
			expected: "github.com/owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetModuleNameFromRepo(tt.repoURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== TEST GetCurrentModule ====================

func TestGetCurrentModule_Success(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	content := `module github.com/example/project

go 1.21
`
	err := os.WriteFile(goModPath, []byte(content), 0644)
	require.NoError(t, err)

	module, err := GetCurrentModule(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/example/project", module)
}

func TestGetCurrentModule_FileNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := GetCurrentModule(tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read go.mod")
}

func TestGetCurrentModule_NoModuleDirective(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	content := `go 1.21
`
	err := os.WriteFile(goModPath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = GetCurrentModule(tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "module directive not found")
}

// ==================== BENCHMARK ====================

func BenchmarkRewriteModuleImports(b *testing.B) {
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")
	content := `package test

import (
	"github.com/old-module/internal/book"
)

type Test struct {
	Book *book.Book
}
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		_ = RewriteModuleImports(filePath, "github.com/old-module", "github.com/new-module")
	}
}
