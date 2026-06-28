package auth

import (
	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"
	"github.com/DVV-15324/witches/pkg/core/response_logger/logger"
	u_jwt "github.com/DVV-15324/witches/pkg/core/utils"

	entityAuth "example/internal/entity/auth"
	"context"
)

type IUsecaseAuth interface {
	UsecaseLogin(ctx context.Context, au *entityAuth.Auth) (*u_jwt.TokenResponse, *u_response.ErrorResponse)
	UsecaseRegister(ctx context.Context, auth *entityAuth.Auth, name string) *u_response.ErrorResponse
}
type HandleAuth struct {
	Log         *logger.EntityLogger
	UsecaseAuth IUsecaseAuth
}

func NewHandleUser(iUsecaseAuth IUsecaseAuth, log *logger.EntityLogger) *HandleAuth {
	return &HandleAuth{
		UsecaseAuth: iUsecaseAuth,
		Log:         log,
	}
}
