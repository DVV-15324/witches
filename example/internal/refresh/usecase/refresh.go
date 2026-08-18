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

type IRefreshTokenRepository interface {
	Create(ctx context.Context, token *domainRefresh.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*domainRefresh.RefreshToken, error)
	Revoke(ctx context.Context, token string, reason string) error
}

type IJwt interface {
	ParseToken(ctx context.Context, tokenStr string) (*w_utils.JwtClaims, error)
	IssueAccessToken(ctx context.Context, sub string, tid string) (*w_utils.Token, error)
	IssueRefreshToken(ctx context.Context, sub string, tid string) (*w_utils.Token, error)
	IssueTokenPair(ctx context.Context, sub string, tid string) (*w_utils.TokenResponse, error)
}

type IAuthUsecase interface {
	GetAuthByUserId(ctx context.Context, userId int) (*domainAuth.Auth, *w_resp.AppError)
	GetAuthByEmail(ctx context.Context, email string) (*domainAuth.Auth, *w_resp.AppError)
}

type IUserUseCase interface {
	CreateUser(ctx context.Context, user *domainUser.User) (int, *w_resp.AppError)
	GetUserById(ctx context.Context, id int) (*domainUser.User, *w_resp.AppError)
}

type RefreshUseCase struct {
	Core             *core.CoreServices
	AuthUsecase      IAuthUsecase
	RefreshTokenRepo IRefreshTokenRepository
	UserUseCase      IUserUseCase
}

func NewRefreshUseCase(core *core.CoreServices, userUseCase IUserUseCase, refreshTokenRepo IRefreshTokenRepository) *RefreshUseCase {
	return &RefreshUseCase{
		UserUseCase:      userUseCase,
		RefreshTokenRepo: refreshTokenRepo,
		Core:             core,
	}
}

func (r *RefreshUseCase) SetAuthUseCase(authUsecase IAuthUsecase) {
	r.AuthUsecase = authUsecase
}
