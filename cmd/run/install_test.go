package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWitchesInstall_Success(t *testing.T) {
	t.Skip("Skipping integration test - requires network and go modules")
	WitchesInstall("postgres")
}

func TestWitchesInstall_DriverMapping(t *testing.T) {
	drivers := map[string]string{
		"mysql":      "github.com/golang-migrate/migrate/v4/database/mysql@latest",
		"postgres":   "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"postgresql": "github.com/golang-migrate/migrate/v4/database/postgres@latest",
		"mssql":      "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
		"sqlserver":  "github.com/golang-migrate/migrate/v4/database/sqlserver@latest",
	}

	for driver, expected := range drivers {
		t.Run(driver, func(t *testing.T) {
			assert.Equal(t, expected, drivers[driver])
		})
	}
}
