package auth

import (
	"context"
	"time"

	c_response "arc-golang/internal/utils"

	entityAuth "arc-golang/internal/entity/auth"
	"arc-golang/internal/entity/user"
	c_hash "arc-golang/internal/utils"
)

func (a *UsecaseAuth) UsecaseRegister(ctx context.Context, auth *entityAuth.Auth, name string) *c_response.ErrorResponse {
	// check tài khoản tồn tại
	authEmail, err := a.ReponsitoryAuth.GetAuthByEmail(ctx, auth.Email)
	if authEmail != nil {
		resp := c_response.NewErrorResponse(409, ErrorEmailIsExisted, time.Now())
		return resp
	}
	if err != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return resp
	}
	//Hash Password
	randomString, err := c_hash.RandomStr()
	if err != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return resp
	}

	hash, err := a.Hash.GenerateFromPassword(auth.Password, randomString)
	if err != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return resp
	}

	// Tạo User
	userUid, err_user := a.UsecaseUser.UsecaseCreateUser(ctx, &user.User{
		Email: auth.Email,
		Name:  name,
	})

	if err_user != nil {
		return err_user
	}

	// Tạo Auth
	err_auth := a.ReponsitoryAuth.CreateAuth(ctx, &entityAuth.Auth{
		Salt:     randomString,
		Email:    auth.Email,
		Password: hash,
		UserId:   userUid,
	})

	if err_auth != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return resp
	}
	return nil
}
