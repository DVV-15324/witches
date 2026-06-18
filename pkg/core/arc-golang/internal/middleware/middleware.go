package middleware

import (
	c_ctx "arc-golang/internal/utils"
	c_jwt "arc-golang/internal/utils"
	c_response "arc-golang/internal/utils"
	"context"
	"errors"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UsecaseAuth interface {
	UsecaseIntrospectToken(ctx context.Context, accessToken string) (*c_jwt.JwtClaims, error)
}

func RequiredAuthAdmin(usecaseAuth UsecaseAuth) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy token ở header
		token, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			resp := c_response.NewErrorResponse(401, err, time.Now())
			c_response.WriteError(c, resp)
			c.Abort()
			return
		}
		// xác thực token
		claims, er := usecaseAuth.UsecaseIntrospectToken(c, token)
		if er != nil {
			resp := c_response.NewErrorResponse(401, err, time.Now())
			c_response.WriteError(c, resp)

			c.Abort()
			return
		}

		//New Request
		requestWithContext := c_ctx.NewRequestResponse(claims.Subject, claims.ID)
		// Lưu vào context
		c.Request = c.Request.WithContext(c_ctx.SaveRequestContext(c, requestWithContext))
		c.Next()
	}
}

// RequiredAuth - Middleware xác thực token từ header
func RequiredAuth(usecaseAuth UsecaseAuth) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy token ở header
		token, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			resp := c_response.NewErrorResponse(401, err, time.Now())
			c_response.WriteError(c, resp)
			c.Abort()
			return
		}
		// xác thực token
		claims, er := usecaseAuth.UsecaseIntrospectToken(c, token)
		if er != nil {
			resp := c_response.NewErrorResponse(401, err, time.Now())
			c_response.WriteError(c, resp)
			c.Abort()
			return
		}

		//New Request
		requestWithContext := c_ctx.NewRequestResponse(claims.Subject, claims.ID)
		// Lưu vào context
		c.Request = c.Request.WithContext(c_ctx.SaveRequestContext(c, requestWithContext))
		c.Next()
	}
}

// extractTokenFromHeader - Trích xuất token từ header Authorization
func extractTokenFromHeader(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("authorization header is required")
	}

	args := strings.Split(accessToken, " ")
	// Thiếu bearer, thiếu token, bị nhiều " "
	if len(args) != 2 || args[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return args[1], nil
}
