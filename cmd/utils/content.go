package cmd_utils

import "fmt"

func CreateContentAccess(DB_DRIVER string) string {
	return fmt.Sprintf(`APP_PORT=                                                                                   # En: Application server port | Vi: Cổng máy chủ ứng dụng

DB_DRIVER=%s                                                                                # En: Database driver (e.g., postgres, mysql) | Vi: Loại database (vd: postgres, mysql)
DB_HOST=                                                                                    # En: Database server host | Vi: Địa chỉ máy chủ database
DB_PORT=                                                                                    # En: Database server port | Vi: Cổng kết nối database
DB_NAME=                                                                                    # En: Database name | Vi: Tên database
DB_PASSWORD=                                                                                # En: Database password | Vi: Mật khẩu database


# JWT
ACCESS_TOKEN_TTL=86400        # 1 day                                                    # En: Access token expiration time (seconds) | Vi: Thời gian hết hạn Access Token (giây)

TYPE=Access
`, DB_DRIVER)
}

func CreateContentRefresh(DB_DRIVER string) string {
	return fmt.Sprintf(`APP_PORT=                                                                                   # En: Application server port | Vi: Cổng máy chủ ứng dụng

DB_DRIVER=%s                                                                                  # En: Database driver (e.g., postgres, mysql) | Vi: Loại database (vd: postgres, mysql)
DB_HOST=                                                                                    # En: Database server host | Vi: Địa chỉ máy chủ database
DB_PORT=    							                                                    # En: Database server port | Vi: Cổng kết nối database                                                                            # En: Database name | Vi: Tên database 										
DB_PASSWORD=                                                                                # En: Database password | Vi: Mật khẩu database
DB_NAME=                                                                                    # En: Database name | Vi: Tên database


REDIS_HOST=                                                                                 # En: Redis server host | Vi: Địa chỉ máy chủ Redis
REDIS_PORT=                                                                                 # En: Redis server port | Vi: Cổng máy chủ Redis
REDIS_PASSWORD=                                                                             # En: Redis password | Vi: Mật khẩu Redis

# Token expiration
ACCESS_TOKEN_TTL=900        # 15 minutes                                                    # En: Access token expiration time (in seconds) | Vi: Thời gian hết hạn Access token (giây)
REFRESH_TOKEN_TTL=604800    # 7 days                                                        # En: Refresh token expiration time (in seconds) | Vi: Thời gian hết hạn Refresh token (giây)

# Session settings
SESSION_TTL=604800          # Session lifetime per device (should match REFRESH_TOKEN_TTL)  # En: Session expiration time (in seconds) | Vi: Thời gian hết hạn phiên (giây)
IDLE_TIMEOUT=1800           # Revoke refresh token after 30 minutes of inactivity           # En: HTTP server idle timeout (in seconds) | Vi: Thời gian chờ khi kết nối không hoạt động (giây)

# Blacklist settings
REVOKED_TTL=900             # Blacklist TTL (should match ACCESS_TOKEN_TTL)                 # En: Revoked token cache expiration time (in seconds) | Vi: Thời gian token bị thu hồi được lưu cache (giây)

# Google OAuth
GOOGLE_CLIENT_ID=                                                                            #

RATE_LIMITD_RIVER=                                                                            #"memory" hoặc "redis"                                                                         # En: Redis password | Vi: Mật khẩu Redis

TYPE=Refresh
`, DB_DRIVER)
}

func CreateContentAccessUsed(DB_DRIVER string, DB_HOST string, DB_PORT string, DB_NAME string, DB_PASSWORD string, DB_URL string, ACCESS_TOKEN_TTL string, TYPE string) string {
	return fmt.Sprintf(`APP_PORT=                                                                                   # En: Application server port | Vi: Cổng máy chủ ứng dụng

DB_DRIVER=%s                                                                                # En: Database driver (e.g., postgres, mysql) | Vi: Loại database (vd: postgres, mysql)
DB_HOST=%s                                                                                  # En: Database server host | Vi: Địa chỉ máy chủ database
DB_PORT=%s                                                                                  # En: Database server port | Vi: Cổng kết nối database
DB_NAME=%s                                                                                  # En: Database name | Vi: Tên database
DB_PASSWORD=%s                                                                              # En: Database password | Vi: Mật khẩu database
DB_URL=%s                                                                                   # En: Database connection string (DSN) | Vi: Chuỗi kết nối database (DSN)

# JWT
ACCESS_TOKEN_TTL=%s         # 1 day                                                         # En: Access token expiration time (seconds) | Vi: Thời gian hết hạn Access Token (giây)

TYPE=%s
`, DB_DRIVER, DB_HOST, DB_PORT, DB_NAME, DB_PASSWORD, DB_URL, ACCESS_TOKEN_TTL, TYPE)
}

