package handler

import (
	"context"
	"example/cmd/server/core"
	domainAuth "example/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

type IUsecaseAuth interface {
	Login(ctx context.Context, au *domainAuth.Auth, DeviceID string, IPAddress string, UserAgent string, TimeZone string, Locale string) (*w_utils.TokenResponse, *w_resp.AppError)
	Register(ctx context.Context, auth *domainAuth.Auth, name string) *w_resp.AppError
	Logout(ctx context.Context, accessToken, refreshTokenStr string, deviceID string) *w_resp.AppError
}
type AuthHandle struct {
	Core        *core.CoreServices
	UsecaseAuth IUsecaseAuth
}

func NewAuthHandle(core *core.CoreServices, authUseCase IUsecaseAuth) *AuthHandle {
	return &AuthHandle{
		UsecaseAuth: authUseCase,
		Core:        core,
	}
}
