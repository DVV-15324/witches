package user

import (
	entity "example/internal/entity/user"
	"context"
)

// Tạo User
func (u *RepositoryUser) CreateUser(cxt context.Context, user *entity.User) (int, error) {
	query := `INSERT INTO users(name, email, role)
			VALUES (?, ?, ?)`
	result, err := u.db.ExecContext(cxt, query,
		user.Name,
		user.Email,
		"user",
	)
	if err != nil {
		return 0, err
	}

	// Lấy ID vừa tạo
	uid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(uid), nil
}
