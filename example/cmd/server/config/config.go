package config

import (
	"log"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Host            string
	MetrictPort     string
	DBHost          string
	DBPort          string
	DBPassword      string
	DBName          string
	DBDriver        string
	DBUser          string
	DBURL           string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
	SessionTTL      int64
	RevokedTTL      int64
	IdleTimeout     int64
	RateLimitPeriod int64
	RateLimitMax    int64
}

// En: Load is a function used to initialize and return the application's configuration.
//
// Vi: Load là hàm khởi tạo và trả về cấu hình của ứng dụng.
func Load() *Config {
	if err := godotenv.Load("witches.env"); err != nil {
		log.Println("Error: not found witches.env")
	}

	return &Config{
		Port:            os.Getenv("APP_PORT"), // En: Application server port | Vi: Cổng máy chủ ứng dụng
		Host:            os.Getenv("APP_HOST"),
		MetrictPort:     os.Getenv("MESTRICT_PORT"),
		DBDriver:        os.Getenv("DB_DRIVER"), // En: Database driver (e.g., postgres, mysql) | Vi: Loại database (vd: postgres, mysql)
		DBUser:          os.Getenv("DB_USER"),
		DBHost:          os.Getenv("DB_HOST"),               // En: Database server host | Vi: Địa chỉ máy chủ database
		DBPort:          os.Getenv("DB_PORT"),               // En: Database server port | Vi: Cổng kết nối database
		DBPassword:      os.Getenv("DB_PASSWORD"),           // En: Database password | Vi: Mật khẩu database
		DBName:          os.Getenv("DB_NAME"),               // En: Database name | Vi: Tên database
		DBURL:           os.Getenv("DB_URL"),                // En: Database connection string (DSN) | Vi: Chuỗi kết nối database (DSN)
		RedisHost:       os.Getenv("REDIS_HOST"),            // En: Redis server host | Vi: Địa chỉ máy chủ Redis
		RedisPort:       os.Getenv("REDIS_PORT"),            // En: Redis server port | Vi: Cổng máy chủ Redis
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),        // En: Redis password | Vi: Mật khẩu Redis
		AccessTokenTTL:  getEnvAsInt64("ACCESS_TOKEN_TTL"),  // En: Access token expiration time (in seconds) | Vi: Thời gian hết hạn Access token (giây)
		RefreshTokenTTL: getEnvAsInt64("REFRESH_TOKEN_TTL"), // En: Refresh token expiration time (in seconds) | Vi: Thời gian hết hạn Refresh token (giây)
		SessionTTL:      getEnvAsInt64("SESSION_TTL"),       // En: Session expiration time (in seconds) | Vi: Thời gian hết hạn phiên (giây)
		RevokedTTL:      getEnvAsInt64("REVOKED_TTL"),       // En: Revoked token cache expiration time (in seconds) | Vi: Thời gian token bị thu hồi được lưu cache (giây)
		IdleTimeout:     getEnvAsInt64("IDLE_TIMEOUT"),      // En: HTTP server idle timeout (in seconds) | Vi: Thời gian chờ khi kết nối không hoạt động (giây)
		RateLimitPeriod: getEnvAsInt64("RATE_LIMIT_PERIOD"), // En: Rate limit time window | Vi: Khoảng thời gian giới hạn tốc độ
		RateLimitMax:    getEnvAsInt64("RATE_LIMIT_MAX"),    // En: Maximum requests allowed | Vi: Số lượng request tối đa
	}
}

// En: Converts string fields to int64
//
// Vi: Chuyển các trường string thành int64
func getEnvAsInt64(key string) int64 {
	strValue := os.Getenv(key)

	val, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Printf("Error: Get key: %v", key)
	}
	return val
}
