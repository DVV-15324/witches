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
	assert.Contains(t, err.Error(), "no module files")
}

func TestRewriteModuleImports_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.go")

	content := `package test

import (
    "github.com/old-module/internal/book"
    "github.com/old-module/internal/shared/domain"
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
	assert.Contains(t, string(newContent), `"github.com/new-module/internal/shared/domain"`)
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
func TestNormalizeImportPath(t *testing.T) {
	tests := []struct {
		name         string
		oldPath      string
		sourceModule string
		targetModule string
		want         string
	}{
		// ==================== TRƯỜNG HỢP GIỮ NGUYÊN ====================
		{
			name:         "skip external pkg from source module",
			oldPath:      "github.com/DVV-15324/witches-book/pkg/redis",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/DVV-15324/witches-book/pkg/redis",
		},
		{
			name:         "skip external pkg from source module with nested path",
			oldPath:      "github.com/DVV-15324/witches-book/pkg/redis/client",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/DVV-15324/witches-book/pkg/redis/client",
		},
		{
			name:         "skip external library github.com",
			oldPath:      "github.com/gin-gonic/gin",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/gin-gonic/gin",
		},
		{
			name:         "skip external library gitlab.com",
			oldPath:      "gitlab.com/some/project",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "gitlab.com/some/project",
		},
		{
			name:         "skip external library golang.org",
			oldPath:      "golang.org/x/text/cases",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "golang.org/x/text/cases",
		},
		{
			name:         "skip external library go.uber.org",
			oldPath:      "go.uber.org/zap",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "go.uber.org/zap",
		},
		{
			name:         "skip external library google.golang.org",
			oldPath:      "google.golang.org/grpc",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "google.golang.org/grpc",
		},
		{
			name:         "skip external library gopkg.in",
			oldPath:      "gopkg.in/yaml.v3",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "gopkg.in/yaml.v3",
		},
		{
			name:         "skip when oldPath is empty",
			oldPath:      "",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "",
		},
		{
			name:         "skip when targetModule is empty",
			oldPath:      "internal/book/handler",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "",
			want:         "internal/book/handler",
		},

		// ==================== TRƯỜNG HỢP SỬA (internal/) ====================
		{
			name:         "replace internal import from source module",
			oldPath:      "github.com/DVV-15324/witches-book/internal/book/handler",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/book/handler",
		},
		{
			name:         "replace internal import with nested path",
			oldPath:      "github.com/DVV-15324/witches-book/internal/book/repository",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/book/repository",
		},
		{
			name:         "replace internal import with usecase alias",
			oldPath:      "github.com/DVV-15324/witches-book/internal/book/usecase",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/book/usecase",
		},
		{
			name:         "replace internal import with shared domain",
			oldPath:      "github.com/DVV-15324/witches-book/internal/shared/domain",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/shared/domain",
		},
		{
			name:         "replace internal import with utils",
			oldPath:      "github.com/DVV-15324/witches-book/internal/shared/utils",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/shared/utils",
		},

		// ==================== TRƯỜNG HỢP SỬA (cmd/) ====================
		{
			name:         "replace cmd import from source module",
			oldPath:      "github.com/DVV-15324/witches-book/cmd/server/config",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/cmd/server/config",
		},
		{
			name:         "replace cmd import with routers",
			oldPath:      "github.com/DVV-15324/witches-book/cmd/server/routers",
			sourceModule: "github.com/DVV-15324/witches-book",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/cmd/server/routers",
		},

		// ==================== TRƯỜNG HỢP SỬA (import tương đối) ====================
		{
			name:         "add target module to internal relative import",
			oldPath:      "internal/book/handler",
			sourceModule: "",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/book/handler",
		},
		{
			name:         "add target module to cmd relative import",
			oldPath:      "cmd/server/config",
			sourceModule: "",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/cmd/server/config",
		},
		{
			name:         "add target module to pkg relative import",
			oldPath:      "pkg/utils",
			sourceModule: "",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/pkg/utils",
		},
		{
			name:         "add target module to internal with nested path",
			oldPath:      "internal/shared/domain",
			sourceModule: "",
			targetModule: "github.com/myapp/backend",
			want:         "github.com/myapp/backend/internal/shared/domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeImportPath(tt.oldPath, tt.sourceModule, tt.targetModule)
			if got != tt.want {
				t.Errorf("NormalizeImportPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
