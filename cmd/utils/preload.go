package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	godotenv "github.com/joho/godotenv"
)

type Config struct {
	Port            int64
	Host            string
	MetrictPort     int64
	DBDriver        string
	DBHost          string
	DBPort          int64
	DBUser          string
	DBName          string
	DBPassword      string
	DBUrl           string
	MaxOpenConns    int64
	MaxIdleConns    int64
	ConnMaxLifetime int64
	ConnMaxIdleTime int64
	RedisHost       string
	RedisPort       int64
	RedisPassword   string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
	SessionTTL      int64
	RevokedTTL      int64
	IdleTimeout     int64
	RateLimitPeriod int64
	RateLimitMax    int64
}

func LoadDbUrl(cfg *Config) *Config {
	cfg.DBUrl = os.Getenv("DB_URL")
	return cfg
}

func PreloadNotDBURL() *Config {
	currentPath := GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Error: not found witches.env")
	}

	return &Config{
		Port:            getEnvAsInt64("APP_PORT"),
		Host:            os.Getenv("APP_HOST"),
		MetrictPort:     getEnvAsInt64("METRIC_PORT"),
		DBDriver:        os.Getenv("DB_DRIVER"),
		DBUser:          os.Getenv("DB_USER"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          getEnvAsInt64("DB_PORT"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		MaxOpenConns:    getEnvAsInt64("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    getEnvAsInt64("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: getEnvAsInt64("DB_CONN_MAX_LIFETIME"),
		ConnMaxIdleTime: getEnvAsInt64("DB_CONN_MAX_IDLE_TIME"),
		DBName:          os.Getenv("DB_NAME"),
		RedisHost:       os.Getenv("REDIS_HOST"),
		RedisPort:       getEnvAsInt64("REDIS_PORT"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		AccessTokenTTL:  getEnvAsInt64("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: getEnvAsInt64("REFRESH_TOKEN_TTL"),
		SessionTTL:      getEnvAsInt64("SESSION_TTL"),
		RevokedTTL:      getEnvAsInt64("REVOKED_TTL"),
		IdleTimeout:     getEnvAsInt64("IDLE_TIMEOUT"),
		RateLimitPeriod: getEnvAsInt64("RATE_LIMIT_PERIOD"),
		RateLimitMax:    getEnvAsInt64("RATE_LIMIT_MAX"),
	}
}

func getEnvAsInt64(key string) int64 {
	strValue := os.Getenv(key)

	val, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Printf("Error: Get key: %v", key)
	}
	return val
}

func Validate(cfg *Config) {
	if cfg.Port == 0 {
		log.Fatal("Error: APP_PORT is not set in environment")
	}
	if cfg.Host == "" {
		log.Fatal("Error: APP_HOST is not set in environment")
	}
	if cfg.MetrictPort == 0 {
		log.Fatal("Error: METRIC_PORT is not set in environment")
	}
	if cfg.DBHost == "" {
		log.Fatal("Error: DB_HOST is not set in environment")
	}
	if cfg.DBPort == 0 {
		log.Fatal("Error: DB_PORT is not set in environment")
	}
	if cfg.DBName == "" {
		log.Fatal("Error: DB_NAME is not set in environment")
	}
	if cfg.DBPassword == "" {
		log.Fatal("Error: DB_PASSWORD is not set in environment")
	}
	if cfg.MaxOpenConns == 0 {
		log.Fatal("Error: DB_MAX_OPEN_CONNS is not set in environment")
	}
	if cfg.MaxIdleConns == 0 {
		log.Fatal("Error: DB_MAX_IDLE_CONNS is not set in environment")
	}
	if cfg.ConnMaxLifetime == 0 {
		log.Fatal("Error: DB_CONN_MAX_LIFETIME is not set in environment")
	}
	if cfg.ConnMaxIdleTime == 0 {
		log.Fatal("Error: DB_CONN_MAX_IDLE_TIME is not set in environment")
	}
	if cfg.RedisHost == "" {
		log.Fatal("Error: REDIS_HOST is required for refresh project")
	}
	if cfg.RedisPort == 0 {
		log.Fatal("Error: REDIS_PORT is required for refresh project")
	}
	if cfg.AccessTokenTTL == 0 {
		log.Fatal("Error: ACCESS_TOKEN_TTL not set")
	}
	if cfg.RefreshTokenTTL == 0 {
		log.Fatal("Error: REFRESH_TOKEN_TTL not set")
	}
	if cfg.SessionTTL == 0 {
		log.Fatal("Error: SESSION_TTL not set")
	}
	if cfg.IdleTimeout == 0 {
		log.Fatal("Error: IDLE_TIMEOUT not set")
	}
	if cfg.RevokedTTL == 0 {
		log.Fatal("Error: REVOKED_TTL not set")
	}

	if cfg.RateLimitPeriod == 0 {
		log.Fatal("Error: RATE_LIMIT_PERIOD not set")
	}
	if cfg.RateLimitMax == 0 {
		log.Fatal("Error: RATE_LIMIT_MAX not set")
	}
}
