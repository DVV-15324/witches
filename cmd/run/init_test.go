package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWitchesInit_Success(t *testing.T) {
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test-module"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	WitchesInit("postgres")

	expectedDirs := []string{
		"cmd",
		"internal",
		"pkg",
		"migrate/migrations",
	}
	for _, dir := range expectedDirs {
		_, err := os.Stat(dir)
		assert.NoError(t, err, "Directory %s should be created", dir)
	}

	_, err = os.Stat("main.go")
	assert.NoError(t, err, "main.go should exist at root")
}

// ✅ Cách 1: Dùng subprocess để test log.Fatal
func TestWitchesInit_InvalidDriver(t *testing.T) {
	// Nếu đang chạy subprocess
	if os.Getenv("TEST_SUBPROCESS_INIT") == "1" {
		WitchesInit("invalid-driver")
		return
	}

	// Chạy subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestWitchesInit_InvalidDriver")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_INIT=1")

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
