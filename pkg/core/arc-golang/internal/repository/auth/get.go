package auth

import (
	entity "arc-golang/internal/entity/auth"
	"context"
	"database/sql"
)

func (a *RepositoryAuth) GetAuthByEmail(ctx context.Context, email string) (*entity.Auth, error) {
	var data entity.Auth
	query := "SELECT id, email, password, user_id, salt, role, banned FROM auths WHERE email = @email"
	err := a.db.QueryRowContext(ctx, query, sql.Named("email", email)).Scan(
		&data.Id,
		&data.Email,
		&data.Password,
		&data.UserId,
		&data.Salt,
		&data.Role,
		&data.Banned,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}
