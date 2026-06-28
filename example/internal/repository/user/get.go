package user

import (
	"context"
	"database/sql"

	entity "example/internal/entity/user"
)

func (u *RepositoryUser) GetAllUser(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT 
		u.id, 
		u.name,
		u.email,
		u.role,
		a.banned
	FROM users u
	LEFT JOIN auths a ON a.user_id = u.id`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var data entity.User
		err := rows.Scan(&data.Id, &data.Name, &data.Email, &data.Role, &data.Banned)
		if err != nil {
			return nil, err
		}
		users = append(users, &data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *RepositoryUser) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	var data entity.User
	query := `SELECT 
		u.id, 
		u.name, 
		u.email,
		u.role
	FROM users u
	WHERE u.id = ?`

	row := u.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&data.Id, &data.Name, &data.Email, &data.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}
