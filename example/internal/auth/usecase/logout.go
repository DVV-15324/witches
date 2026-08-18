package usecase

import (
	"context"
	"errors"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
)

func (a *AuthUseCase) Logout(ctx context.Context, accessToken, refreshTokenStr string, deviceID string) *w_resp.AppError {
	// 1. Get refresh token info
	token, err := a.RefreshTokenUseCase.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return err
	}
	if token == nil {
		return w_resp.NewAppError(500, errors.New("refresh token không hợp lệ"), time.Now())
	}

	// 2. Revoke refresh token trong DB (CHỈ CẦN NÀY LÀ ĐỦ CHO REFRESH TOKEN)
	if err := a.RefreshTokenUseCase.Revoke(ctx, refreshTokenStr, "user logout", deviceID); err != nil {
		return err
	}
	// 3. Blacklist access token
	a.Core.Blacklist.BlacklistToken(ctx, accessToken)
	// 4. Delete session
	a.Core.Session.DeleteSession(ctx, token.UserID, deviceID)

	return nil
}
