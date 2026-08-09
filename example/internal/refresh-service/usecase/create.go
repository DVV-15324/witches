package usecase

import (
	"context"
	"errors"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"time"
	modelRefresh "example/internal/shared/model"
)

// Create - Tạo refresh token mới
func (uc *RefreshUseCase) Create(ctx context.Context, token *modelRefresh.RefreshToken) *w_resp.AppError {
	if token == nil {
		return w_resp.NewAppError(400, errors.New("token cannot be nil"), time.Now())
	}

	// Validate required fields
	if token.UserID == 0 {
		return w_resp.NewAppError(400, errors.New("user_id is required"), time.Now())
	}
	if token.Token == "" {
		return w_resp.NewAppError(400, errors.New("token is required"), time.Now())
	}
	if token.ExpiresAt == 0 {
		token.ExpiresAt = time.Now().Unix() + uc.Config.RefreshTokenTTL
	}

	// Set default values
	token.Revoked = false
	token.RevokedAt = 0

	// Save to repository (có cache)
	if err := uc.RefreshTokenRepo.Create(ctx, token); err != nil {
		return w_resp.NewAppError(500, err, time.Now())
	}

	return nil
}
