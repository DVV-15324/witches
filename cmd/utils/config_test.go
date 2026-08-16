package utils

import (
	"log"
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

func TestPreloadNotDBURL(t *testing.T) {
	// Tạo file .env tạm
	tmpDir := t.TempDir()
	envContent := `APP_PORT=9090
APP_HOST=127.0.0.1
METRIC_PORT=9091
UID_BITS=30
REQUEST_KEY=custom_key
LOG_PATH=/var/log
CORS_ALLOW_ORIGINS=http://example.com
CORS_ALLOW_METHODS=GET,POST
CORS_ALLOW_HEADERS=X-Custom
TIMEZONE=Asia/Ho_Chi_Minh
SUPPORTED_LANGUAGES=en,vi
JWT_SECRET=my_super_secret
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
ACCESS_TOKEN_TTL=1800
REFRESH_TOKEN_TTL=86400
SESSION_TTL=86400
IDLE_TIMEOUT=3600
REVOKED_TTL=600
RATE_LIMIT_PERIOD=120
RATE_LIMIT_MAX=200`

	envPath := filepath.Join(tmpDir, "witches.env")
	err := os.WriteFile(envPath, []byte(envContent), 0644)
	require.NoError(t, err)

	// Chuyển vào temp dir
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			log.Printf("failed to chdir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	cfg := PreloadNotDBURL()

	// Kiểm tra tất cả các field
	assert.Equal(t, int64(9090), cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, int64(9091), cfg.MetrictPort)
	assert.Equal(t, 30, cfg.UIDBits)
	assert.Equal(t, contextKey("custom_key"), contextKey(cfg.RequestKey))
	assert.Equal(t, "/var/log", cfg.LogPath)
	assert.Equal(t, "http://example.com", cfg.CorsAllowOrigins)
	assert.Equal(t, "GET,POST", cfg.CorsAllowMethods)
	assert.Equal(t, "X-Custom", cfg.CorsAllowHeaders)
	assert.Equal(t, "Asia/Ho_Chi_Minh", cfg.Timezone)
	assert.Equal(t, []string{"en", "vi"}, cfg.SupportedLanguages)
	assert.Equal(t, "my_super_secret", cfg.JWTSecret)
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
	assert.Equal(t, int64(1800), cfg.AccessTokenTTL)
	assert.Equal(t, int64(86400), cfg.RefreshTokenTTL)
	assert.Equal(t, int64(86400), cfg.SessionTTL)
	assert.Equal(t, int64(3600), cfg.IdleTimeout)
	assert.Equal(t, int64(600), cfg.RevokedTTL)
	assert.Equal(t, int64(120), cfg.RateLimitPeriod)
	assert.Equal(t, int64(200), cfg.RateLimitMax)
}

func TestLoadDbUrl(t *testing.T) {
	os.Setenv("DB_URL", "postgres://admin:securepass@db.example.com:5432/prod_db")
	defer os.Unsetenv("DB_URL")

	cfg := &Config{}
	cfg = LoadDbUrl(cfg)
	assert.Equal(t, "postgres://admin:securepass@db.example.com:5432/prod_db", cfg.DBUrl)
}
