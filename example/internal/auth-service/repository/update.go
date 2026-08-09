package repository

import (
	"context"

	entityAuth "example/internal/auth-service/entity"
	"example/internal/auth-service/mapping"
	modelAuth "example/internal/shared/model"
)

func (r *AuthRepository) UpdateAuth(ctx context.Context, modleAuth *modelAuth.Auth) error {
	auth := mapping.FromModelToEntityAuth(modleAuth)
	// Lấy old auth để invalidate cache
	var oldAuth entityAuth.Auth
	r.db.WithContext(ctx).Where("email = ?", auth.Email).First(&oldAuth)

	result := r.db.WithContext(ctx).Model(&entityAuth.Auth{}).
		Where("email = ?", auth.UserId).
		Updates(map[string]interface{}{
			"password":  auth.Password,
			"salt":      auth.Salt,
			"banned":    auth.Banned,
			"auth_type": auth.AuthType,
		})
	if result.Error != nil {
		return result.Error
	}

	// // Invalidate cache cũ và cache mới
	go func() {
		r.invalidateAuthCache(context.Background(), oldAuth.Email, uint32(oldAuth.UserId))
		r.cacheAuth(context.Background(), auth)
	}()

	return nil
}
