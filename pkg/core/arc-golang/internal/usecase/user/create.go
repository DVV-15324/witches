package user

import (
	c_response "arc-golang/internal/utils"
	"time"

	entityUser "arc-golang/internal/entity/user"
	"context"
)

func (u *UsecaseUser) UsecaseCreateUser(ctx context.Context, user *entityUser.User) (int, *c_response.ErrorResponse) {
	uid, err := u.ReponsitoryUser.CreateUser(ctx, user)
	if err != nil {
		resp := c_response.NewErrorResponse(500, err, time.Now())
		return 0, resp
	}
	return uid, nil
}
