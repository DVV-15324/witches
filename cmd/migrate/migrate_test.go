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

func setupTestWithPostgres(t *testing.T) (dbURL string, migrationPath string, cleanup func()) {
	ctx := context.Background()

	// 1. Tạo thư mục tạm
	tmpDir := t.TempDir()

	// 2. Tạo thư mục migrations trong thư mục tạm
	migrationsDir := filepath.Join(tmpDir, "migrate", "migrations")
	migrationsDir = filepath.ToSlash(migrationsDir)
	err := os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err)

	// 3. Tạo migration files...
	upFile := filepath.Join(migrationsDir, "000001_test_migration.up.sql")
	downFile := filepath.Join(migrationsDir, "000001_test_migration.down.sql")

	upContent := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	downContent := `
DROP TABLE IF EXISTS users;
`

	err = os.WriteFile(upFile, []byte(upContent), 0644)
	require.NoError(t, err)
	err = os.WriteFile(downFile, []byte(downContent), 0644)
	require.NoError(t, err)

	// 4. Chạy PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15",
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

	// 5. Lấy connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// 6. Trả về URL và cleanup
	cleanup = func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
		_ = os.RemoveAll(tmpDir)
	}
	connStr = strings.TrimPrefix(connStr, "postgres://")
	return connStr, migrationsDir, cleanup
}

// Test WitchesMigrateUp
func TestWitchesMigrateUp(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()
	WitchesMigrateUp(dbURL, "postgres", migrationPath)
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

// Test WitchesMigrateVersion
func TestWitchesMigrateVersion(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Up trước
	WitchesMigrateUp(dbURL, "postgres", migrationPath)

	// Gọi hàm version
	WitchesMigrateVersion(dbURL, "postgres", migrationPath)

	// Kiểm tra version trong DB
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

// Test WitchesMigrateDown
func TestWitchesMigrateDown(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	WitchesMigrateUp(dbURL, "postgres", migrationPath)
	WitchesMigrateDown(dbURL, "postgres", migrationPath)

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

// Test WitchesMigrateDrop
func TestWitchesMigrateDrop(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	WitchesMigrateUp(dbURL, "postgres", migrationPath)
	WitchesMigrateDrop(dbURL, "postgres", migrationPath)

	ctx := context.Background()
	connStr := utils.BuildDatabaseURL("postgres", dbURL)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	func() { _ = db.Close() }()

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

// Test WitchesMigrateForce
func TestWitchesMigrateForce(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Force version 1 mà không cần migration thực tế
	WitchesMigrateForce(dbURL, "postgres", migrationPath, "1")

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
