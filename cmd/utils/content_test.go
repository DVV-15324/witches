package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContentRefresh(t *testing.T) {
	content := CreateContentRefresh()
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "APP_PORT=8080")
	assert.Contains(t, content, "DB_DRIVER=your_driver")
	assert.Contains(t, content, "REDIS_HOST=localhost")
	assert.Contains(t, content, "UID_BITS=26")
	assert.Contains(t, content, "REQUEST_KEY=request_context")
	assert.Contains(t, content, "LOG_PATH=./logs")
	assert.Contains(t, content, "CORS_ALLOW_ORIGINS=*")
	assert.Contains(t, content, "SUPPORTED_LANGUAGES=en-US,vi-VN")
}

func TestCreateContentRefreshUsed(t *testing.T) {
	cfg := &Config{
		Port:               8080,
		Host:               "localhost",
		Env:                "development",
		MetrictPort:        8088,
		UIDBits:            26,
		HashLen:            16,
		RequestKey:         "request_context",
		LogPath:            "./logs",
		CorsAllowOrigins:   "*",
		CorsAllowMethods:   "GET,POST,PUT,DELETE,OPTIONS",
		CorsAllowHeaders:   "Content-Type,Authorization",
		Timezone:           "UTC",
		SupportedLanguages: []string{"en-US", "vi-VN"},
		JWTSecret:          "your_secret_key",
		DBDriver:           "mysql",
		DBUser:             "root",
		DBHost:             "localhost",
		DBPort:             3306,
		DBName:             "testdb",
		DBPassword:         "password",
		MaxOpenConns:       100,
		MaxIdleConns:       10,
		ConnMaxLifetime:    60,
		ConnMaxIdleTime:    600,
		SlowThreshold:      5,
		RedisHost:          "localhost",
		RedisPort:          6379,
		RedisPassword:      "",
		AccessTokenTTL:     900,
		RefreshTokenTTL:    604800,
		IdleTimeout:        1800,
		RevokedTTL:         300,
		RateLimitPeriod:    60,
		RateLimitMax:       100,
	}

	dbURL := "mysql://root:password@localhost:3306/testdb"
	content := CreateContentRefreshUsed(dbURL, cfg)

	assert.Contains(t, content, "APP_PORT=8080")
	assert.Contains(t, content, "APP_HOST=localhost")
	assert.Contains(t, content, "METRIC_PORT=8088")
	assert.Contains(t, content, "DB_DRIVER=mysql")
	assert.Contains(t, content, "DB_USER=root")
	assert.Contains(t, content, "DB_URL=mysql://root:password@localhost:3306/testdb")
	assert.Contains(t, content, "REDIS_HOST=localhost")
	assert.Contains(t, content, "ACCESS_TOKEN_TTL=900")
	assert.Contains(t, content, "RATE_LIMIT_MAX=100")
}
