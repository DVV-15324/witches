package utils

import (
	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
)

func Cors(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		allowOrigin := cfg.CorsAllowOrigins
		if allowOrigin == "" {
			allowOrigin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)

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
