package sql

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DVV-15324/witches/pkg/core/response/logger"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ==================== HELPERS ====================

func setupPostgresContainer(tb testing.TB) (string, func()) {
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
	require.NoError(tb, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(tb, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			tb.Logf("Failed to terminate container: %v", err)
		}
	}

	return connStr, cleanup
}

// 👇 KHÔNG dùng t.TempDir() - tự quản lý thư mục
func setupTestLogger(tb testing.TB) (*logger.EntityLogger, func()) {
	// Tạo thư mục tạm thủ công
	tmpDir, err := os.MkdirTemp("", "test-log-*")
	require.NoError(tb, err)

	logPath := filepath.Join(tmpDir, "test.log")
	logg, err := logger.NewFileLogger(logPath, 1, 20, 30)
	require.NoError(tb, err)

	cleanup := func() {
		// Đóng logger trước (nếu có method)
		// if closer, ok := logg.(interface{ Close() error }); ok {
		//     closer.Close()
		// }

		// Xóa thư mục sau khi test
		if err := os.RemoveAll(tmpDir); err != nil {
			tb.Logf("Failed to remove temp dir: %v", err)
		}
	}

	return logg, cleanup
}

// ==================== TESTS ====================

func TestDatabaseInstance_WithPostgres(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger() // 👈 ĐẢM BẢO cleanup sau test

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	assert.NoError(t, err)
	defer instance.Close()

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_Success(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)

	assert.NoError(t, err)
	assert.NotNil(t, instance)

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)

	err = instance.Close()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_WithPoolConfig(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	maxOpen := 10
	maxIdle := 5
	maxLifetime := 30 * time.Minute
	maxIdleTime := 5 * time.Minute

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		int64(maxOpen), int64(maxIdle), maxLifetime, maxIdleTime,
	)
	assert.NoError(t, err)
	defer instance.Close()

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	assert.Equal(t, maxOpen, sqlDB.Stats().MaxOpenConnections)

	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_Close(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	require.NoError(t, err)

	err = instance.Close()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_WithLogger(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	assert.NoError(t, err)
	defer instance.Close()

	assert.NotNil(t, instance.Log)
	assert.Equal(t, logg, instance.Log)

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestDatabaseInstance_Close_Error(t *testing.T) {
	instance := &DatabaseInstance{
		DB: nil,
	}

	err := instance.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database instance is nil")
}

func TestNewDatabaseInstance_DefaultConfig(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		0, 0, 0, 0,
	)
	assert.NoError(t, err)
	defer instance.Close()

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	assert.Equal(t, 0, sqlDB.Stats().MaxOpenConnections)

	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestDatabaseInstance_InvalidDriver(t *testing.T) {
	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"invalid_driver",
		"dsn",
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	assert.Error(t, err)
	assert.Nil(t, instance)
}

func TestDatabaseInstance_InvalidDSN(t *testing.T) {
	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		"invalid-dsn",
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	assert.Error(t, err)
	assert.Nil(t, instance)
}

// ==================== BENCHMARKS ====================

func BenchmarkDatabaseInstance_Ping(b *testing.B) {
	connStr, cleanupContainer := setupPostgresContainer(b)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(b)
	defer cleanupLogger()

	instance, err := NewDatabaseInstance(
		"postgres",
		connStr,
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)
	require.NoError(b, err)
	defer instance.Close()

	sqlDB, err := instance.DB.DB()
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		sqlDB.Ping()
	}
}
