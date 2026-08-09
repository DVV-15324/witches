package utils

import (
	"os"
	"path/filepath"
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
			os.Setenv(tt.key, tt.value)
			defer os.Unsetenv(tt.key)
			result := getEnvAsInt64(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidate(t *testing.T) {
	validConfig := &Config{
		Port:            8080,
		Host:            "localhost",
		MetrictPort:     8088,
		DBHost:          "localhost",
		DBPort:          3306,
		DBName:          "testdb",
		DBPassword:      "password",
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: 60,
		ConnMaxIdleTime: 600,
		SlowThreshold:   5,
		RedisHost:       "localhost",
		RedisPort:       6379,
		AccessTokenTTL:  900,
		RefreshTokenTTL: 604800,
		SessionTTL:      604800,
		IdleTimeout:     1800,
		RevokedTTL:      300,
		RateLimitPeriod: 60,
		RateLimitMax:    100,
	}

	t.Run("valid config - should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Validate(validConfig)
		})
	})

	t.Run("missing APP_PORT - should panic", func(t *testing.T) {
		cfg := *validConfig
		cfg.Port = 0
		assert.Panics(t, func() {
			Validate(&cfg)
		})
	})

	t.Run("missing DB_HOST - should panic", func(t *testing.T) {
		cfg := *validConfig
		cfg.DBHost = ""
		assert.Panics(t, func() {
			Validate(&cfg)
		})
	})

	t.Run("missing REDIS_HOST - should panic", func(t *testing.T) {
		cfg := *validConfig
		cfg.RedisHost = ""
		assert.Panics(t, func() {
			Validate(&cfg)
		})
	})
}

func TestPreloadNotDBURL(t *testing.T) {
	// Tạo file witches.env tạm
	tmpDir := t.TempDir()
	envContent := `APP_PORT=8080
APP_HOST=localhost
METRIC_PORT=8088
DB_DRIVER=mysql
DB_USER=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=testdb
DB_PASSWORD=password
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=60
DB_CONN_MAX_IDLE_TIME=600
SLOW_THRESHOLD=5
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
ACCESS_TOKEN_TTL=900
REFRESH_TOKEN_TTL=604800
SESSION_TTL=604800
IDLE_TIMEOUT=1800
REVOKED_TTL=300
RATE_LIMIT_PERIOD=60
RATE_LIMIT_MAX=100`

	envPath := filepath.Join(tmpDir, "witches.env")
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	require.NoError(t, err)

	// Change to temp dir
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	os.Chdir(tmpDir)

	cfg := PreloadNotDBURL()
	assert.Equal(t, int64(8080), cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, int64(3306), cfg.DBPort)
	assert.Equal(t, "mysql", cfg.DBDriver)
}

func TestLoadDbUrl(t *testing.T) {
	os.Setenv("DB_URL", "mysql://root:pass@localhost:3306/testdb")
	defer os.Unsetenv("DB_URL")

	cfg := &Config{}
	cfg = LoadDbUrl(cfg)
	assert.Equal(t, "mysql://root:pass@localhost:3306/testdb", cfg.DBUrl)
}
