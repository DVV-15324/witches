package usecase

import (
	"context"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	modelBook "example/internal/shared/domain"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"time"
)

func (u *BookUsecase) GetByID(ctx context.Context, id int) (*modelBook.Book, *w_resp.AppError) {
	e, err := u.repo.GetByID(ctx, id)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, resp
	}
	return e, nil
}

func (u *BookUsecase) GetAll(ctx context.Context, req *w_utils.PaginationRequest) ([]*modelBook.Book, int64, *w_resp.AppError) {
	list, total, err := u.repo.GetAll(ctx, req)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return nil, 0, resp
	}
	return list, total, nil
}