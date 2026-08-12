// cmd/create/create_test.go
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWitchesCreate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	projectName := "test-project"
	WitchesCreate(projectName)

	projectPath := filepath.Join(tmpDir, projectName)
	_, err = os.Stat(projectPath)
	assert.NoError(t, err, "Project directory should be created")

	envPath := filepath.Join(projectPath, "witches.env")
	_, err = os.Stat(envPath)
	assert.NoError(t, err, "witches.env should be created")

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "DB_")
}

// ✅ Test với subprocess
func TestWitchesCreate_ProjectExists(t *testing.T) {
	// Kiểm tra nếu đang chạy subprocess
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		// Tạo project đã tồn tại
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		if err != nil {
			os.Exit(1)
		}

		projectName := "existing-project"
		err = os.Mkdir(projectName, 0755)
		if err != nil {
			os.Exit(1)
		}

		// Gọi hàm sẽ log.Fatal và exit
		WitchesCreate(projectName)
		return
	}

	// Chạy subprocess để test
	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesCreate_ProjectExists")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")

	err := cmd.Run()

	// ✅ log.Fatal sẽ exit với code 1
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	} else {
		t.Errorf("Expected ExitError, got %T", err)
	}
}

func TestWitchesCreate_EnvFileExists(t *testing.T) {
	// ✅ Dùng subprocess cho test này luôn vì cũng log.Fatal
	if os.Getenv("TEST_SUBPROCESS_ENV") == "1" {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		if err != nil {
			os.Exit(1)
		}

		projectName := "test-project"
		projectPath := filepath.Join(tmpDir, projectName)

		err = os.MkdirAll(projectPath, 0755)
		if err != nil {
			os.Exit(1)
		}

		envPath := filepath.Join(projectPath, "witches.env")
		err = os.WriteFile(envPath, []byte("existing env"), 0644)
		if err != nil {
			os.Exit(1)
		}

		WitchesCreate(projectName)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesCreate_EnvFileExists")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_ENV=1")

	err := cmd.Run()
	assert.Error(t, err, "Should exit with error")
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")
	}
}
