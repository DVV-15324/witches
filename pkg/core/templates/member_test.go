package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMemberGoArc_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	projectName := "testmember"
	CreateMemberGoArc(projectName, "postgres")

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
		"pkg/redis/client.go",
	}

	for _, relPath := range expectedFiles {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err, "file should exist: %s", relPath)
		assertGolden(t, "member_init", string(actual), relPath)
	}
}
