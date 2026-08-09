package utils

import "fmt"

func CreateContentRefresh() string {
	return `# SERVER CONFIGURATION
APP_PORT=8080
APP_HOST=localhost
METRIC_PORT=8088

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

# REDIS CONFIGURATION
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# TOKEN EXPIRATION
ACCESS_TOKEN_TTL=900
REFRESH_TOKEN_TTL=604800

# SESSION SETTINGS
SESSION_TTL=604800
IDLE_TIMEOUT=1800

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
METRIC_PORT=%d

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

# REDIS CONFIGURATION
REDIS_HOST=%s
REDIS_PORT=%d
REDIS_PASSWORD=%s

# TOKEN EXPIRATION
ACCESS_TOKEN_TTL=%d
REFRESH_TOKEN_TTL=%d

# SESSION SETTINGS
SESSION_TTL=%d
IDLE_TIMEOUT=%d

# BLACKLIST SETTINGS
REVOKED_TTL=%d

# RATE LIMIT
RATE_LIMIT_PERIOD=%d
RATE_LIMIT_MAX=%d
`, cfg.Port, cfg.Host, cfg.MetrictPort,
		cfg.DBDriver, cfg.DBUser, cfg.DBHost,
		cfg.DBPort, cfg.DBName, cfg.DBPassword, DB_URL,
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime,
		cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword,
		cfg.AccessTokenTTL, cfg.RefreshTokenTTL,
		cfg.SessionTTL, cfg.IdleTimeout,
		cfg.RevokedTTL,
		cfg.RateLimitPeriod,
		cfg.RateLimitMax,
	)
}
