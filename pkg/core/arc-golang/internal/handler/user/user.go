package user

import (
	c_response "arc-golang/internal/utils"

	entityUser "arc-golang/internal/entity/user"
	"context"
)

type IUsecaseUser interface {
	//UsecaseCreateUser(ctx context.Context, user *entityUser.User) (int, *c_response.Response)
	UsecaseGetAllUser(ctx context.Context) ([]*entityUser.User, *c_response.ErrorResponse)
	UsecaseGetUserById(ctx context.Context, id int) (*entityUser.User, *c_response.ErrorResponse)
}
type HandleUser struct {
	UsecaseUser IUsecaseUser
}

func NewHandleUser(iUsecaseUser IUsecaseUser) *HandleUser {
	return &HandleUser{
		UsecaseUser: iUsecaseUser,
	}
}
