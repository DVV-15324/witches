package cmd_migrate

import (
	"context"
	"database/sql"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestWithPostgres(t *testing.T) (dbURL string, cleanup func()) {
	ctx := context.Background()

	// 1. Tạo thư mục tạm
	tmpDir := t.TempDir()

	// 2. Tạo thư mục migrations trong thư mục tạm
	migrationsDir := filepath.Join(tmpDir, "migrate", "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	os.Setenv("WITCHES_MIGRATIONS_PATH", migrationsDir)
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

	// 3. Chạy PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
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

	// 4. Lấy connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// 5. Trả về URL và cleanup
	cleanup = func() {
		// Dừng container sau test
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
		os.RemoveAll(tmpDir)
	}

	return connStr, cleanup
}

// Test WitchesMigrateUp - Up migration với PostgreSQL
func TestWitchesMigrateUp(t *testing.T) {
	// Kiểm tra binary migrate có sẵn
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Gọi hàm migrate up
	WitchesMigrateUp(dbURL, "postgres")

	// Kiểm tra bảng đã được tạo
	ctx := context.Background()
	connStr := dbURL + "&sslmode=disable"
	// Có thể dùng sql.Open hoặc GORM để kiểm tra
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

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

// Test WitchesMigrateVersion - Kiểm tra version
func TestWitchesMigrateVersion(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Up trước
	WitchesMigrateUp(dbURL, "postgres")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	WitchesMigrateVersion(dbURL, "postgres")

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	assert.Contains(t, string(out), "1", "Output should contain version 1")
}

// Test WitchesMigrateDown - Down migration
func TestWitchesMigrateDown(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Up trước
	WitchesMigrateUp(dbURL, "postgres")

	// Down sau
	WitchesMigrateDown(dbURL, "postgres")

	// Kiểm tra bảng đã bị xóa
	ctx := context.Background()
	connStr := dbURL + "&sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

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

// Test WitchesMigrateDrop - Drop database
func TestWitchesMigrateDrop(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Up trước
	WitchesMigrateUp(dbURL, "postgres")

	// Drop sau
	WitchesMigrateDrop(dbURL, "postgres")

	// Kiểm tra bảng đã bị xóa
	ctx := context.Background()
	connStr := dbURL + "&sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

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

// Test WitchesMigrateForce - Force version
func TestWitchesMigrateForce(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	dbURL, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// Force version 1 mà không cần migration
	WitchesMigrateForce(dbURL, "postgres", "1")

	// Kiểm tra version đã được set
	ctx := context.Background()
	connStr := dbURL + "&sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	var version int
	err = db.QueryRowContext(ctx, `
		SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1
	`).Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, 1, version, "Migration version should be forced to 1")
}
