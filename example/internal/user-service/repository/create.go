package repository

import (
	"context"
	modelUser "example/internal/shared/model"
	mapping "example/internal/user-service/mapping"
)

// Tạo User
func (u *UserRepository) CreateUser(ctx context.Context, user *modelUser.User) (int, error) {

	// map user *modelUser.User -> Entity
	entityU := mapping.FromModelToEntityUser(user)
	// Tạo mới với GORM
	result := u.db.WithContext(ctx).Select("email", "name").Create(entityU)
	if result.Error != nil {
		return 0, result.Error
	}

	// Lấy ID vừa tạo
	return int(user.Id), nil
}
