// cmd/utils/content_test.go
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
}

func TestCreateContentRefreshUsed(t *testing.T) {
	cfg := &Config{
		Port:            8080,
		Host:            "localhost",
		MetrictPort:     8088,
		DBDriver:        "mysql",
		DBUser:          "root",
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
		RedisPassword:   "",
		AccessTokenTTL:  900,
		RefreshTokenTTL: 604800,
		SessionTTL:      604800,
		IdleTimeout:     1800,
		RevokedTTL:      300,
		RateLimitPeriod: 60,
		RateLimitMax:    100,
	}

	dbURL := "mysql://root:password@localhost:3306/testdb"
	content := CreateContentRefreshUsed(dbURL, cfg)

	assert.Contains(t, content, "APP_PORT=8080")
	assert.Contains(t, content, "DB_DRIVER=mysql")
	assert.Contains(t, content, "DB_URL=mysql://root:password@localhost:3306/testdb")
	assert.Contains(t, content, "REDIS_HOST=localhost")
}
