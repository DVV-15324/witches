package repository

import (
	"context"

	"example/internal/auth-service/mapping"
	modelAuth "example/internal/shared/model"
)

func (r *AuthRepository) CreateAuth(ctx context.Context, modelAuth *modelAuth.Auth) error {
	auth := mapping.FromModelToEntityAuth(modelAuth)
	result := r.db.WithContext(ctx).
		Select("email", "password", "user_id", "salt").
		Create(auth)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
