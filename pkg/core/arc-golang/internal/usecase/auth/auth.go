package auth

import (
	entityAuth "arc-golang/internal/entity/auth"
	entityUser "arc-golang/internal/entity/user"
	c_jwt "arc-golang/internal/utils"
	c_response "arc-golang/internal/utils"
	"context"
)

type IReponsitoryAuth interface {
	CreateAuth(cxt context.Context, auth *entityAuth.Auth) error
	GetAuthByEmail(ctx context.Context, email string) (*entityAuth.Auth, error)
}

type IHash interface {
	GenerateFromPassword(password string, salt string) (string, error)
	CompareHashAndPassword(passwordStr string, password string, salt string) bool
}

type IUsecaseUser interface {
	UsecaseCreateUser(ctx context.Context, user *entityUser.User) (int, *c_response.ErrorResponse)
	UsecaseGetUserById(ctx context.Context, id int) (*entityUser.User, *c_response.ErrorResponse)
}

type IJwt interface {
	ParseToken(ctx context.Context, tokenStr string) (*c_jwt.JwtClaims, error)
	IssueToken(ctx context.Context, sub string, tid string) (*c_jwt.TokenResponse, error)
}

type UsecaseAuth struct {
	Jwt             IJwt
	UsecaseUser     IUsecaseUser
	Hash            IHash
	ReponsitoryAuth IReponsitoryAuth
}

func NewUsecaseAuth(jwt IJwt, iUsecaseUser IUsecaseUser, hash IHash, reponsitoryAuth IReponsitoryAuth) *UsecaseAuth {
	return &UsecaseAuth{
		Jwt:             jwt,
		UsecaseUser:     iUsecaseUser,
		Hash:            hash,
		ReponsitoryAuth: reponsitoryAuth,
	}
}
