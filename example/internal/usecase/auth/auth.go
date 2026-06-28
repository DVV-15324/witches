package auth

import (
	entityAuth "example/internal/entity/auth"
	entityUser "example/internal/entity/user"
	"context"
	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"
	u_jwt "github.com/DVV-15324/witches/pkg/core/utils"
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
	UsecaseCreateUser(ctx context.Context, user *entityUser.User) (int, *u_response.ErrorResponse)
	UsecaseGetUserById(ctx context.Context, id int) (*entityUser.User, *u_response.ErrorResponse)
}

type IJwt interface {
	ParseToken(ctx context.Context, tokenStr string) (*u_jwt.JwtClaims, error)
	IssueToken(ctx context.Context, sub string, tid string, role u_jwt.Role) (*u_jwt.TokenResponse, error)
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
