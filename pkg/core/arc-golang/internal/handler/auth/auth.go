package auth

import (
	c_jwt "arc-golang/internal/utils"
	c_response "arc-golang/internal/utils"

	entityAuth "arc-golang/internal/entity/auth"
	"context"
)

type IUsecaseAuth interface {
	UsecaseLogin(ctx context.Context, au *entityAuth.Auth) (*c_jwt.TokenResponse, *c_response.ErrorResponse)
	UsecaseRegister(ctx context.Context, auth *entityAuth.Auth, name string) *c_response.ErrorResponse
}
type HandleAuth struct {
	UsecaseAuth IUsecaseAuth
}

func NewHandleUser(iUsecaseAuth IUsecaseAuth) *HandleAuth {
	return &HandleAuth{
		UsecaseAuth: iUsecaseAuth,
	}
}
