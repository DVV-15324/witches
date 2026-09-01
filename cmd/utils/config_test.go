package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvAsInt64(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected int64
	}{
		{"valid int", "TEST_INT", "123", 123},
		{"zero", "TEST_ZERO", "0", 0},
		{"negative", "TEST_NEG", "-5", -5},
		{"empty string", "TEST_EMPTY", "", 0},
		{"invalid string", "TEST_INVALID", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv(tt.key, tt.value)
			defer func() {
				_ = os.Unsetenv(tt.key)
			}()
			result := getEnvAsInt64(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPreloadNotDBURL(t *testing.T) {
	tmpDir := t.TempDir()
	envContent := `APP_PORT=9090
APP_HOST=127.0.0.1
APP_ENV=development
METRIC_PORT=9091
UID_BITS=26
REQUEST_KEY=custom_key
LOG_PATH=/var/log
CORS_ALLOW_ORIGINS=http://example.com
CORS_ALLOW_METHODS=GET,POST
CORS_ALLOW_HEADERS=X-Custom
TIMEZONE=Asia/Ho_Chi_Minh
SUPPORTED_LANGUAGES=en,vi
DB_DRIVER=postgres
DB_USER=admin
DB_HOST=db.example.com
DB_PORT=5432
DB_NAME=prod_db
DB_PASSWORD=securepass
DB_MAX_OPEN_CONNS=200
DB_MAX_IDLE_CONNS=20
DB_CONN_MAX_LIFETIME=120
DB_CONN_MAX_IDLE_TIME=900
SLOW_THRESHOLD=10
REDIS_HOST=cache.example.com
REDIS_PORT=6380
REDIS_PASSWORD=redispass
RATE_LIMIT_PERIOD=120
RATE_LIMIT_MAX=200
CUSTOM_KEY=custom_value`

	envPath := filepath.Join(tmpDir, "witches.env")
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	require.NoError(t, err)

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	cfg := PreloadNotDBURL()

	assert.Equal(t, int64(9090), cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, int64(9091), cfg.MetrictPort)
	assert.Equal(t, int64(26), cfg.UIDBits)
	assert.Equal(t, contextKey("custom_key"), cfg.RequestKey)
	assert.Equal(t, "/var/log", cfg.LogPath)
	assert.Equal(t, "http://example.com", cfg.CorsAllowOrigins)
	assert.Equal(t, "GET,POST", cfg.CorsAllowMethods)
	assert.Equal(t, "X-Custom", cfg.CorsAllowHeaders)
	assert.Equal(t, "Asia/Ho_Chi_Minh", cfg.Timezone)
	assert.Equal(t, []string{"en", "vi"}, cfg.SupportedLanguages)
	assert.Equal(t, "postgres", cfg.DBDriver)
	assert.Equal(t, "admin", cfg.DBUser)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, int64(5432), cfg.DBPort)
	assert.Equal(t, "prod_db", cfg.DBName)
	assert.Equal(t, "securepass", cfg.DBPassword)
	assert.Equal(t, int64(200), cfg.MaxOpenConns)
	assert.Equal(t, int64(20), cfg.MaxIdleConns)
	assert.Equal(t, int64(120), cfg.ConnMaxLifetime)
	assert.Equal(t, int64(900), cfg.ConnMaxIdleTime)
	assert.Equal(t, int64(10), cfg.SlowThreshold)
	assert.Equal(t, "cache.example.com", cfg.RedisHost)
	assert.Equal(t, int64(6380), cfg.RedisPort)
	assert.Equal(t, "redispass", cfg.RedisPassword)
	assert.Equal(t, int64(120), cfg.RateLimitPeriod)
	assert.Equal(t, int64(200), cfg.RateLimitMax)
	assert.Equal(t, "custom_value", cfg.Others["CUSTOM_KEY"])
}

func TestLoadDbUrl(t *testing.T) {
	_ = os.Setenv("DB_URL", "postgres://admin:securepass@db.example.com:5432/prod_db")
	defer func() {
		_ = os.Unsetenv("DB_URL")
	}()

	cfg := &Config{}
	LoadDbUrl(cfg)
	assert.Equal(t, "postgres://admin:securepass@db.example.com:5432/prod_db", cfg.DBUrl)
}

func TestLoadModule(t *testing.T) {
	tmpDir := t.TempDir()
	modulePath := filepath.Join(tmpDir, "internal", "book")
	err := os.MkdirAll(modulePath, 0755)
	require.NoError(t, err)

	envContent := `MODULE_KEY=module_value
ANOTHER_KEY=another_value`
	err = os.WriteFile(filepath.Join(modulePath, "module.env"), []byte(envContent), 0644)
	require.NoError(t, err)

	cfg := &Config{}
	LoadModule(modulePath, cfg)

	assert.NotNil(t, cfg.Domains)
	assert.Equal(t, "module_value", cfg.Domains[modulePath]["MODULE_KEY"])
	assert.Equal(t, "another_value", cfg.Domains[modulePath]["ANOTHER_KEY"])
}

func TestLoadModule_NoEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	modulePath := filepath.Join(tmpDir, "internal", "book")
	err := os.MkdirAll(modulePath, 0755)
	require.NoError(t, err)

	cfg := &Config{}
	LoadModule(modulePath, cfg)

	// Không panic, Domains vẫn nil
	assert.Nil(t, cfg.Domains)
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Tạo witches.env
	envContent := `APP_PORT=9090
DB_URL=postgres://user:pass@localhost:5432/db`
	err = os.WriteFile("witches.env", []byte(envContent), 0644)
	require.NoError(t, err)

	// Tạo internal folder và module
	internalDir := filepath.Join(tmpDir, "internal")
	err = os.MkdirAll(internalDir, 0755)
	require.NoError(t, err)

	bookDir := filepath.Join(internalDir, "book")
	err = os.MkdirAll(bookDir, 0755)
	require.NoError(t, err)

	moduleEnv := `BOOK_KEY=book_value`
	err = os.WriteFile(filepath.Join(bookDir, "module.env"), []byte(moduleEnv), 0644)
	require.NoError(t, err)

	cfg := Load()

	assert.Equal(t, int64(9090), cfg.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DBUrl)
	assert.NotNil(t, cfg.Domains)

	// Debug: in ra tất cả keys
	t.Logf("Domains keys: %v", getKeys(cfg.Domains))

	// Tìm key chứa "book" trong path
	var foundKey string
	for k := range cfg.Domains {
		if strings.Contains(k, "book") {
			foundKey = k
			break
		}
	}

	if foundKey != "" {
		assert.Equal(t, "book_value", cfg.Domains[foundKey]["BOOK_KEY"])
	} else {
		// Nếu không tìm thấy, dùng bookDir làm key
		assert.Equal(t, "book_value", cfg.Domains[bookDir]["BOOK_KEY"])
	}
}

// Helper function
func getKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, int64(8080), cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, int64(8088), cfg.MetrictPort)
	assert.Equal(t, int64(26), cfg.UIDBits)
	assert.Equal(t, contextKey("request_context"), cfg.RequestKey)
	assert.Equal(t, "./logs", cfg.LogPath)
	assert.Equal(t, "*", cfg.CorsAllowOrigins)
	assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", cfg.CorsAllowMethods)
	assert.Equal(t, "Content-Type,Authorization", cfg.CorsAllowHeaders)
	assert.Equal(t, "UTC", cfg.Timezone)
	assert.Equal(t, []string{"en-US", "vi-VN"}, cfg.SupportedLanguages)
	assert.Equal(t, "mysql", cfg.DBDriver)
	assert.Equal(t, "root", cfg.DBUser)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, int64(3306), cfg.DBPort)
	assert.Equal(t, "your_database", cfg.DBName)
	assert.Equal(t, "your_password", cfg.DBPassword)
	assert.Equal(t, int64(100), cfg.MaxOpenConns)
	assert.Equal(t, int64(10), cfg.MaxIdleConns)
	assert.Equal(t, int64(60), cfg.ConnMaxLifetime)
	assert.Equal(t, int64(600), cfg.ConnMaxIdleTime)
	assert.Equal(t, int64(5), cfg.SlowThreshold)
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, int64(6379), cfg.RedisPort)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, int64(60), cfg.RateLimitPeriod)
	assert.Equal(t, int64(100), cfg.RateLimitMax)
}

func TestIsKnownKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"APP_PORT", "APP_PORT", true},
		{"APP_HOST", "APP_HOST", true},
		{"DB_DRIVER", "DB_DRIVER", true},
		{"UNKNOWN_KEY", "UNKNOWN_KEY", false},
		{"CUSTOM_KEY", "CUSTOM_KEY", false},
		{"REDIS_HOST", "REDIS_HOST", true},
		{"RATE_LIMIT_MAX", "RATE_LIMIT_MAX", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKnownKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		expected []string
	}{
		{"normal", "a,b,c", ",", []string{"a", "b", "c"}},
		{"with spaces", "a, b, c", ",", []string{"a", "b", "c"}},
		{"empty", "", ",", []string{}},
		{"single", "hello", ",", []string{"hello"}},
		{"multiple sep", "a;b;c", ";", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitAndTrim(tt.s, tt.sep)
			assert.Equal(t, tt.expected, result)
		})
	}
}
