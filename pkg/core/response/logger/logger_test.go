package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLogger_File(t *testing.T) {
	path := "./test.log"
	logger, err := NewFileLogger(path, 1, 20, 30)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	logger.Info("Hello, I'am Info")
	logger.Warn("Hello, I'am Warn")
	logger.Error("Hello, I'am Error")

	logger.InfoWithFields("Info", zap.String("Field", "hello"))
	logger.WarnWithFields("Info", zap.String("Field", "hello"))
	logger.ErrorWithFields("Info", zap.String("Field", "hello"))
}

// Test các config khác nhau
func TestNewFileLogger_WithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		maxSize   int
		maxBackup int
		maxAge    int
	}{
		{
			name:      "valid config with small size",
			path:      "./test1.log",
			maxSize:   1,
			maxBackup: 5,
			maxAge:    10,
		},
		{
			name:      "valid config with large size",
			path:      "./test2.log",
			maxSize:   100,
			maxBackup: 10,
			maxAge:    30,
		},
		{
			name:      "zero max size",
			path:      "./test3.log",
			maxSize:   0,
			maxBackup: 3,
			maxAge:    7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewFileLogger(tt.path, tt.maxSize, tt.maxBackup, tt.maxAge)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer func() { _ = logger.Sync() }()
			defer func() { _ = os.Remove(tt.path) }()

			if logger == nil {
				t.Error("Logger should not be nil")
			}
			if logger.Log == nil {
				t.Error("Logger.Log should not be nil")
			}
		})
	}
}

// Test tất cả methods
func TestLoggerMethods(t *testing.T) {
	path := "./test_methods.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	t.Run("Info", func(t *testing.T) {
		logger.Info("info message")
	})

	t.Run("Warn", func(t *testing.T) {
		logger.Warn("warn message")
	})

	t.Run("Error", func(t *testing.T) {
		logger.Error("error message")
	})

	t.Run("InfoWithFields", func(t *testing.T) {
		logger.InfoWithFields("info with fields",
			zap.String("key1", "value1"),
			zap.Int("key2", 123),
		)
	})

	t.Run("WarnWithFields", func(t *testing.T) {
		logger.WarnWithFields("warn with fields",
			zap.String("key", "value"),
			zap.Bool("flag", true),
		)
	})

	t.Run("ErrorWithFields", func(t *testing.T) {
		logger.ErrorWithFields("error with fields",
			zap.String("error", "something went wrong"),
			zap.Int("code", 500),
		)
	})
}

// Test Sync
func TestLoggerSync(t *testing.T) {
	path := "./test_sync.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	logger.Info("test message")

	err = logger.Sync()
	if err != nil {
		t.Logf("Sync returned: %v", err)
	}
}

// Test nhiều messages
func TestLogger_MultipleMessages(t *testing.T) {
	path := "./test_multiple.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	for i := 0; i < 10; i++ {
		logger.InfoWithFields("message", zap.Int("index", i))
		logger.WarnWithFields("warning", zap.Int("index", i))
		logger.ErrorWithFields("error", zap.Int("index", i))
	}
}

// Test với các loại field khác nhau
func TestLogger_WithDifferentFieldTypes(t *testing.T) {
	path := "./test_fields.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	fields := []zap.Field{
		zap.String("string", "value"),
		zap.Int("int", 123),
		zap.Int64("int64", 1234567890),
		zap.Float64("float", 3.14),
		zap.Bool("bool", true),
		zap.Strings("strings", []string{"a", "b", "c"}),
		zap.Ints("ints", []int{1, 2, 3}),
	}

	logger.InfoWithFields("all field types", fields...)
}

// Test concurrent
func TestLogger_Concurrent(t *testing.T) {
	path := "./test_concurrent.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				logger.InfoWithFields("concurrent",
					zap.Int("goroutine", id),
					zap.Int("iteration", j),
				)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
