package utils

import (
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
