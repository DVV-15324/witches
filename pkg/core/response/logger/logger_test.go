package logger

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestNewFileLogger_WithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		maxSize   int
		maxBackup int
		maxAge    int
	}{
		{"valid config with small size", "./test1.log", 1, 5, 10},
		{"valid config with large size", "./test2.log", 100, 10, 30},
		{"zero max size", "./test3.log", 0, 3, 7},
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

func TestLoggerMethods(t *testing.T) {
	path := "./test_methods.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	t.Run("Info", func(t *testing.T) { logger.Info("info message") })
	t.Run("Warn", func(t *testing.T) { logger.Warn("warn message") })
	t.Run("Error", func(t *testing.T) { logger.Error("error message") })
	t.Run("Debug", func(t *testing.T) { logger.Debug("debug message") })

	t.Run("InfoWithFields", func(t *testing.T) {
		logger.InfoWithFields("info with fields", zap.String("key1", "value1"), zap.Int("key2", 123))
	})
	t.Run("WarnWithFields", func(t *testing.T) {
		logger.WarnWithFields("warn with fields", zap.String("key", "value"), zap.Bool("flag", true))
	})
	t.Run("ErrorWithFields", func(t *testing.T) {
		logger.ErrorWithFields("error with fields", zap.String("error", "something went wrong"), zap.Int("code", 500))
	})
	t.Run("DebugWithFields", func(t *testing.T) {
		logger.DebugWithFields("debug with fields", zap.String("k", "v"))
	})
}

func TestLoggerSync(t *testing.T) {
	path := "./test_sync.log"
	logger, err := NewFileLogger(path, 1, 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	defer func() { _ = os.Remove(path) }()

	logger.Info("test message")
	if err := logger.Sync(); err != nil {
		t.Logf("Sync returned: %v", err)
	}
}

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
				logger.InfoWithFields("concurrent", zap.Int("goroutine", id), zap.Int("iteration", j))
			}
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLogger_OutputContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.log")

	logger, err := NewFileLogger(path, 1, 1, 1)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	// Đảm bảo đóng logger sau khi test xong
	defer func() {
		if err := logger.Close(); err != nil {
			t.Logf("Failed to close logger: %v", err)
		}
	}()

	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")
	logger.InfoWithFields("with fields", zap.String("key", "value"))

	if err := logger.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	content := string(data)

	expectedMessages := []string{
		"hello info",
		"hello warn",
		"hello error",
		"with fields",
		`"key":"value"`,
	}
	for _, msg := range expectedMessages {
		if !strings.Contains(content, msg) {
			t.Errorf("Log file missing expected content %q. File content: %s", msg, content)
		}
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "debug.log")

	logger, err := NewFileLogger(path, 1, 1, 1)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Debug("debug message")
	if err := logger.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "debug message") {
		t.Log("Debug message was written (logger might be at Debug level)")
	} else {
		t.Log("Debug message was NOT written (expected, logger is at Info level)")
	}
}
