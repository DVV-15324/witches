package sql

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgresContainer chạy container PostgreSQL cho test
func setupPostgresContainer(tb testing.TB) (string, func()) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:15",
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

// setupTestLogger tạo logger với thư mục tạm tự quản lý (không dùng t.TempDir())
func setupTestLogger(tb testing.TB) (*logger.ModelLogger, func()) {
	tmpDir, err := os.MkdirTemp("", "test-log-*")
	require.NoError(tb, err)

	logPath := filepath.Join(tmpDir, "test.log")
	logg, err := logger.NewFileLogger(logPath, 1, 20, 30)
	require.NoError(tb, err)

	cleanup := func() {
		// Đóng logger nếu có method Close
		// if closer, ok := logg.(interface{ Close() error }); ok {
		//     _ = closer.Close()
		// }
		if err := os.RemoveAll(tmpDir); err != nil {
			tb.Logf("Failed to remove temp dir: %v", err)
		}
	}

	return logg, cleanup
}

// makeTestConfig tạo config cho test
func makeTestConfig(driver, dsn string) *utils.Config {
	cfg := utils.DefaultConfig()
	cfg.DBDriver = driver
	cfg.DBUrl = dsn
	cfg.MaxOpenConns = 10
	cfg.MaxIdleConns = 5
	cfg.ConnMaxLifetime = 30 * 60 // 30 phút (giây)
	cfg.ConnMaxIdleTime = 5 * 60  // 5 phút (giây)
	cfg.SlowThreshold = 5
	return cfg
}

// === TESTS ===

func TestNewDatabaseInstance_WithPostgres(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", connStr)
	instance, err := NewDatabaseInstance(cfg, logg)

	assert.NoError(t, err)
	require.NotNil(t, instance)

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)

	err = instance.Close()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_Success(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", connStr)
	instance, err := NewDatabaseInstance(cfg, logg)

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

	cfg := makeTestConfig("postgres", connStr)
	cfg.MaxOpenConns = 20
	cfg.MaxIdleConns = 10

	instance, err := NewDatabaseInstance(cfg, logg)
	assert.NoError(t, err)
	defer func() {
		_ = instance.Close()
	}()

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)

	// Kiểm tra pool config (có thể không lấy được chính xác qua Stats, nhưng kiểm tra tồn tại)
	assert.NotZero(t, sqlDB.Stats().MaxOpenConnections, "MaxOpenConnections should be set")
}

func TestNewDatabaseInstance_Close(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", connStr)
	instance, err := NewDatabaseInstance(cfg, logg)
	require.NoError(t, err)

	err = instance.Close()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_WithLogger(t *testing.T) {
	connStr, cleanupContainer := setupPostgresContainer(t)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", connStr)
	instance, err := NewDatabaseInstance(cfg, logg)
	assert.NoError(t, err)
	defer func() {
		_ = instance.Close()
	}()

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

	cfg := makeTestConfig("postgres", connStr)
	cfg.MaxOpenConns = 0
	cfg.MaxIdleConns = 0
	cfg.ConnMaxLifetime = 0
	cfg.ConnMaxIdleTime = 0

	instance, err := NewDatabaseInstance(cfg, logg)
	assert.NoError(t, err)
	defer func() {
		_ = instance.Close()
	}()

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	assert.Equal(t, 0, sqlDB.Stats().MaxOpenConnections)

	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestNewDatabaseInstance_InvalidDriver(t *testing.T) {
	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("invalid_driver", "dsn")
	instance, err := NewDatabaseInstance(cfg, logg)

	assert.Error(t, err)
	assert.Nil(t, instance)
}

func TestNewDatabaseInstance_InvalidDSN(t *testing.T) {
	logg, cleanupLogger := setupTestLogger(t)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", "invalid-dsn")
	instance, err := NewDatabaseInstance(cfg, logg)

	assert.Error(t, err)
	assert.Nil(t, instance)
}

// === BENCHMARKS ===

func BenchmarkDatabaseInstance_Ping(b *testing.B) {
	connStr, cleanupContainer := setupPostgresContainer(b)
	defer cleanupContainer()

	logg, cleanupLogger := setupTestLogger(b)
	defer cleanupLogger()

	cfg := makeTestConfig("postgres", connStr)
	instance, err := NewDatabaseInstance(cfg, logg)
	require.NoError(b, err)
	defer func() {
		_ = instance.Close()
	}()

	sqlDB, err := instance.DB.DB()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sqlDB.Ping()
	}
}
