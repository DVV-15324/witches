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
	RequestKey         contextKey
	LogPath            string
	CorsAllowOrigins   string
	CorsAllowMethods   string
	CorsAllowHeaders   string
	Timezone           string
	SupportedLanguages []string
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
	RateLimitPeriod    int64
	RateLimitMax       int64
	Others             map[string]string            // cfg.Others["CUSTOM_KEY"] in witches expan
	Domains            map[string]map[string]string // cfg.Domains["book"]["KEY"] in domain.env
}

func LoadDbUrl(cfg *Config) {
	cfg.DBUrl = os.Getenv("DB_URL")
}

func LoadModule(modulePath string, cfg *Config) {
	envPath := filepath.Join(modulePath, "module.env")
	if m, err := godotenv.Read(envPath); err == nil {
		if cfg.Domains == nil {
			cfg.Domains = make(map[string]map[string]string)
		}
		cfg.Domains[modulePath] = m
	}
}

func PreloadNotDBURL() *Config {
	currentPath := GetCurrentPath()
	envPath := filepath.Join(currentPath, "witches.env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Error: not found witches.env, using default values")
	}

	cfg := DefaultConfig()
	cfg.Others = make(map[string]string)

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
	if val := os.Getenv("RATE_LIMIT_PERIOD"); val != "" {
		cfg.RateLimitPeriod = getEnvAsInt64("RATE_LIMIT_PERIOD")
	}
	if val := os.Getenv("RATE_LIMIT_MAX"); val != "" {
		cfg.RateLimitMax = getEnvAsInt64("RATE_LIMIT_MAX")
	}
	// Đọc tất cả biến trong witches.env vào Others
	if m, err := godotenv.Read(envPath); err == nil {
		for k, v := range m {
			// Chỉ lưu những key không có trong struct field
			if !isKnownKey(k) {
				cfg.Others[k] = v
			}
		}
	}
	return cfg
}
func isKnownKey(key string) bool {
	knownKeys := map[string]bool{
		"APP_PORT": true, "APP_HOST": true, "ENV": true,
		"METRIC_PORT": true, "UID_BITS": true, "REQUEST_KEY": true,
		"LOG_PATH": true, "CORS_ALLOW_ORIGINS": true,
		"CORS_ALLOW_METHODS": true, "CORS_ALLOW_HEADERS": true,
		"TIMEZONE": true, "SUPPORTED_LANGUAGES": true,
		"DB_DRIVER": true, "DB_USER": true, "DB_HOST": true,
		"DB_PORT": true, "DB_NAME": true, "DB_PASSWORD": true,
		"DB_MAX_OPEN_CONNS": true, "DB_MAX_IDLE_CONNS": true,
		"DB_CONN_MAX_LIFETIME": true, "DB_CONN_MAX_IDLE_TIME": true,
		"SLOW_THRESHOLD": true,
		"REDIS_HOST":     true, "REDIS_PORT": true, "REDIS_PASSWORD": true,
		"RATE_LIMIT_PERIOD": true, "RATE_LIMIT_MAX": true,
	}
	return knownKeys[key]
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
		RequestKey:         "request_context",
		LogPath:            "./logs",
		CorsAllowOrigins:   "*",
		CorsAllowMethods:   "GET,POST,PUT,DELETE,OPTIONS",
		CorsAllowHeaders:   "Content-Type,Authorization",
		Timezone:           "UTC",
		SupportedLanguages: []string{"en-US", "vi-VN"},
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
