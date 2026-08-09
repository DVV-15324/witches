package utils

import (
	godotenv "github.com/joho/godotenv"
	"log"
	"os"
)

type PreConfig struct {
	Port            string
	Host            string
	MetrictPort     string
	DBHost          string
	DBPort          string
	DBPassword      string
	DBName          string
	DBDriver        string
	DBUser          string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	AccessTokenTTL  string
	RefreshTokenTTL string
	SessionTTL      string
	RevokedTTL      string
	IdleTimeout     string
	RateLimitPeriod string
	RateLimitMax    string
}

func PreloadDBURL() *PreConfig {
	if err := godotenv.Load("witches.env"); err != nil {
		log.Println("Error: not found witches.env")
	}

	return &PreConfig{
		Port:            os.Getenv("APP_PORT"),
		Host:            os.Getenv("APP_HOST"),
		MetrictPort:     os.Getenv("MESTRICT_PORT"),
		DBDriver:        os.Getenv("DB_DRIVER"),
		DBUser:          os.Getenv("DB_USER"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		RedisHost:       os.Getenv("REDIS_HOST"),
		RedisPort:       os.Getenv("REDIS_PORT"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		AccessTokenTTL:  os.Getenv("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: os.Getenv("REFRESH_TOKEN_TTL"),
		SessionTTL:      os.Getenv("SESSION_TTL"),
		RevokedTTL:      os.Getenv("REVOKED_TTL"),
		IdleTimeout:     os.Getenv("IDLE_TIMEOUT"),
		RateLimitPeriod: os.Getenv("RATE_LIMIT_PERIOD"),
		RateLimitMax:    os.Getenv("RATE_LIMIT_MAX"),
	}
}

func Validate(cfg *PreConfig) {
	if cfg.Port == "" {
		log.Fatal("Error: APP_PORT is not set in environment")
	}
	if cfg.Host == "" {
		log.Fatal("Error: APP_HOST is not set in environment")
	}
	if cfg.MetrictPort == "" {
		log.Fatal("Error: METRIC_PORT is not set in environment")
	}
	if cfg.DBHost == "" {
		log.Fatal("Error: DB_HOST is not set in environment")
	}
	if cfg.DBPort == "" {
		log.Fatal("Error: DB_PORT is not set in environment")
	}
	if cfg.DBName == "" {
		log.Fatal("Error: DB_NAME is not set in environment")
	}
	if cfg.DBPassword == "" {
		log.Fatal("Error: DB_PASSWORD is not set in environment")
	}
	if cfg.RedisHost == "" {
		log.Fatal("Error: REDIS_HOST is required for refresh project")
	}
	if cfg.RedisPort == "" {
		log.Fatal("Error: REDIS_PORT is required for refresh project")
	}
	if cfg.RedisPassword == "" {
		log.Fatal("Error: REDIS_PASSWORD is required for refresh project")
	}
	if cfg.AccessTokenTTL == "" {
		log.Fatal("Error: ACCESS_TOKEN_TTL not set")
	}
	if cfg.RefreshTokenTTL == "" {
		log.Fatal("Error: REFRESH_TOKEN_TTL not set")
	}
	if cfg.SessionTTL == "" {
		log.Fatal("Error: SESSION_TTL not set")
	}
	if cfg.IdleTimeout == "" {
		log.Fatal("Error: IDLE_TIMEOUT not set")
	}
	if cfg.RevokedTTL == "" {
		log.Fatal("Error: REVOKED_TTL not set")
	}

	if cfg.RateLimitPeriod == "" {
		log.Fatal("Error: RATE_LIMIT_PERIOD not set")
	}
	if cfg.RateLimitMax == "" {
		log.Fatal("Error: RATE_LIMIT_MAX not set")
	}
}
