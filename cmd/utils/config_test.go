package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("APP_PORT", "8080")
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("METRIC_PORT", "8088")
	os.Setenv("DB_DRIVER", "mysql")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_MAX_OPEN_CONNS", "100")
	os.Setenv("DB_MAX_IDLE_CONNS", "10")
	os.Setenv("DB_CONN_MAX_LIFETIME", "60")
	os.Setenv("DB_CONN_MAX_IDLE_TIME", "600")
	os.Setenv("SLOW_THRESHOLD", "5")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("ACCESS_TOKEN_TTL", "900")
	os.Setenv("REFRESH_TOKEN_TTL", "604800")
	os.Setenv("SESSION_TTL", "604800")
	os.Setenv("IDLE_TIMEOUT", "1800")
	os.Setenv("REVOKED_TTL", "300")
	os.Setenv("RATE_LIMIT_PERIOD", "60")
	os.Setenv("RATE_LIMIT_MAX", "100")
	os.Setenv("TEST_EMPTY", "")
	os.Setenv("TEST_INVALID", "abc")
}

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

func TestPreloadNotDBURL(t *testing.T) {
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

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, int64(8080), cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, int64(8088), cfg.MetrictPort)
	assert.Equal(t, "mysql", cfg.DBDriver)
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
	assert.Equal(t, int64(900), cfg.AccessTokenTTL)
	assert.Equal(t, int64(604800), cfg.RefreshTokenTTL)
	assert.Equal(t, int64(604800), cfg.SessionTTL)
	assert.Equal(t, int64(1800), cfg.IdleTimeout)
	assert.Equal(t, int64(300), cfg.RevokedTTL)
	assert.Equal(t, int64(60), cfg.RateLimitPeriod)
	assert.Equal(t, int64(100), cfg.RateLimitMax)
}
