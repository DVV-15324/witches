package handler

import (
	"context"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	modelAuth "example/internal/shared/model"
)

type IUsecaseAuth interface {
	Login(ctx context.Context, au *modelAuth.Auth, DeviceID string, IPAddress string, UserAgent string, TimeZone string, Locale string) (*w_utils.TokenResponse, *w_resp.AppError)
	Register(ctx context.Context, auth *modelAuth.Auth, name string) *w_resp.AppError
	LoginWithGoogle(ctx context.Context, accessToken string, deviceID string, sipAddress string, userAgent string, TimeZone string, Locale string) (*w_utils.TokenResponse, *w_resp.AppError)
	Logout(ctx context.Context, accessToken, refreshTokenStr string, deviceID string) *w_resp.AppError
}
type AuthHandle struct {
	Log         *logger.EntityLogger
	UsecaseAuth IUsecaseAuth
}

func NewAuthHandle(authUseCase IUsecaseAuth, log *logger.EntityLogger) *AuthHandle {
	return &AuthHandle{
		UsecaseAuth: authUseCase,
		Log:         log,
	}
}
