package auth

import (
	entityAuth "arc-golang/internal/entity/auth"
	c_jwt "arc-golang/internal/utils"
	c_response "arc-golang/internal/utils"
	c_uid "arc-golang/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

func (a *UsecaseAuth) UsecaseLogin(ctx context.Context, au *entityAuth.Auth) (*c_jwt.TokenResponse, *c_response.ErrorResponse) {
	//Check tai khoan ko ton tai
	auth, err_auth := a.ReponsitoryAuth.GetAuthByEmail(ctx, au.Email)
	if err_auth != nil {
		resp := c_response.NewErrorResponse(404, err_auth, time.Now())
		return nil, resp
	}
	//So sánh Hash Password
	ss := a.Hash.CompareHashAndPassword(auth.Password, au.Password, auth.Salt)
	if !ss {
		resp := c_response.NewErrorResponse(404, err_auth, time.Now())
		return nil, resp
	}
	//Tạo token
	s := c_uid.NewUID(uint32(auth.UserId), 1)
	// Mã hóa Base58
	sub := s.ToBase58()
	// Tạo mã định danh cho tokwn
	tid := uuid.New().String()

	//Tạo token
	token, err_token := a.Jwt.IssueToken(ctx, sub, tid)
	if err_token != nil {
		resp := c_response.NewErrorResponse(500, err_token, time.Now())
		return nil, resp
	}
	return token, nil
}
