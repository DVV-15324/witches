package utils

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetMigrationsURL(t *testing.T) {
	url := GetMigrationsURL()
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "migrate")
	assert.Contains(t, url, "migrations")
}

func TestGetFrameworkPath(t *testing.T) {
	path := GetFrameworkPath()
	assert.NotEmpty(t, path)
}
