package user

import (
	entityUser "example/internal/entity/user"
	"context"
)

// interface
type IResponsitoryUser interface {
	CreateUser(cxt context.Context, user *entityUser.User) (int, error)
	GetUserById(ctx context.Context, id int) (*entityUser.User, error)
	GetAllUser(ctx context.Context) ([]*entityUser.User, error)
}

type UsecaseUser struct {
	ReponsitoryUser IResponsitoryUser
}

func NewUsecaseUser(iReponsitoryUser IResponsitoryUser) *UsecaseUser {
	return &UsecaseUser{
		ReponsitoryUser: iReponsitoryUser,
	}
}
