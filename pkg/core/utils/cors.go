package utils

import (
	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func Cors(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		isAllowed := false

		if cfg.CorsAllowOrigins == "*" || cfg.CorsAllowOrigins == "" {
			isAllowed = true
		} else {
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
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Max-Age", "43200")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		}

		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		allowMethods := cfg.CorsAllowMethods
		if allowMethods == "" {
			allowMethods = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)

		allowHeaders := cfg.CorsAllowHeaders
		if allowHeaders == "" {
			allowHeaders = "Content-Type, Authorization, X-Client-Platform"
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
