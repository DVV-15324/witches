package usecase

import (
	"context"
	"time"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	domainBook "example/internal/shared/domain"
)

func (u *BookUsecase) Create(ctx context.Context, req *domainBook.Book) (*w_resp.AppError) {
	 err := u.repo.Create(ctx, req)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return resp
	}

	return nil
}