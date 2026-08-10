package logger

import (
	"context"
	"os"
	"testing"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

func TestNewGormLogger(t *testing.T) {
	path := "./test_gorm.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	defer os.Remove(path)

	slowThreshold := 100 * time.Millisecond
	keyReq := "trace_id"
	gormLogger := NewGormLogger(logger, slowThreshold, keyReq)

	if gormLogger == nil {
		t.Error("GormLogger should not be nil")
	}
	if gormLogger.zapLogger == nil {
		t.Error("zapLogger should not be nil")
	}
	if gormLogger.SlowThreshold != slowThreshold {
		t.Errorf("SlowThreshold = %v, want %v", gormLogger.SlowThreshold, slowThreshold)
	}
	if gormLogger.KeyReq != keyReq {
		t.Errorf("KeyReq = %v, want %v", gormLogger.KeyReq, keyReq)
	}
}

func TestGormLogger_LogMode(t *testing.T) {
	path := "./test_gorm_mode.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	// Test LogMode
	newLogger := gormLogger.LogMode(gormlogger.Warn)
	if newLogger == nil {
		t.Error("LogMode should return a valid logger")
	}

	// Check that LogLevel was updated
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	// Test Info with different log levels
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
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
		threshold time.Duration
	}{
		{
			name:      "Normal query - Info level",
			logLevel:  gormlogger.Info,
			slowQuery: false,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users",
			threshold: 100 * time.Millisecond,
		},
		{
			name:      "Slow query - Warn level",
			logLevel:  gormlogger.Warn,
			slowQuery: true,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users WHERE id = 1",
			threshold: 10 * time.Millisecond,
		},
		{
			name:      "Query with error - Error level",
			logLevel:  gormlogger.Error,
			slowQuery: false,
			hasError:  true,
			rows:      -1,
			sql:       "SELECT * FROM users WHERE id = 999",
			threshold: 100 * time.Millisecond,
		},
		{
			name:      "Silent level - should not log",
			logLevel:  gormlogger.Silent,
			slowQuery: false,
			hasError:  false,
			rows:      10,
			sql:       "SELECT * FROM users",
			threshold: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormLogger.LogLevel = tt.logLevel
			gormLogger.SlowThreshold = tt.threshold

			var err error
			if tt.hasError {
				err = context.DeadlineExceeded
			}

			// Sleep to simulate slow query if needed
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	// Test with context that has trace_id
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")
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
	defer logger.Sync()
	defer os.Remove(path)

	gormLogger := NewGormLogger(logger, 100*time.Millisecond, "trace_id")
	if gormLogger == nil {
		t.Fatal("GormLogger should not be nil")
	}

	ctx := context.Background()

	// Test all methods with different log levels
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
