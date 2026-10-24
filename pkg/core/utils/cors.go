package utils

import (
	"strings"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
)

func Cors(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Kiểm tra xem Origin này có được phép không
		isAllowed := false

		// Nếu cấu hình trong .env là "*", tự động cho phép tất cả các port localhost
		if cfg.CorsAllowOrigins == "*" || cfg.CorsAllowOrigins == "" {
			isAllowed = true
		} else {
			// 2. Nếu cấu hình cụ thể, tách chuỗi bằng dấu phẩy và check xem có khớp không
			allowedOrigins := strings.Split(cfg.CorsAllowOrigins, ",")
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					isAllowed = true
					break
				}
			}
		}

		if isAllowed && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Add("Vary", "Origin")
		}

		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		allowMethods := cfg.CorsAllowMethods
		if allowMethods == "" {
			allowMethods = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)

		allowHeaders := cfg.CorsAllowHeaders
		if allowHeaders == "" {
			allowHeaders = "Content-Type, Authorization"
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
