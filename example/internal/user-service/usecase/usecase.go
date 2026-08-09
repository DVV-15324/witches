package usecase

import (
	"context"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	modelUser "example/internal/shared/model"
)

// interface
type IUserResponsitory interface {
	CreateUser(cxt context.Context, user *modelUser.User) (int, error)
	GetUserById(ctx context.Context, id int) (*modelUser.User, error)
	GetAllUser(ctx context.Context, req *w_utils.PaginationRequest) ([]*modelUser.User, int64, error)
}

type UserUseCase struct {
	TxManager        w_utils.TxManager
	UserResponsitory IUserResponsitory
}

func NewUserUsecase(UserResponsitory IUserResponsitory, txManager w_utils.TxManager) *UserUseCase {
	return &UserUseCase{
		UserResponsitory: UserResponsitory,
		TxManager:        txManager,
	}
}
