package cmd_migrate

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	os.Setenv("TESTCONTAINERS_REUSE_ENABLE", "true")
	code := m.Run()
	os.Exit(code)
}

// SỬA: Trả về *utils.Config thay vì dbURL string
func setupTestWithPostgres(t *testing.T) (cfg *utils.Config, migrationPath string, cleanup func()) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrate", "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	migrationsDir = filepath.ToSlash(migrationsDir)
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

	// Parse connStr để lấy thông tin gán vào Config
	u, err := url.Parse(connStr)
	require.NoError(t, err)
	port, _ := strconv.Atoi(u.Port())

	cfg = utils.DefaultConfig()
	cfg.DBDriver = "postgres"
	cfg.DBHost = u.Hostname()
	cfg.DBPort = int64(port)
	cfg.DBUser = u.User.Username()
	cfg.DBPassword, _ = u.User.Password()
	cfg.DBName = u.Path[1:] // Bỏ dấu '/' ở đầu

	return cfg, migrationsDir, cleanup
}

func TestWitchesMigrateUp(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	// SỬA: Truyền cfg thay vì dbURL
	err = WitchesMigrateUp(cfg, migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.BuildDatabaseURL(cfg))
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
	assert.True(t, tableExists, "Table 'users' should exist")
}

func TestWitchesMigrateVersion(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(cfg, migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateVersion(cfg, migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.BuildDatabaseURL(cfg))
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

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

	cfg, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(cfg, migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateDown(cfg, migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.BuildDatabaseURL(cfg))
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

func TestWitchesMigrateDrop(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateUp(cfg, migrationPath)
	require.NoError(t, err)

	err = WitchesMigrateDrop(cfg, migrationPath)
	require.NoError(t, err)

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.BuildDatabaseURL(cfg))
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

	cfg, migrationPath, cleanup := setupTestWithPostgres(t)
	defer cleanup()

	err = WitchesMigrateForce(cfg, migrationPath, "1")
	require.NoError(t, err)

	ctx := context.Background()
	db, err := sql.Open("postgres", utils.BuildDatabaseURL(cfg))
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

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

	// Tạo config giả
	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"

	err = WitchesMigrateUp(cfg, "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate up failed")
}

func TestWitchesMigrateUp_InvalidDBConfig(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}
	tmpDir := t.TempDir()

	// Config sai host để gây lỗi kết nối
	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"
	cfg.DBHost = "invalid_url"

	err = WitchesMigrateUp(cfg, tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate up failed")
}

func TestWitchesMigrateDown_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"

	err = WitchesMigrateDown(cfg, "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate down failed")
}

func TestWitchesMigrateDrop_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"

	err = WitchesMigrateDrop(cfg, "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate drop failed")
}

func TestWitchesMigrateForce_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"

	err = WitchesMigrateForce(cfg, "/invalid/path", "1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate force failed")
}

func TestWitchesMigrateVersion_InvalidPath(t *testing.T) {
	_, err := exec.LookPath("migrate")
	if err != nil {
		t.Skip("Skipping test: 'migrate' binary not found in PATH")
	}

	cfg := utils.DefaultConfig()
	cfg.DBDriver = "postgres"

	err = WitchesMigrateVersion(cfg, "/invalid/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migrate version failed")
}
