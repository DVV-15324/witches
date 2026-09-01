package template

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var githubAPIBase = "https://api.github.com"

func TestAddGoDomainFromLink_Golden(t *testing.T) {
	// Tạo mock GitHub API server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/owner/repo":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"default_branch": "main"}`))
		case r.URL.Path == "/repos/owner/repo/git/trees/main":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"tree": [
					{"path": "internal/book/module.go", "type": "blob"},
					{"path": "internal/book/model/model.go", "type": "blob"},
					{"path": "internal/book/handler/handler.go", "type": "blob"},
					{"path": "internal/book/repository/repository.go", "type": "blob"},
					{"path": "internal/book/usecase/usecase.go", "type": "blob"},
					{"path": "internal/book/migrate/migrations/1_create_table.up.sql", "type": "blob"},
					{"path": "internal/book/migrate/migrations/1_drop_table.down.sql", "type": "blob"},
					{"path": "internal/shared/domain/book.go", "type": "blob"}
				]
			}`))
		case strings.Contains(r.URL.Path, "/raw.githubusercontent.com/owner/repo/main/internal/"):
			fileName := filepath.Base(r.URL.Path)
			switch fileName {
			case "module.go":
				w.Write([]byte(`package book
type BookModule struct{}`))
			case "model.go":
				w.Write([]byte(`package model
type Book struct{ID string}`))
			case "handler.go":
				w.Write([]byte(`package handler
type BookHandler struct{}`))
			case "repository.go":
				w.Write([]byte(`package repository
type BookRepository struct{}`))
			case "usecase.go":
				w.Write([]byte(`package usecase
type BookUsecase struct{}`))
			case "1_create_table.up.sql":
				w.Write([]byte(`CREATE TABLE books (id SERIAL PRIMARY KEY, title TEXT);`))
			case "1_drop_table.down.sql":
				w.Write([]byte(`DROP TABLE books;`))
			case "book.go":
				w.Write([]byte(`package domain
const BookCollection = "books"`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	//  Override GitHub API base để dùng mock server
	originalBase := githubAPIBase
	githubAPIBase = mockServer.URL
	defer func() {
		githubAPIBase = originalBase
	}()

	tmpDir := t.TempDir()
	moduleName := "github.com/example/project"

	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module "+moduleName),
		0644,
	))

	// Tạo routers
	routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
	require.NoError(t, os.MkdirAll(routersDir, 0755))

	modulesContent := `package routers

type Modules struct {
	User *user.UserModule
}

func InitModules() *Modules {
	return &Modules{
		User: &user.UserModule{},
	}
}
`
	require.NoError(t, os.WriteFile(
		filepath.Join(routersDir, "modules.go"),
		[]byte(modulesContent),
		0644,
	))

	routersContent := `package routers

func initModule(modules *Modules) {
}
`
	require.NoError(t, os.WriteFile(
		filepath.Join(routersDir, "routers.go"),
		[]byte(routersContent),
		0644,
	))

	// Tạo utils
	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	require.NoError(t, os.MkdirAll(utilsDir, 0755))

	keyContent := `package utils

var (
	ObjectUser int64 = 1
)
`
	require.NoError(t, os.WriteFile(
		filepath.Join(utilsDir, "key_object.go"),
		[]byte(keyContent),
		0644,
	))

	// Chạy AddGoDomainFromLink với mock server
	err := AddGoDomainFromLink(
		tmpDir,
		moduleName,
		"book",
		"https://github.com/DVV-15324/witches-book",
	)
	assert.NoError(t, err)

	// Golden test cho các file đã tạo/sửa
	filesToCheck := []string{
		"internal/book/module.go",
		"internal/book/model/model.go",
		"internal/book/handler/handler.go",
		"internal/book/repository/repository.go",
		"internal/book/usecase/usecase.go",
		"internal/book/migrate/migrations/1_create_table.up.sql",
		"internal/book/migrate/migrations/1_drop_table.down.sql",
		"internal/shared/domain/book.go",
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
		"internal/shared/utils/key_object.go",
	}

	for _, relPath := range filesToCheck {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err, "file should exist: %s", relPath)
		assertGolden(t, "add_domain_link", string(actual), relPath)
	}
}

// ==================== TEST LỖI ====================

func TestAddGoDomainFromLink_InvalidRepo(t *testing.T) {
	tmpDir := t.TempDir()
	moduleName := "github.com/example/project"

	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module "+moduleName),
		0644,
	))

	routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
	require.NoError(t, os.MkdirAll(routersDir, 0755))

	modulesContent := `package routers
type Modules struct{}
func InitModules() *Modules { return &Modules{} }
`
	require.NoError(t, os.WriteFile(
		filepath.Join(routersDir, "modules.go"),
		[]byte(modulesContent),
		0644,
	))

	routersContent := `package routers
func initModule(modules *Modules) {}
`
	require.NoError(t, os.WriteFile(
		filepath.Join(routersDir, "routers.go"),
		[]byte(routersContent),
		0644,
	))

	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	require.NoError(t, os.MkdirAll(utilsDir, 0755))

	keyContent := `package utils
var (ObjectUser int64 = 1)
`
	require.NoError(t, os.WriteFile(
		filepath.Join(utilsDir, "key_object.go"),
		[]byte(keyContent),
		0644,
	))

	err := AddGoDomainFromLink(tmpDir, moduleName, "book", "https://invalid-url")
	assert.Error(t, err)
}

func TestAddGoDomainFromLink_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	err := AddGoDomainFromLink(tmpDir, "github.com/example/project", "book", "https://github.com/owner/repo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get target module")
}

// ==================== BENCHMARK ====================

func BenchmarkAddGoDomainFromLink(b *testing.B) {
	tmpDir := b.TempDir()
	moduleName := "github.com/example/project"

	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module "+moduleName), 0644)

	routersDir := filepath.Join(tmpDir, "cmd", "server", "routers")
	os.MkdirAll(routersDir, 0755)

	modulesContent := `package routers
type Modules struct {
	User *user.UserModule
}
func InitModules() *Modules { return &Modules{User: &user.UserModule{}} }
`
	os.WriteFile(filepath.Join(routersDir, "modules.go"), []byte(modulesContent), 0644)

	routersContent := `package routers
func initModule(modules *Modules) {}
`
	os.WriteFile(filepath.Join(routersDir, "routers.go"), []byte(routersContent), 0644)

	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	os.MkdirAll(utilsDir, 0755)

	keyContent := `package utils
var (ObjectUser int64 = 1)
`
	os.WriteFile(filepath.Join(utilsDir, "key_object.go"), []byte(keyContent), 0644)

	b.ResetTimer()
	for b.Loop() {
		_ = AddGoDomainFromLink(tmpDir, moduleName, "book", "https://github.com/owner/repo")
	}
}
