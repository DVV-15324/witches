package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDatabaseURL(t *testing.T) {
	tests := []struct {
		name   string
		driver string
		dbURL  string
		want   string
	}{
		{
			name:   "mysql",
			driver: "mysql",
			dbURL:  "user:pass@localhost:3306/db",
			want:   "mysql://user:pass@localhost:3306/db",
		},
		{
			name:   "postgres",
			driver: "postgres",
			dbURL:  "user:pass@localhost:5432/db",
			want:   "postgres://user:pass@localhost:5432/db",
		},
		{
			name:   "sqlserver",
			driver: "sqlserver",
			dbURL:  "user:pass@localhost:1433?database=db",
			want:   "sqlserver://user:pass@localhost:1433?database=db",
		},
		{
			name:   "mysql without prefix",
			driver: "mysql",
			dbURL:  "user:pass@tcp(localhost:3306)/db",
			want:   "mysql://user:pass@tcp(localhost:3306)/db",
		},
		{
			name:   "unknown driver",
			driver: "unknown",
			dbURL:  "user:pass@localhost/db",
			want:   "unknown://user:pass@localhost/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildDatabaseURL(tt.driver, tt.dbURL)
			assert.Equal(t, tt.want, got)
		})
	}
}
