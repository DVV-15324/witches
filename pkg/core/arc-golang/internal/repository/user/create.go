package user

import (
	entity "arc-golang/internal/entity/user"
	"context"
	"database/sql"
)

// Tạo User
func (u *RepositoryUser) CreateUser(cxt context.Context, user *entity.User) (int, error) {
	query := `INSERT INTO users(name, email, role)
			values(@name, @email, @role)
			RETURNING id;`
	var uid int
	err := u.db.QueryRowContext(cxt, query,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("role", "user"),
	).Scan(&uid)
	if err != nil {
		return 0, err
	}
	return uid, nil
}
