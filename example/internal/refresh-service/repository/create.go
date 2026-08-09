package repository

import (
	"context"
	mapping "example/internal/refresh-service/mapping"
	modelRefresh "example/internal/shared/model"
)

func (r *RefreshTokenRepository) Create(ctx context.Context, token *modelRefresh.RefreshToken) error {
	entityToken := mapping.FromModelToEntityRefresh(token)
	if err := r.db.WithContext(ctx).Create(entityToken).Error; err != nil {
		return err
	}
	return r.cacheToken(ctx, entityToken)
}
