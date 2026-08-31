package template

import (
	"os"
	"os/exec"
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
	CreateCaptainGoArc(projectName, "postgres")

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
		assertGolden(t, "captain_init", string(actual), relPath)
	}
}

func TestCreateCaptainGoArc_InvalidDriver(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	//BẮT LỖI: Sử dụng recover để bắt panic hoặc chạy trong subprocess
	if os.Getenv("TEST_SUBPROCESS_INVALID_DRIVER") == "1" {
		CreateCaptainGoArc("testproject", "invalid")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCreateCaptainGoArc_InvalidDriver")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_INVALID_DRIVER=1")

	err = cmd.Run()
	assert.Error(t, err, "Should exit with error for invalid driver")
}