func CreateContentRefreshUsed(
	DB_DRIVER string,
	DB_HOST string,
	DB_PORT string,
	DB_NAME string,
	DB_PASSWORD string,
	DB_URL string,
	REDIS_HOST string,
	REDIS_PORT string,
	REDIS_PASSWORD string,
	ACCESS_TOKEN_TTL string,
	REFRESH_TOKEN_TTL string,
	SESSION_TTL string,
	IDLE_TIMEOUT string,
	REVOKED_TTL string,
	GOOGLE_CLIENT_ID string,
	RATE_LIMIT_DRIVER string,
	TYPE string,
) string {
	return fmt.Sprintf(`APP_PORT=                                                                                   # En: Application server port | Vi: Cổng máy chủ ứng dụng

DB_DRIVER=%s                                                                                # En: Database driver (e.g., postgres, mysql) | Vi: Loại database (vd: postgres, mysql)
DB_HOST=%s                                                                                  # En: Database server host | Vi: Địa chỉ máy chủ database
DB_PORT=%s                                                                                  # En: Database server port | Vi: Cổng kết nối database
DB_NAME=%s                                                                                  # En: Database name | Vi: Tên database
DB_PASSWORD=%s                                                                              # En: Database password | Vi: Mật khẩu database
DB_URL=%s                                                                                   # En: Database connection string (DSN) | Vi: Chuỗi kết nối database (DSN)

REDIS_HOST=%s                                                                               # En: Redis server host | Vi: Địa chỉ máy chủ Redis
REDIS_PORT=%s                                                                               # En: Redis server port | Vi: Cổng máy chủ Redis
REDIS_PASSWORD=%s                                                                           # En: Redis password | Vi: Mật khẩu Redis

# Token expiration
ACCESS_TOKEN_TTL=%s         # Access token TTL                                               # En: Access token expiration time (in seconds) | Vi: Thời gian hết hạn Access token (giây)
REFRESH_TOKEN_TTL=%s        # Refresh token TTL                                             # En: Refresh token expiration time (in seconds) | Vi: Thời gian hết hạn Refresh token (giây)

# Session settings
SESSION_TTL=%s              # Session lifetime per device (should match REFRESH_TOKEN_TTL)  # En: Session expiration time (in seconds) | Vi: Thời gian hết hạn phiên (giây)
IDLE_TIMEOUT=%s             # Revoke refresh token after inactivity                         # En: HTTP server idle timeout (in seconds) | Vi: Thời gian chờ khi kết nối không hoạt động (giây)

# Blacklist settings
REVOKED_TTL=%s              # Blacklist TTL (should match ACCESS_TOKEN_TTL)                 # En: Revoked token cache expiration time (in seconds) | Vi: Thời gian token bị thu hồi được lưu cache (giây)

# Google OAuth
GOOGLE_CLIENT_ID=%s                                                                         # En: Google OAuth Client ID | Vi: ID ứng dụng Google OAuth
GOOGLE_CLIENT_SECRET=                                                                       # En: Google OAuth Client Secret | Vi: Mật khẩu ứng dụng Google OAuth
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback                              # En: Google OAuth Redirect URL | Vi: URL chuyển hướng OAuth Google

# Rate Limit
RATE_LIMIT_DRIVER=%s                                                                        # En: Rate limit driver (memory or redis) | Vi: Công cụ giới hạn tốc độ (memory hoặc redis)
RATE_LIMIT_REQUESTS=100                                                                     # En: Max requests per window | Vi: Số request tối đa trong khoảng thời gian
RATE_LIMIT_WINDOW=60                                                                        # En: Time window in seconds | Vi: Khoảng thời gian (giây)

TYPE=%s
`,
		DB_DRIVER, DB_HOST, DB_PORT, DB_NAME, DB_PASSWORD, DB_URL,
		REDIS_HOST, REDIS_PORT, REDIS_PASSWORD,
		ACCESS_TOKEN_TTL, REFRESH_TOKEN_TTL,
		SESSION_TTL, IDLE_TIMEOUT,
		REVOKED_TTL,
		GOOGLE_CLIENT_ID,
		RATE_LIMIT_DRIVER,
		TYPE,
	)
}
