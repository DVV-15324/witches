package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDomainNameProcessing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "user", "User"},
		{"uppercase", "USER", "User"},
		{"mixed case", "UsEr", "User"},
		{"with spaces", " user ", "User"},
		{"multiple words", "user profile", "Userprofile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domainName := strings.TrimSpace(tt.input)
			domainName = strings.ToLower(domainName)
			domainName = strings.ReplaceAll(domainName, " ", "")
			domainNameCap := cases.Title(language.English).String(domainName)
			domainNameCap = strings.ReplaceAll(domainNameCap, " ", "")

			assert.Equal(t, tt.expected, domainNameCap)
		})
	}
}

func TestDomainConfig(t *testing.T) {
	config := ModuleConfig{
		NameCap:     "User",
		Name:        "user",
		ProjectName: "github.com/example/myproject",
	}

	assert.Equal(t, "github.com/example/myproject", config.GetProjectName())
	assert.Equal(t, "user", config.Name)
	assert.Equal(t, "User", config.NameCap)
}

func TestAddDomain_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	moduleName := "github.com/example/myproject"

	// Tạo go.mod
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module "+moduleName),
		0644,
	))

	// Tạo các thư mục cần thiết
	dirs := []string{
		"cmd/server/routers",
		"internal/shared/domain",
		"internal/shared/utils",
	}
	for _, dir := range dirs {
		require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, dir), 0755))
	}

	// Tạo modules.go
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
		filepath.Join(tmpDir, "cmd", "server", "routers", "modules.go"),
		[]byte(modulesContent),
		0644,
	))

	// Tạo routers.go
	routersContent := `package routers

func initModule(modules *Modules) {
	// Routes will be registered here
}
`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "cmd", "server", "routers", "routers.go"),
		[]byte(routersContent),
		0644,
	))

	// Tạo key_object.go
	keyContent := `package utils

var (
	ObjectUser int64 = 1
)
`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "internal", "shared", "utils", "key_object.go"),
		[]byte(keyContent),
		0644,
	))

	// Chạy AddDomain
	AddModule(tmpDir, moduleName, "book", "postgres")

	// Kiểm tra các file đã được tạo
	expectedFiles := []string{
		"internal/book/module.go",
		"internal/book/model/model.go",
		"internal/book/handler/handler.go",
		"internal/book/repository/repository.go",
		"internal/book/usecase/usecase.go",
		"internal/book/migrate/migrations/1_create_table.up.sql",
		"internal/book/migrate/migrations/1_drop_table.down.sql",
		"internal/shared/domain/book.go",
	}

	for _, relPath := range expectedFiles {
		actualPath := filepath.Join(tmpDir, relPath)
		_, err := os.Stat(actualPath)
		assert.NoError(t, err, "file should exist: %s", relPath)
	}

	// Golden test cho các file đã sửa
	filesToCheck := []string{
		"cmd/server/routers/modules.go",
		"cmd/server/routers/routers.go",
		"internal/shared/utils/key_object.go",
	}

	for _, relPath := range filesToCheck {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err)
		assertGolden(t, "add_domain", string(actual), relPath)
	}
}
func TestAddModule_MySQL_Golden(t *testing.T) {
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

	err := AddModule(tmpDir, moduleName, "book", "mysql")
	assert.NoError(t, err)

	// Kiểm tra migration files với MySQL
	_, err = os.Stat(filepath.Join(tmpDir, "internal", "book", "migrate", "migrations", "1_create_table.up.sql"))
	assert.NoError(t, err)
}

func TestAddModule_UnsupportedDriver(t *testing.T) {
	tmpDir := t.TempDir()
	moduleName := "github.com/example/project"

	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module "+moduleName),
		0644,
	))

	err := AddModule(tmpDir, moduleName, "book", "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database")
}
