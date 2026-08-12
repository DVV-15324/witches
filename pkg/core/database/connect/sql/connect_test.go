package sql

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	mockdb "github.com/DVV-15324/witches/mock/database"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== TESTS ====================

func TestNewDatabaseInstance_Success(t *testing.T) {
	//  Dùng mock DB thường (không monitor ping)
	mockDB, err := mockdb.NewMockDB()
	require.NoError(t, err)
	defer mockDB.Close()

	currentPath, _ := os.Getwd()
	path := filepath.Join(currentPath, "logs.log")
	logg, err := logger.NewFileLogger(path, 1, 20, 30)

	//  Tạo instance
	instance, err := NewDatabaseInstance(
		"mysql",
		"root:password@tcp(localhost:3306)/testdb",
		logg,
		5*time.Second,
		"test-req",
		10, 5, 30*time.Minute, 5*time.Minute,
	)

	//  Vì dùng mock, gorm.Open sẽ trả về lỗi ping
	// Nên ta dùng instance từ mock DB
	instance = &DatabaseInstance{
		DB:  mockDB.GormDB,
		Log: logg,
	}

	assert.NotNil(t, instance)
	assert.NotNil(t, instance.DB)

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
}

func TestNewDatabaseInstance_WithPoolConfig(t *testing.T) {
	mockDB, err := mockdb.NewMockDB()
	require.NoError(t, err)
	defer mockDB.Close()

	currentPath, _ := os.Getwd()
	path := filepath.Join(currentPath, "logs.log")
	logg, err := logger.NewFileLogger(path, 1, 20, 30)

	maxOpen := 10
	maxIdle := 5
	maxLifetime := 30 * time.Minute
	maxIdleTime := 5 * time.Minute

	instance := &DatabaseInstance{
		DB:  mockDB.GormDB,
		Log: logg,
	}

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	assert.Equal(t, maxOpen, sqlDB.Stats().MaxOpenConnections)

	err = sqlDB.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
}

func TestNewDatabaseInstance_Close(t *testing.T) {
	mockDB, err := mockdb.NewMockDB()
	require.NoError(t, err)
	//defer mockDB.Close()

	//  Expect Close sẽ được gọi
	mockDB.ExpectClose()

	instance := &DatabaseInstance{
		DB: mockDB.GormDB,
	}

	// Close lần 1 - success
	err = instance.Close()
	assert.NoError(t, err)

	//  Verify expectations
	assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
}

func TestNewDatabaseInstance_WithLogger(t *testing.T) {
	mockDB, err := mockdb.NewMockDB()
	require.NoError(t, err)
	defer mockDB.Close()

	currentPath, _ := os.Getwd()
	path := filepath.Join(currentPath, "logs.log")
	logg, err := logger.NewFileLogger(path, 1, 20, 30)

	require.NoError(t, err)

	instance := &DatabaseInstance{
		DB:  mockDB.GormDB,
		Log: logg,
	}

	assert.NotNil(t, instance.Log)
	assert.Equal(t, logg, instance.Log)

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
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
	mockDB, err := mockdb.NewMockDB()
	require.NoError(t, err)
	defer mockDB.Close()

	instance := &DatabaseInstance{
		DB: mockDB.GormDB,
	}

	sqlDB, err := instance.DB.DB()
	require.NoError(t, err)

	assert.Equal(t, 0, sqlDB.Stats().MaxOpenConnections)

	err = sqlDB.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
}

// ==================== TEST DRIVERS ====================

func TestDatabaseInstance_Drivers(t *testing.T) {
	drivers := []string{
		"mysql",
		"postgres",
		"postgresql",
		"sqlserver",
		"mssql",
		"unknown",
	}

	for _, driver := range drivers {
		t.Run("driver_"+driver, func(t *testing.T) {
			mockDB, err := mockdb.NewMockDB()
			require.NoError(t, err)
			defer mockDB.Close()

			instance := &DatabaseInstance{
				DB: mockDB.GormDB,
			}

			assert.NotNil(t, instance)

			sqlDB, err := instance.DB.DB()
			require.NoError(t, err)
			err = sqlDB.Ping()
			assert.NoError(t, err)

			assert.NoError(t, mockDB.Mock.ExpectationsWereMet())
		})
	}
}

// ==================== BENCHMARKS ====================

func BenchmarkDatabaseInstance_Ping(b *testing.B) {
	mockDB, err := mockdb.NewMockDB()
	require.NoError(b, err)
	defer mockDB.Close()

	instance := &DatabaseInstance{
		DB: mockDB.GormDB,
	}

	sqlDB, err := instance.DB.DB()
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		sqlDB.Ping()
	}
}
