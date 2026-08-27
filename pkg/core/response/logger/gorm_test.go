package logger

import (
	"context"
	"os"
	"testing"
	"time"

	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	gormlogger "gorm.io/gorm/logger"
)

func TestNewGormLogger(t *testing.T) {
	path := "./test_gorm.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	defer func() {
		_ = os.Remove(path)
	}()

	config := wcmd_utils.DefaultConfig()
	config.SlowThreshold = 5
	gormLogger := NewGormLogger(logger, config)

	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}
	if gormLogger.zapLogger == nil {
		t.Fatal("zapLogger should not be nil")
	}
	if gormLogger.config == nil {
		t.Fatal("config should not be nil")
	}
	if gormLogger.config.RequestKey != "request_context" {
		t.Fatalf("RequestKey = %v, want %v", gormLogger.config.RequestKey, "request_context")
	}
	if gormLogger.config.SlowThreshold != 5 {
		t.Fatalf("SlowThreshold = %v, want %v", gormLogger.config.SlowThreshold, 5)
	}
}

func TestGormLogger_LogMode(t *testing.T) {
	path := "./test_gorm_mode.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	defer func() {
		_ = os.Remove(path)
	}()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	newLogger := gormLogger.LogMode(gormlogger.Warn)
	if newLogger == nil {
		t.Error("LogMode should return a valid logger")
	}
	if newLogger.(*GormZapLogger).LogLevel != gormlogger.Warn {
		t.Errorf("LogLevel = %v, want %v", newLogger.(*GormZapLogger).LogLevel, gormlogger.Warn)
	}
}

func TestGormLogger_Info(t *testing.T) {
	path := "./test_gorm_info.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	defer func() {
		_ = os.Remove(path)
	}()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		logLevel gormlogger.LogLevel
		msg      string
		data     []interface{}
	}{
		{
			name:     "Info with Info level",
			logLevel: gormlogger.Info,
			msg:      "test info message",
			data:     []interface{}{"key", "value"},
		},
		{
			name:     "Info with Warn level (should not log)",
			logLevel: gormlogger.Warn,
			msg:      "test info message with warn level",
			data:     []interface{}{"key", "value"},
		},
		{
			name:     "Info with Silent level (should not log)",
			logLevel: gormlogger.Silent,
			msg:      "test info message with silent level",
			data:     []interface{}{"key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormLogger.LogLevel = tt.logLevel
			gormLogger.Info(ctx, tt.msg, tt.data...)
		})
	}
}

func TestGormLogger_Warn(t *testing.T) {
	path := "./test_gorm_warn.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	defer func() { _ = os.Remove(path) }()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		logLevel gormlogger.LogLevel
		msg      string
		data     []interface{}
	}{
		{
			name:     "Warn with Warn level",
			logLevel: gormlogger.Warn,
			msg:      "test warn message",
			data:     []interface{}{"key", "value"},
		},
		{
			name:     "Warn with Error level (should not log)",
			logLevel: gormlogger.Error,
			msg:      "test warn message with error level",
			data:     []interface{}{"key", "value"},
		},
		{
			name:     "Warn with Silent level (should not log)",
			logLevel: gormlogger.Silent,
			msg:      "test warn message with silent level",
			data:     []interface{}{"key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormLogger.LogLevel = tt.logLevel
			gormLogger.Warn(ctx, tt.msg, tt.data...)
		})
	}
}

func TestGormLogger_Error(t *testing.T) {
	path := "./test_gorm_error.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		logLevel gormlogger.LogLevel
		msg      string
		data     []interface{}
	}{
		{
			name:     "Error with Error level",
			logLevel: gormlogger.Error,
			msg:      "test error message",
			data:     []interface{}{"key", "value"},
		},
		{
			name:     "Error with Silent level (should not log)",
			logLevel: gormlogger.Silent,
			msg:      "test error message with silent level",
			data:     []interface{}{"key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormLogger.LogLevel = tt.logLevel
			gormLogger.Error(ctx, tt.msg, tt.data...)
		})
	}
}

func TestGormLogger_Trace(t *testing.T) {
	path := "./test_gorm_trace.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	config := wcmd_utils.DefaultConfig()
	config.SlowThreshold = 100
	config.SlowThreshold = 0
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()
	begin := time.Now()

	tests := []struct {
		name      string
		logLevel  gormlogger.LogLevel
		slowQuery bool
		hasError  bool
		rows      int64
		sql       string
	}{
		{
			name:      "Normal query - Info level",
			logLevel:  gormlogger.Info,
			slowQuery: false,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users",
		},
		{
			name:      "Slow query - Warn level",
			logLevel:  gormlogger.Warn,
			slowQuery: true,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users WHERE id = 1",
		},
		{
			name:      "Query with error - Error level",
			logLevel:  gormlogger.Error,
			slowQuery: false,
			hasError:  true,
			rows:      -1,
			sql:       "SELECT * FROM users WHERE id = 999",
		},
		{
			name:      "Silent level - should not log",
			logLevel:  gormlogger.Silent,
			slowQuery: false,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormLogger.LogLevel = tt.logLevel

			if tt.slowQuery {
				config.SlowThreshold = 1
				config.SlowThreshold = 0
			} else {
				config.SlowThreshold = 100
			}

			var err error
			if tt.hasError {
				err = context.DeadlineExceeded
			}

			if tt.slowQuery {
				time.Sleep(20 * time.Millisecond)
			}

			gormLogger.Trace(ctx, begin, func() (sql string, rowsAffected int64) {
				return tt.sql, tt.rows
			}, err)
		})
	}
}

func TestGormLogger_TraceWithContext(t *testing.T) {
	path := "./test_gorm_context.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.WithValue(context.Background(), config.RequestKey, "test-trace-123")
	begin := time.Now()

	gormLogger.Trace(ctx, begin, func() (sql string, rowsAffected int64) {
		return "SELECT * FROM users WHERE id = 1", 1
	}, nil)
}

func TestGormLogger_AllLevelsCombined(t *testing.T) {
	path := "./test_gorm_all.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	config := wcmd_utils.DefaultConfig()
	gormLogger := NewGormLogger(logger, config)
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	levels := []gormlogger.LogLevel{
		gormlogger.Silent,
		gormlogger.Error,
		gormlogger.Warn,
		gormlogger.Info,
	}

	for i, level := range levels {
		t.Run("Level_"+string(rune('0'+i)), func(t *testing.T) {
			gormLogger.LogLevel = level

			gormLogger.Info(ctx, "info message", "key", "value")
			gormLogger.Warn(ctx, "warn message", "key", "value")
			gormLogger.Error(ctx, "error message", "key", "value")

			begin := time.Now()
			gormLogger.Trace(ctx, begin, func() (sql string, rowsAffected int64) {
				return "SELECT * FROM users", 10
			}, nil)
		})
	}
}
