package database

import (
	"os"
	"path/filepath"
	"testing"

	utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDBURL_MySQL(t *testing.T) {
	url, err := GenerateDBURL(
		"mysql",
		"root",
		"password",
		"localhost",
		"testdb",
		3306,
	)

	assert.NoError(t, err)
	expected := "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(t, expected, url)
}

func TestGenerateDBURL_Postgres(t *testing.T) {
	url, err := GenerateDBURL(
		"postgres",
		"postgres",
		"password",
		"localhost",
		"testdb",
		5432,
	)

	assert.NoError(t, err)
	expected := "postgres:password@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expected, url)
}

func TestGenerateDBURL_Postgresql(t *testing.T) {
	url, err := GenerateDBURL(
		"postgresql",
		"postgres",
		"password",
		"localhost",
		"testdb",
		5432,
	)

	assert.NoError(t, err)
	expected := "postgres:password@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expected, url)
}

func TestGenerateDBURL_SQLServer(t *testing.T) {
	url, err := GenerateDBURL(
		"sqlserver",
		"sa",
		"password",
		"localhost",
		"testdb",
		1433,
	)

	assert.NoError(t, err)
	expected := "sa:password@localhost:1433?database=testdb&encrypt=disable"
	assert.Equal(t, expected, url)
}

func TestGenerateDBURL_MSSQL(t *testing.T) {
	url, err := GenerateDBURL(
		"mssql",
		"sa",
		"password",
		"localhost",
		"testdb",
		1433,
	)

	assert.NoError(t, err)
	expected := "sa:password@localhost:1433?database=testdb&encrypt=disable"
	assert.Equal(t, expected, url)
}

func TestGenerateDBURL_InvalidDriver(t *testing.T) {
	url, err := GenerateDBURL(
		"invalid",
		"user",
		"pass",
		"localhost",
		"db",
		3306,
	)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Error(), "unsupported database")
}

func TestWitchesDBURL_Success(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	err = WitchesDBURL("mysql", config)
	assert.NoError(t, err)

	envPath := filepath.Join(tmpDir, "witches.env")
	assert.FileExists(t, envPath)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "DB_URL")
	assert.Contains(t, string(content), "root:password@tcp(localhost:3306)/testdb")
}

func TestWitchesDBURL_WithPostgres(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "postgres",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     5432,
	}

	err = WitchesDBURL("postgres", config)
	assert.NoError(t, err)

	envPath := filepath.Join(tmpDir, "witches.env")
	assert.FileExists(t, envPath)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "postgres:password@localhost:5432/testdb")
}

func TestWitchesDBURL_WithSQLServer(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "sa",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     1433,
	}

	err = WitchesDBURL("sqlserver", config)
	assert.NoError(t, err)

	envPath := filepath.Join(tmpDir, "witches.env")
	assert.FileExists(t, envPath)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "sa:password@localhost:1433")
}

func TestWitchesDBURL_InvalidDriver(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	err = WitchesDBURL("invalid", config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database")
}

func TestWitchesDBURL_CreateDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	err = os.Chdir(subDir)
	require.NoError(t, err)

	config := &utils.Config{
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	err = WitchesDBURL("mysql", config)
	assert.NoError(t, err)

	envPath := filepath.Join(subDir, "witches.env")
	assert.FileExists(t, envPath)
}

func TestGenerateDBURL_AllDrivers(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		user     string
		password string
		host     string
		dbName   string
		port     int64
		expected string
	}{
		{
			name:     "mysql",
			driver:   "mysql",
			user:     "root",
			password: "pass",
			host:     "localhost",
			dbName:   "mydb",
			port:     3306,
			expected: "root:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name:     "postgres",
			driver:   "postgres",
			user:     "postgres",
			password: "pass",
			host:     "localhost",
			dbName:   "mydb",
			port:     5432,
			expected: "postgres:pass@localhost:5432/mydb?sslmode=disable",
		},
		{
			name:     "postgresql",
			driver:   "postgresql",
			user:     "postgres",
			password: "pass",
			host:     "localhost",
			dbName:   "mydb",
			port:     5432,
			expected: "postgres:pass@localhost:5432/mydb?sslmode=disable",
		},
		{
			name:     "sqlserver",
			driver:   "sqlserver",
			user:     "sa",
			password: "pass",
			host:     "localhost",
			dbName:   "mydb",
			port:     1433,
			expected: "sa:pass@localhost:1433?database=mydb&encrypt=disable",
		},
		{
			name:     "mssql",
			driver:   "mssql",
			user:     "sa",
			password: "pass",
			host:     "localhost",
			dbName:   "mydb",
			port:     1433,
			expected: "sa:pass@localhost:1433?database=mydb&encrypt=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := GenerateDBURL(
				tt.driver,
				tt.user,
				tt.password,
				tt.host,
				tt.dbName,
				tt.port,
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func BenchmarkGenerateDBURL_MySQL(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		GenerateDBURL(
			"mysql",
			"root",
			"password",
			"localhost",
			"testdb",
			3306,
		)
	}
}

func BenchmarkGenerateDBURL_Postgres(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		GenerateDBURL(
			"postgres",
			"postgres",
			"password",
			"localhost",
			"testdb",
			5432,
		)
	}
}

func BenchmarkWitchesDBURL(b *testing.B) {
	tmpDir := b.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(b, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)

	config := &utils.Config{
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	b.ResetTimer()
	for b.Loop() {
		WitchesDBURL("mysql", config)
	}
}
