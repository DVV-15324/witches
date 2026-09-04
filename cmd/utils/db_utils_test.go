package utils

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		dbURL    string
		expected string
	}{
		{
			name:     "PostgreSQL",
			driver:   "postgres",
			dbURL:    "user:pass@localhost:5432/db",
			expected: "postgres://user:pass@localhost:5432/db",
		},
		{
			name:     "MySQL",
			driver:   "mysql",
			dbURL:    "user:pass@tcp(localhost:3306)/db",
			expected: "mysql://user:pass@tcp(localhost:3306)/db",
		},
		{
			name:     "SQL Server",
			driver:   "sqlserver",
			dbURL:    "user:pass@localhost:1433/db",
			expected: "sqlserver://user:pass@localhost:1433/db",
		},
		{
			name:     "Empty driver",
			driver:   "",
			dbURL:    "user:pass@localhost:5432/db",
			expected: "user:pass@localhost:5432/db",
		},
		{
			name:     "Unknown driver",
			driver:   "unknown",
			dbURL:    "user:pass@localhost:5432/db",
			expected: "unknown://user:pass@localhost:5432/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDatabaseURL(tt.driver, tt.dbURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCurrentPath_GetwdFail(t *testing.T) {
	if os.Getenv("TEST_GET_CURRENT_PATH_FAIL") == "1" {
		getwdUtils = func() (string, error) {
			return "", errors.New("mock getwd error")
		}

		GetCurrentPath()
		return
	}

	cmd := exec.Command(
		os.Args[0],
		"-test.run=TestGetCurrentPath_GetwdFail",
	)

	cmd.Env = append(
		os.Environ(),
		"TEST_GET_CURRENT_PATH_FAIL=1",
	)

	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode())
}
func TestGetMigrationsPath_GetwdFail(t *testing.T) {
	if os.Getenv("TEST_GET_MIGRATIONS_PATH_FAIL") == "1" {
		getwdUtils = func() (string, error) {
			return "", errors.New("mock getwd error")
		}

		GetMigrationsPath()
		return
	}

	cmd := exec.Command(
		os.Args[0],
		"-test.run=TestGetMigrationsPath_GetwdFail",
	)

	cmd.Env = append(
		os.Environ(),
		"TEST_GET_MIGRATIONS_PATH_FAIL=1",
	)

	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode())
}
func TestGetFrameworkPath_ExecutableFail(t *testing.T) {
	if os.Getenv("TEST_GET_FRAMEWORK_PATH_FAIL") == "1" {
		executableUtils = func() (string, error) {
			return "", errors.New("mock executable error")
		}

		GetFrameworkPath()
		return
	}

	cmd := exec.Command(
		os.Args[0],
		"-test.run=TestGetFrameworkPath_ExecutableFail",
	)

	cmd.Env = append(
		os.Environ(),
		"TEST_GET_FRAMEWORK_PATH_FAIL=1",
	)

	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode())
}
