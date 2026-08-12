package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	utils "github.com/DVV-15324/witches/cmd/utils"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper: Setup MySQL container
func setupMySQLContainer(t *testing.T) (string, func()) {
	ctx := context.Background()

	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0"),
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx, "parseTime=true")
	require.NoError(t, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return connStr, cleanup
}

// Helper: Setup PostgreSQL container
func setupPostgresContainer(t *testing.T) (string, func()) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return connStr, cleanup
}

// Test GenerateDBURL với MySQL container thật
func TestGenerateDBURL_MySQL_WithContainer(t *testing.T) {
	// 1. Tạo MySQL container
	connStr, cleanup := setupMySQLContainer(t)
	defer cleanup()

	// 2. Kết nối đến DB
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	// 3. Test GenerateDBURL
	url, err := GenerateDBURL(
		"mysql",
		"testuser",
		"testpass",
		"localhost",
		"testdb",
		3306,
	)
	assert.NoError(t, err)
	assert.Contains(t, url, "testuser:testpass@tcp(localhost:3306)/testdb")
}

// Test GenerateDBURL với PostgreSQL container thật
func TestGenerateDBURL_Postgres_WithContainer(t *testing.T) {
	// 1. Tạo PostgreSQL container
	connStr, cleanup := setupPostgresContainer(t)
	defer cleanup()

	// 2. Kết nối đến DB
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	// 3. Test GenerateDBURL
	url, err := GenerateDBURL(
		"postgres",
		"testuser",
		"testpass",
		"localhost",
		"testdb",
		5432,
	)
	assert.NoError(t, err)
	assert.Contains(t, url, "testuser:testpass@localhost:5432/testdb")
}

// Test WitchesDBURL với MySQL container thật (không cần file witches.env)
func TestWitchesDBURL_MySQL_WithContainer(t *testing.T) {
	// 1. Tạo MySQL container
	connStr, cleanup := setupMySQLContainer(t)
	defer cleanup()

	// 2. Kết nối đến DB
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	// 3. Test WitchesDBURL - chỉ kiểm tra URL, không cần file
	config := &utils.Config{
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     3306,
	}

	err = WitchesDBURL("mysql", config)
	assert.NoError(t, err)
}

// Test WitchesDBURL với PostgreSQL container thật
func TestWitchesDBURL_Postgres_WithContainer(t *testing.T) {
	// 1. Tạo PostgreSQL container
	connStr, cleanup := setupPostgresContainer(t)
	defer cleanup()

	// 2. Kết nối đến DB
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	// 3. Test WitchesDBURL
	config := &utils.Config{
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBHost:     "localhost",
		DBName:     "testdb",
		DBPort:     5432,
	}

	err = WitchesDBURL("postgres", config)
	assert.NoError(t, err)
}

// Test kết nối thật với container
func TestDatabaseConnection_WithContainer(t *testing.T) {
	tests := []struct {
		name      string
		driver    string
		setup     func(*testing.T) (string, func())
		generator func() (string, error)
	}{
		{
			name:   "MySQL",
			driver: "mysql",
			setup:  setupMySQLContainer,
			generator: func() (string, error) {
				return GenerateDBURL("mysql", "testuser", "testpass", "localhost", "testdb", 3306)
			},
		},
		{
			name:   "PostgreSQL",
			driver: "postgres",
			setup:  setupPostgresContainer,
			generator: func() (string, error) {
				return GenerateDBURL("postgres", "testuser", "testpass", "localhost", "testdb", 5432)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Setup container
			connStr, cleanup := tt.setup(t)
			defer cleanup()

			// 2. Connect DB
			db, err := sql.Open(tt.driver, connStr)
			require.NoError(t, err)
			defer db.Close()

			err = db.Ping()
			require.NoError(t, err)

			// 3. Test generator
			url, err := tt.generator()
			assert.NoError(t, err)
			assert.NotEmpty(t, url)
		})
	}
}

func BenchmarkGenerateDBURL_MySQL(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		GenerateDBURL(
			"mysql",
			"testuser",
			"testpass",
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
			"testuser",
			"testpass",
			"localhost",
			"testdb",
			5432,
		)
	}
}

func BenchmarkGenerateDBURL_SQLServer(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		GenerateDBURL(
			"sqlserver",
			"sa",
			"testpass",
			"localhost",
			"testdb",
			1433,
		)
	}
}
