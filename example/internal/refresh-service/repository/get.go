package repository

import (
	"context"

	"errors"
	"time"

	entityRefresh "example/internal/refresh-service/entity"
	mapping "example/internal/refresh-service/mapping"
	modelRefresh "example/internal/shared/model"

	"gorm.io/gorm"
)

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*modelRefresh.RefreshToken, error) {
	// 1. Check cache
	cached, err := r.getCachedToken(ctx, token)
	if err == nil && cached != nil {
		modelRefresh := mapping.FromEntityToModelRefresh(cached)
		return modelRefresh, nil
	}

	// 2. Query DB
	var rt entityRefresh.RefreshToken
	err = r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Cache "not found" để tránh cache penetration
			r.cacheNotFound(ctx, token)
			return nil, nil
		}
		return nil, err
	}

	// 3. Cache kết quả
	if !rt.Revoked && time.Now().Before(time.Unix(rt.ExpiresAt, 0)) {
		r.cacheToken(ctx, &rt)
	}
	modelRefresh := mapping.FromEntityToModelRefresh(&rt)
	return modelRefresh, nil
}
