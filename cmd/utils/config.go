package utils

import (
	godotenv "github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type contextKey string

type Config struct {
	Port               int64
	Host               string
	Env                string
	MetrictPort        int64
	UIDBits            int64
	HashLen            int64
	RequestKey         contextKey
	LogPath            string
	CorsAllowOrigins   string
	CorsAllowMethods   string
	CorsAllowHeaders   string
	Timezone           string
	SupportedLanguages []string
	JWTSecret          string
	DBDriver           string
	DBHost             string
	DBPort             int64
	DBUser             string
	DBName             string
	DBPassword         string
	DBUrl              string
	MaxOpenConns       int64
	MaxIdleConns       int64
	ConnMaxLifetime    int64
	ConnMaxIdleTime    int64
	SlowThreshold      int64
	RedisHost          string
	RedisPort          int64
	RedisPassword      string
	AccessTokenTTL     int64
	RefreshTokenTTL    int64
	RevokedTTL         int64
	IdleTimeout        int64
	RateLimitPeriod    int64
	RateLimitMax       int64
}

func LoadDbUrl(cfg *Config) *Config {
	cfg.DBUrl = os.Getenv("DB_URL")
	return cfg
}

func PreloadNotDBURL() *Config {
	currentPath := GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Error: not found witches.env, using default values")
	}

	cfg := DefaultConfig()

	if val := os.Getenv("APP_PORT"); val != "" {
		cfg.Port = getEnvAsInt64("APP_PORT")
	}
	if val := os.Getenv("APP_HOST"); val != "" {
		cfg.Host = os.Getenv("APP_HOST")
	}

	if val := os.Getenv("ENV"); val != "" {
		if isValidEnv(val) {
			cfg.Env = val
		}
	}
	if val := os.Getenv("METRIC_PORT"); val != "" {
		cfg.MetrictPort = getEnvAsInt64("METRIC_PORT")
	}
	if val := os.Getenv("UID_BITS"); val != "" {
		cfg.UIDBits = int64(getEnvAsInt64("UID_BITS"))
	}
	if val := os.Getenv("HASH_LEN"); val != "" {
		cfg.HashLen = int64(getEnvAsInt64("HASH_LEN"))
	}
	if val := os.Getenv("REQUEST_KEY"); val != "" {
		cfg.RequestKey = contextKey(os.Getenv("REQUEST_KEY"))
	}
	if val := os.Getenv("LOG_PATH"); val != "" {
		cfg.LogPath = os.Getenv("LOG_PATH")
	}
	if val := os.Getenv("CORS_ALLOW_ORIGINS"); val != "" {
		cfg.CorsAllowOrigins = os.Getenv("CORS_ALLOW_ORIGINS")
	}
	if val := os.Getenv("CORS_ALLOW_METHODS"); val != "" {
		cfg.CorsAllowMethods = os.Getenv("CORS_ALLOW_METHODS")
	}
	if val := os.Getenv("CORS_ALLOW_HEADERS"); val != "" {
		cfg.CorsAllowHeaders = os.Getenv("CORS_ALLOW_HEADERS")
	}
	if val := os.Getenv("TIMEZONE"); val != "" {
		cfg.Timezone = os.Getenv("TIMEZONE")
	}
	if val := os.Getenv("SUPPORTED_LANGUAGES"); val != "" {
		cfg.SupportedLanguages = splitAndTrim(os.Getenv("SUPPORTED_LANGUAGES"), ",")
	}
	if val := os.Getenv("JWT_SECRET"); val != "" {
		cfg.JWTSecret = os.Getenv("JWT_SECRET")
	}
	if val := os.Getenv("DB_DRIVER"); val != "" {
		cfg.DBDriver = os.Getenv("DB_DRIVER")
	}
	if val := os.Getenv("DB_USER"); val != "" {
		cfg.DBUser = os.Getenv("DB_USER")
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		cfg.DBHost = os.Getenv("DB_HOST")
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		cfg.DBPort = getEnvAsInt64("DB_PORT")
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		cfg.DBName = os.Getenv("DB_NAME")
	}
	if val := os.Getenv("DB_PASSWORD"); val != "" {
		cfg.DBPassword = os.Getenv("DB_PASSWORD")
	}
	if val := os.Getenv("DB_MAX_OPEN_CONNS"); val != "" {
		cfg.MaxOpenConns = getEnvAsInt64("DB_MAX_OPEN_CONNS")
	}
	if val := os.Getenv("DB_MAX_IDLE_CONNS"); val != "" {
		cfg.MaxIdleConns = getEnvAsInt64("DB_MAX_IDLE_CONNS")
	}
	if val := os.Getenv("DB_CONN_MAX_LIFETIME"); val != "" {
		cfg.ConnMaxLifetime = getEnvAsInt64("DB_CONN_MAX_LIFETIME")
	}
	if val := os.Getenv("DB_CONN_MAX_IDLE_TIME"); val != "" {
		cfg.ConnMaxIdleTime = getEnvAsInt64("DB_CONN_MAX_IDLE_TIME")
	}
	if val := os.Getenv("SLOW_THRESHOLD"); val != "" {
		cfg.SlowThreshold = getEnvAsInt64("SLOW_THRESHOLD")
	}
	if val := os.Getenv("REDIS_HOST"); val != "" {
		cfg.RedisHost = os.Getenv("REDIS_HOST")
	}
	if val := os.Getenv("REDIS_PORT"); val != "" {
		cfg.RedisPort = getEnvAsInt64("REDIS_PORT")
	}
	if val := os.Getenv("REDIS_PASSWORD"); val != "" {
		cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	}
	if val := os.Getenv("ACCESS_TOKEN_TTL"); val != "" {
		cfg.AccessTokenTTL = getEnvAsInt64("ACCESS_TOKEN_TTL")
	}
	if val := os.Getenv("REFRESH_TOKEN_TTL"); val != "" {
		cfg.RefreshTokenTTL = getEnvAsInt64("REFRESH_TOKEN_TTL")
	}
	if val := os.Getenv("IDLE_TIMEOUT"); val != "" {
		cfg.IdleTimeout = getEnvAsInt64("IDLE_TIMEOUT")
	}
	if val := os.Getenv("REVOKED_TTL"); val != "" {
		cfg.RevokedTTL = getEnvAsInt64("REVOKED_TTL")
	}
	if val := os.Getenv("RATE_LIMIT_PERIOD"); val != "" {
		cfg.RateLimitPeriod = getEnvAsInt64("RATE_LIMIT_PERIOD")
	}
	if val := os.Getenv("RATE_LIMIT_MAX"); val != "" {
		cfg.RateLimitMax = getEnvAsInt64("RATE_LIMIT_MAX")
	}

	return cfg
}

func getEnvAsInt64(key string) int64 {
	strValue := os.Getenv(key)
	if strValue == "" {
		return 0
	}
	val, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Printf("Error: Get key: %v, using default value 0", key)
		return 0
	}
	return val
}

func DefaultConfig() *Config {
	return &Config{
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
		DBName:             "your_database",
		DBPassword:         "your_password",
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
}

func splitAndTrim(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func isValidEnv(env string) bool {
	switch env {
	case "development", "test", "production":
		return true
	default:
		return false
	}
}
