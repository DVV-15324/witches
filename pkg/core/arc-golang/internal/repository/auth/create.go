package auth

import (
	entity "arc-golang/internal/entity/auth"
	"context"
	"database/sql"
)

// tạo auth
func (u *RepositoryAuth) CreateAuth(cxt context.Context, auth *entity.Auth) error {
	query := `INSERT INTO auths(salt, email, password, user_id, role, banned)
			values(@salt, @email, @password, @user_id, @role, @banned);`
	_, err := u.db.ExecContext(cxt, query,
		sql.Named("salt", auth.Salt),
		sql.Named("email", auth.Email),
		sql.Named("password", auth.Password),
		sql.Named("user_id", auth.UserId),
		sql.Named("role", "user"),
		sql.Named("banned", false),
	)
	if err != nil {
		return err
	}
	return nil
}
