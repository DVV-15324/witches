package utils

import (
	"fmt"
	"strings"
)

func CreateContentRefresh() string {
	return `# SERVER CONFIGURATION
APP_PORT=8080
APP_HOST=localhost
APP_ENV=development
METRIC_PORT=8088
UID_BITS=26
HASH_LEN=16
REQUEST_KEY=request_context

# LOGGING
LOG_PATH=./logs

# CORS 
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization

# TIMEZONE 
TIMEZONE=UTC
SUPPORTED_LANGUAGES=en-US,vi-VN

# SECURITY 
JWT_SECRET=your_secret_key

# DATABASE CONFIGURATION
DB_DRIVER=your_driver
DB_USER=your_user
DB_HOST=localhost
DB_PORT=3306
DB_NAME=your_database
DB_PASSWORD=your_password

# Database Connection Pool
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=60
DB_CONN_MAX_IDLE_TIME=600
SLOW_THRESHOLD=5

# REDIS CONFIGURATION
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# TOKEN EXPIRATION
ACCESS_TOKEN_TTL=3600
REFRESH_TOKEN_TTL=604800

# SESSION SETTINGS
IDLE_TIMEOUT=3600

# BLACKLIST SETTINGS
REVOKED_TTL=300

# RATE LIMIT
RATE_LIMIT_PERIOD=60
RATE_LIMIT_MAX=100 
`
}

func CreateContentRefreshUsed(DB_URL string, cfg *Config) string {
	return fmt.Sprintf(`# SERVER CONFIGURATION
APP_PORT=%d
APP_HOST=%s
APP_ENV=%s
METRIC_PORT=%d
UID_BITS=%d
HASH_LEN=%d
REQUEST_KEY=%s

# LOGGING
LOG_PATH=%s

# CORS
CORS_ALLOW_ORIGINS=%s
CORS_ALLOW_METHODS=%s
CORS_ALLOW_HEADERS=%s

# TIMEZONE
TIMEZONE=%s
SUPPORTED_LANGUAGES=%s

# SECURITY
JWT_SECRET=%s

# DATABASE CONFIGURATION
DB_DRIVER=%s
DB_USER=%s
DB_HOST=%s
DB_PORT=%d
DB_NAME=%s
DB_PASSWORD=%s
DB_URL=%s

# Database Connection Pool
DB_MAX_OPEN_CONNS=%d
DB_MAX_IDLE_CONNS=%d
DB_CONN_MAX_LIFETIME=%d
DB_CONN_MAX_IDLE_TIME=%d
SLOW_THRESHOLD=%d

# REDIS CONFIGURATION
REDIS_HOST=%s
REDIS_PORT=%d
REDIS_PASSWORD=%s

# TOKEN EXPIRATION
ACCESS_TOKEN_TTL=%d
REFRESH_TOKEN_TTL=%d

# SESSION SETTINGS
IDLE_TIMEOUT=%d

# BLACKLIST SETTINGS
REVOKED_TTL=%d

# RATE LIMIT
RATE_LIMIT_PERIOD=%d
RATE_LIMIT_MAX=%d
`,
		cfg.Port, cfg.Host, cfg.Env, cfg.MetrictPort,
		cfg.UIDBits, cfg.HashLen, cfg.RequestKey,
		cfg.LogPath,
		cfg.CorsAllowOrigins, cfg.CorsAllowMethods, cfg.CorsAllowHeaders,
		cfg.Timezone, strings.Join(cfg.SupportedLanguages, ","),
		cfg.JWTSecret,
		cfg.DBDriver, cfg.DBUser, cfg.DBHost,
		cfg.DBPort, cfg.DBName, cfg.DBPassword, DB_URL,
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime,
		cfg.ConnMaxIdleTime, cfg.SlowThreshold,
		cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword,
		cfg.AccessTokenTTL, cfg.RefreshTokenTTL,
		cfg.IdleTimeout,
		cfg.RevokedTTL,
		cfg.RateLimitPeriod, cfg.RateLimitMax,
	)
}
