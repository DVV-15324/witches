package template

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

// helper: đọc hoặc ghi golden file – luôn dùng thư mục testdata của source
func readOrUpdateGolden(t *testing.T, name string, actual []byte) []byte {
	// Lấy đường dẫn tuyệt đối của file golden_test.go
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller")
	}
	baseDir := filepath.Dir(filename)
	goldenPath := filepath.Join(baseDir, "testdata", name+".golden")

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, actual, 0644); err != nil {
			t.Fatal(err)
		}
		return actual
	}
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden file %s not found, run with -update", goldenPath)
		}
		t.Fatal(err)
	}
	return expected
}

// ------------------------------------------------------------
// Test AddGoService với golden
// ------------------------------------------------------------
func TestAddGoService_Golden(t *testing.T) {
	tmpDir := t.TempDir()

	// tạo go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module github.com/example/test"), 0644); err != nil {
		t.Fatal(err)
	}

	// tạo internal/shared/utils/key_object.go (cần cho updateKeyObject)
	utilsDir := filepath.Join(tmpDir, "internal", "shared", "utils")
	if err := os.MkdirAll(utilsDir, 0755); err != nil {
		t.Fatal(err)
	}
	keyContent := `package utils

var (
	ObjectUser uint = 10
)
`
	if err := os.WriteFile(filepath.Join(utilsDir, "key_object.go"), []byte(keyContent), 0644); err != nil {
		t.Fatal(err)
	}

	// gọi hàm AddGoService với service "book"
	AddGoService(tmpDir, "github.com/example/test", "book")

	// danh sách file cần so sánh với golden
	expectedFiles := []string{
		"internal/book-service/handler/handler.go",
		"internal/book-service/entity/entity.go",
		"internal/book-service/usecase/usecase.go",
		"internal/shared/model/book.go",
	}

	for _, relPath := range expectedFiles {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err)

		expected := readOrUpdateGolden(t, relPath, actual)
		assert.Equal(t, string(expected), string(actual), "file %s mismatch", relPath)
	}
}

// ------------------------------------------------------------
// Test createProjectStructure với golden (kiểm tra một vài file đại diện)
// ------------------------------------------------------------
func TestCreateProjectStructure_Golden(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origWd)

	require.NoError(t, os.Chdir(tmpDir))

	config := ProjectConfig{
		ModuleName: "github.com/example/project",
	}

	// test với postgres
	err = createProjectStructure(config, "postgres")
	require.NoError(t, err)

	// danh sách file đại diện để so với golden
	importantFiles := []string{
		"main.go",
		"go.mod",
		"cmd/root.go",
		"internal/shared/utils/key_object.go",
		"internal/auth-service/handler/handler.go",
		"pkg/redis/client.go",
		"migrate/migrations/1_create_table.up.sql",
		"migrate/migrations/1_drop_table.down.sql",
	}

	for _, relPath := range importantFiles {
		actualPath := filepath.Join(tmpDir, relPath)
		actual, err := os.ReadFile(actualPath)
		require.NoError(t, err)

		expected := readOrUpdateGolden(t, relPath, actual)
		assert.Equal(t, string(expected), string(actual), "file %s mismatch", relPath)
	}
}

func TestCreateGoArcRefresh(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	projectName := "github.com/example/myproject"
	dbType := "postgres"

	// Capture stdout/stderr để tránh làm phiền output test, hoặc không capture cũng được
	// Chúng ta có thể capture để kiểm tra log, nhưng không cần thiết

	CreateGoArcRefresh(projectName, dbType)

	// Kiểm tra file đã được tạo (dựa trên createProjectStructure)
	importantFiles := []string{
		"main.go",
		"go.mod",
		"cmd/root.go",
		"internal/shared/utils/key_object.go",
	}
	for _, f := range importantFiles {
		assert.FileExists(t, filepath.Join(tmpDir, f), "missing file: %s", f)
	}
}
