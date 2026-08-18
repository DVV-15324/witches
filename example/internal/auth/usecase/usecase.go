package usecase

import (
	"context"
	core "example/cmd/server/core"
	domainAuth "example/internal/shared/domain"
	domainRefresh "example/internal/shared/domain"
	domainUser "example/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

type IAuthReponsitory interface {
	CreateAuth(ctx context.Context, auth *domainAuth.Auth) error
	GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, error)
	GetAuthByUserId(ctx context.Context, userID uint32) (*domainAuth.Auth, error)
}
type IRefreshTokenUsecase interface {
	Create(ctx context.Context, token *domainRefresh.RefreshToken) *w_resp.AppError
	GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, *w_resp.AppError)
	Revoke(ctx context.Context, token string, reason string, deviceID string) *w_resp.AppError
}

type IHash interface {
	GenerateFromPassword(password string, salt string) (string, error)
	CompareHashAndPassword(passwordStr string, password string, salt string) bool
}

type IUserUseCase interface {
	CreateUser(ctx context.Context, user *domainUser.User) (int, *w_resp.AppError)
	GetUserById(ctx context.Context, id int) (*domainUser.User, *w_resp.AppError)
}

type IJwt interface {
	ParseToken(ctx context.Context, tokenStr string) (*w_utils.JwtClaims, error)
	IssueAccessToken(ctx context.Context, sub string, tid string) (*w_utils.Token, error)
	IssueRefreshToken(ctx context.Context, sub string, tid string) (*w_utils.Token, error)
	IssueTokenPair(ctx context.Context, sub string, tid string) (*w_utils.TokenResponse, error)
}

type AuthUseCase struct {
	Core                *core.CoreServices
	UserUseCase         IUserUseCase
	AuthReponsitory     IAuthReponsitory
	RefreshTokenUseCase IRefreshTokenUsecase
}

func NewAuthUseCase(
	core *core.CoreServices,
	authReponsitory IAuthReponsitory,
	userUseCase IUserUseCase,
	refreshTokenUseCase IRefreshTokenUsecase,
) *AuthUseCase {
	return &AuthUseCase{
		Core:                core,
		AuthReponsitory:     authReponsitory,
		RefreshTokenUseCase: refreshTokenUseCase,
	}
}
