package template

import (
	"os"

	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCaptainGoArc_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	projectName := "testproject"
	err = CreateTemplateGoArc(projectName, "postgres")
	assert.NoError(t, err)

	// Kiểm tra các file đã được tạo
	expectedFiles := []string{
		"main.go",
		"README.md",
		"go.mod",
		"cmd/root.go",
		"cmd/server/config/config.go",
		"cmd/server/routers/composer.go",
		"cmd/server/routers/routers.go",
		"cmd/server/routers/modules.go",
		"cmd/server/core/core.go",
		"internal/shared/middleware/limit.go",
		"internal/shared/middleware/timing.go",
		"internal/shared/utils/dummy.go",
		"internal/shared/utils/key_object.go",
		"internal/shared/utils/uid.go",
		"migrate/migrations/1_create_table.up.sql",
		"migrate/migrations/1_drop_table.down.sql",
		"pkg/redis/client.go",
	}

	for _, relPath := range expectedFiles {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err, "file should exist: %s", relPath)
		assertGolden(t, "create_project", string(actual), relPath)
	}
}

func TestCreateGoArc_InvalidDriver(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// CreateTemplateGoArc trả về error, test trực tiếp
	err = CreateTemplateGoArc("testproject", "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database")
}
func TestCreateTemplateGoArc_MySQL_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = CreateTemplateGoArc("testproject", "mysql")
	assert.NoError(t, err)

	// Kiểm tra migration files với MySQL
	_, err = os.Stat("migrate/migrations/1_create_table.up.sql")
	assert.NoError(t, err)
}

func TestCreateTemplateGoArc_MSSQL_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = CreateTemplateGoArc("testproject", "mssql")
	assert.NoError(t, err)

	_, err = os.Stat("migrate/migrations/1_create_table.up.sql")
	assert.NoError(t, err)
}
func TestModuleConfig_GetProjectNameTemplate(t *testing.T) {
	config := TemplateConfig{
		ProjectName: "github.com/example/project",
	}

	assert.Equal(t, "github.com/example/project", config.GetProjectName())
}
