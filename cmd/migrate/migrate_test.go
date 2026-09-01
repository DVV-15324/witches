package cmd_migrate

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	// Tắt Ryuk để tránh lỗi
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	os.Setenv("TESTCONTAINERS_REUSE_ENABLE", "true")
	code := m.Run()
	os.Exit(code)
}

func setupTestWithPostgres(t *testing.T) (dbURL string, migrationPath string, cleanup func()) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrate", "migrations")
	migrationsDir = filepath.ToSlash(migrationsDir)
	err := os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err)

	upFile := filepath.Join(migrationsDir, "000001_test_migration.up.sql")
	downFile := filepath.Join(migrationsDir, "000001_test_migration.down.sql")

	upContent := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	downContent := `DROP TABLE IF EXISTS users;`

	err = os.WriteFile(upFile, []byte(upContent), 0644)
	require.NoError(t, err)
	err = os.WriteFile(downFile, []byte(downContent), 0644)
	require.NoError(t, err)

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	cleanup = func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
		_ = os.RemoveAll(tmpDir)
	}
	connStr = strings.TrimPrefix(connStr, "postgres://")
	return connStr, migrationsDir, cleanup
}

func TestWitchesMigrateUp(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	var tableExists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists, "Table 'users' should exist")
}

func TestWitchesMigrateVersion(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateVersion(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	var version int
	err = db.QueryRowContext(ctx, `
		SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1
	`).Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, 1, version, "Migration version should be 1")
}

func TestWitchesMigrateDown(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateDown(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	var tableExists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err)
	assert.False(t, tableExists, "Table 'users' should not exist")
}

func TestWitchesMigrateDrop(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateDrop(dbURL, "postgres", migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	var tableExists bool
	err = db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err)
	assert.False(t, tableExists, "Table 'users' should not exist")
}

func TestWitchesMigrateForce(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateForce(dbURL, "postgres", migrationPath, "1")
	require.NoError(t, err)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	var version int
	err = db.QueryRowContext(ctx, `
		SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1
	`).Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, 1, version, "Migration version should be forced to 1")
}

func TestWitchesMigrateUp_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	// Test với path không tồn tại
	err = WitchesMigrateUp("user:pass@localhost:5432/db", "postgres", "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate up failed")
}

func TestWitchesMigrateUp_InvalidDBURL(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	tmpDir := t.TempDir()
	err = WitchesMigrateUp("invalid_url", "postgres", tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate up failed")
}

func TestWitchesMigrateDown_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	err = WitchesMigrateDown("user:pass@localhost:5432/db", "postgres", "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate down failed")
}

func TestWitchesMigrateDrop_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	err = WitchesMigrateDrop("user:pass@localhost:5432/db", "postgres", "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate drop failed")
}

func TestWitchesMigrateForce_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	err = WitchesMigrateForce("user:pass@localhost:5432/db", "postgres", "/invalid/path", "1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate force failed")
}

func TestWitchesMigrateVersion_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	err = WitchesMigrateVersion("user:pass@localhost:5432/db", "postgres", "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate version failed")
}

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
			name:     "PostgreSQL with sslmode",
			driver:   "postgres",
			dbURL:    "user:pass@localhost:5432/db?sslmode=disable",
			expected: "postgres://user:pass@localhost:5432/db?sslmode=disable",
		},
		{
			name:     "Empty driver",
			driver:   "",
			dbURL:    "user:pass@localhost:5432/db",
			expected: "user:pass@localhost:5432/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.BuildDatabaseURL(tt.driver, tt.dbURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}
