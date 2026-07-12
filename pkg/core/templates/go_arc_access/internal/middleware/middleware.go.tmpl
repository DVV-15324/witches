package middleware

import (
	"context"
	"errors"

	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"
	"github.com/DVV-15324/witches/pkg/core/response_logger/logger"
	u_cxt "github.com/DVV-15324/witches/pkg/core/utils"
	u_jwt "github.com/DVV-15324/witches/pkg/core/utils"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UsecaseAuth interface {
	UsecaseIntrospectToken(ctx context.Context, accessToken string) (*u_jwt.JwtClaims, error)
}

func RequiredAuth(usecaseAuth UsecaseAuth, log *logger.EntityLogger) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy token ở header
		token, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			resp := u_response.NewErrorResponse(401, err, time.Now())
			u_response.WriteError(c, log, resp)
			c.Abort()
			return
		}
		// xác thực token
		claims, er := usecaseAuth.UsecaseIntrospectToken(c, token)
		if er != nil {
			resp := u_response.NewErrorResponse(401, err, time.Now())
			u_response.WriteError(c, log, resp)

			c.Abort()
			return
		}

		//New Request
		requestWithContext := u_cxt.NewRequestResponse(claims.Subject, claims.ID)
		// Lưu vào context
		c.Request = c.Request.WithContext(u_cxt.SaveRequestContext(c, *requestWithContext))
		c.Next()
	}
}

// RequiredAuth - Middleware xác thực token từ header
func RequiredAuthHeader(usecaseAuth UsecaseAuth, log *logger.EntityLogger) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy token ở header
		token, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			resp := u_response.NewErrorResponse(401, err, time.Now())
			u_response.WriteError(c, log, resp)
			c.Abort()
			return
		}
		// xác thực token
		claims, er := usecaseAuth.UsecaseIntrospectToken(c, token)
		if er != nil {
			resp := u_response.NewErrorResponse(401, err, time.Now())
			u_response.WriteError(c, log, resp)

			c.Abort()
			return
		}

		//New Request
		requestWithContext := u_cxt.NewRequestResponse(claims.Subject, claims.ID)
		// Lưu vào context
		c.Request = c.Request.WithContext(u_cxt.SaveRequestContext(c, *requestWithContext))
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
