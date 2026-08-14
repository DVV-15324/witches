package template

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateFS_ContainsFiles(t *testing.T) {
	// Kiểm tra embed FS có chứa các file template cần thiết
	entries, err := templateFS.ReadDir("template")
	require.NoError(t, err)
	assert.NotEmpty(t, entries)

	// Kiểm tra các thư mục con
	expectedDirs := []string{
		"cmd",
		"internal",
		"migrate",
		"pkg",
	}

	foundDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			foundDirs[entry.Name()] = true
		}
	}

	for _, dir := range expectedDirs {
		assert.True(t, foundDirs[dir], "Missing directory: %s", dir)
	}
}

func TestTemplateFS_ReadTemplateFile(t *testing.T) {
	// Kiểm tra đọc được file template cụ thể
	testFiles := []string{
		"template/main.go.tmpl",
		"template/cmd/root.go.tmpl",
		"template/internal/auth/handler/handler.go.tmpl",
	}

	for _, file := range testFiles {
		content, err := templateFS.ReadFile(file)
		assert.NoError(t, err, "Failed to read: %s", file)
		assert.NotEmpty(t, content, "File is empty: %s", file)
	}
}

func TestTemplateFS_AllTemplateFiles(t *testing.T) {
	// Đếm số lượng file template bằng fs.WalkDir
	var totalFiles int
	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})
	require.NoError(t, err)
	assert.Greater(t, totalFiles, 50, "Should have at least 50 template files")
}

func TestTemplateFS_CountTmplFiles(t *testing.T) {
	// Đếm số lượng file .tmpl
	var tmplFiles int
	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tmpl" {
			tmplFiles++
		}
		return nil
	})
	require.NoError(t, err)
	assert.Greater(t, tmplFiles, 40, "Should have at least 40 .tmpl files")
}

func TestCreateProjectStructure_Success(t *testing.T) {
	// Tạo temp dir và chuyển vào đó
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := ProjectConfig{
		ModuleName: "github.com/example/test-project",
	}

	// Test với từng loại database
	dbTypes := []string{"mysql", "postgres", "postgresql", "mssql", "sqlserver"}

	for _, dbType := range dbTypes {
		t.Run("db_"+dbType, func(t *testing.T) {
			// Tạo subdir cho mỗi test
			testDir := filepath.Join(tmpDir, dbType)
			err := os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			err = os.Chdir(testDir)
			require.NoError(t, err)

			err = createProjectStructure(config, dbType)
			assert.NoError(t, err)

			// Kiểm tra các file quan trọng đã được tạo
			criticalFiles := []string{
				"main.go",
				"go.mod",
				"cmd/root.go",
				"internal/auth/handler/handler.go",
				"internal/shared/utils/key_object.go",
				"migrate/migrations/1_create_table.up.sql",
				"migrate/migrations/1_drop_table.down.sql",
				"pkg/redis/client.go",
			}

			for _, file := range criticalFiles {
				assert.FileExists(t, filepath.Join(testDir, file), "Missing file: %s", file)
			}
		})
	}
}

func TestCreateProjectStructure_InvalidDB(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := ProjectConfig{
		ModuleName: "github.com/example/test",
	}

	// Test với DB không hỗ trợ
	err = createProjectStructure(config, "invalid_db")
	assert.Error(t, err)
}

func TestRenderTemplate_AllTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	var templateFiles []string
	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tmpl" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	require.NoError(t, err)

	for _, tmpl := range templateFiles {
		t.Run("render_"+filepath.Base(tmpl), func(t *testing.T) {
			dest := filepath.Join(tmpDir, "output", strings.TrimSuffix(filepath.Base(tmpl), ".tmpl"))
			destDir := filepath.Dir(dest)

			err := os.MkdirAll(destDir, 0755)
			require.NoError(t, err)

			// utils.RenderTemplate(templateFS, tmpDir, dest, tmpl, config)
		})
	}
}

func TestProjectConfig_GetModuleName(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		expected   string
	}{
		{
			name:       "github module",
			moduleName: "github.com/user/project",
			expected:   "github.com/user/project",
		},
		{
			name:       "local module",
			moduleName: "example.com/myapp",
			expected:   "example.com/myapp",
		},
		{
			name:       "empty module",
			moduleName: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ProjectConfig{ModuleName: tt.moduleName}
			assert.Equal(t, tt.expected, config.GetMuduleName())
		})
	}
}

func TestCreateProjectStructure_EmbedContents(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := ProjectConfig{
		ModuleName: "github.com/example/embed-test",
	}

	err = createProjectStructure(config, "mysql")
	require.NoError(t, err)

	// Kiểm tra nội dung file main.go
	mainContent, err := os.ReadFile("main.go")
	require.NoError(t, err)
	assert.Contains(t, string(mainContent), "package main")
	assert.Contains(t, string(mainContent), "github.com/example/embed-test")

	// Kiểm tra key_object.go
	keyContent, err := os.ReadFile("internal/shared/utils/key_object.go")
	require.NoError(t, err)
	assert.Contains(t, string(keyContent), "package utils")
}

func BenchmarkCreateProjectStructure(b *testing.B) {
	tmpDir := b.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(b, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)

	config := ProjectConfig{
		ModuleName: "github.com/example/bench",
	}

	b.ResetTimer()
	for b.Loop() {
		_ = createProjectStructure(config, "mysql")
		// Clean up after each iteration
		os.RemoveAll(tmpDir)
		os.MkdirAll(tmpDir, 0755)
		os.Chdir(tmpDir)
	}
}

func BenchmarkTemplateFS_WalkAll(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var count int
		fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				count++
			}
			return nil
		})
		_ = count
	}
}

func BenchmarkTemplateFS_ReadAllFiles(b *testing.B) {
	// Lấy danh sách file trước
	var files []string
	fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	b.ResetTimer()
	for b.Loop() {
		for _, f := range files {
			templateFS.ReadFile(f)
		}
	}
}
