package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentPath(t *testing.T) {
	path := GetCurrentPath()
	assert.NotEmpty(t, path)
	_, err := os.Stat(path)
	assert.NoError(t, err)
}

func TestGetMigrationsPath(t *testing.T) {
	path := GetMigrationsPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "migrate")
	assert.Contains(t, path, "migrations")
}

func TestGetMigrationsURL_Empty(t *testing.T) {
	url := GetMigrationsURL("")
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "migrate")
	assert.Contains(t, url, "migrations")
	// Kiểm tra không có dấu // dư thừa
	assert.NotContains(t, url, "//")
}

func TestGetMigrationsURL_Domain(t *testing.T) {
	// Tạo thư mục test
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo thư mục domain
	domainPath := filepath.Join(tmpDir, "internal", "book", "migrate", "migrations")
	err = os.MkdirAll(domainPath, 0755)
	require.NoError(t, err)

	url := GetMigrationsURL("book")
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "internal/book/migrate/migrations")
}

func TestGetMigrationsURL_CustomPath(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Tạo thư mục custom
	customPath := filepath.Join(tmpDir, "custom", "migrations")
	err = os.MkdirAll(customPath, 0755)
	require.NoError(t, err)

	url := GetMigrationsURL("./custom/migrations")
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "custom/migrations")
}

func TestGetMigrationsURL_DomainNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Domain chưa tồn tại, sẽ tự tạo
	url := GetMigrationsURL("newdomain")
	assert.NotEmpty(t, url)

	// SỬA: Kiểm tra đường dẫn đúng
	expected := filepath.ToSlash(filepath.Join(tmpDir, "internal", "newdomain", "migrate", "migrations"))
	assert.Equal(t, expected, url)

	// Kiểm tra thư mục đã được tạo
	_, err = os.Stat(filepath.Join(tmpDir, "internal", "newdomain", "migrate", "migrations"))
	assert.NoError(t, err, "Directory should be created")
}

func TestGetFrameworkPath(t *testing.T) {
	path := GetFrameworkPath()
	assert.NotEmpty(t, path)
}
