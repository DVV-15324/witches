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
		cfg      *Config
		expected string
	}{
		{
			name: "PostgreSQL",
			cfg: &Config{
				DBDriver:   "postgres",
				DBUser:     "user",
				DBPassword: "pass",
				DBHost:     "localhost",
				DBPort:     5432,
				DBName:     "db",
			},
			expected: "postgres://user:pass@localhost:5432/db?sslmode=disable",
		},
		{
			name: "MySQL",
			cfg: &Config{
				DBDriver:   "mysql",
				DBUser:     "user",
				DBPassword: "pass",
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "db",
			},
			expected: "mysql://user:pass@tcp(localhost:3306)/db",
		},
		{
			name: "SQL Server",
			cfg: &Config{
				DBDriver:   "sqlserver",
				DBUser:     "user",
				DBPassword: "pass",
				DBHost:     "localhost",
				DBPort:     1433,
				DBName:     "db",
			},
			expected: "sqlserver://user:pass@localhost:1433?database=db",
		},
		{
			name: "Unknown driver",
			cfg: &Config{
				DBDriver: "unknown",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDatabaseURL(tt.cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ... (Giữ nguyên phần TestGetCurrentPath, TestGetMigrationsPath, TestGetFrameworkPath ở dưới)
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
