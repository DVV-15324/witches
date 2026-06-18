package user

import (
	c_response "arc-golang/internal/utils"
	"time"

	entityUser "arc-golang/internal/entity/user"
	"context"
)

func (u *UsecaseUser) UsecaseGetAllUser(ctx context.Context) ([]*entityUser.User, *c_response.ErrorResponse) {
	users, err := u.ReponsitoryUser.GetAllUser(ctx)
	if err != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return nil, resp
	}
	return users, nil
}

func (u *UsecaseUser) UsecaseGetUserById(ctx context.Context, id int) (*entityUser.User, *c_response.ErrorResponse) {
	user, err := u.ReponsitoryUser.GetUserById(ctx, id)
	if err != nil {
		resp := c_response.NewErrorResponse(404, err, time.Now())
		return nil, resp
	}
	return user, nil
}
