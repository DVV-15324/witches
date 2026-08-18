package usecase

import (
	"context"
	"time"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	domainBook "example/internal/shared/domain"
)

func (u *BookUsecase) Update(ctx context.Context, id int, req *domainBook.Book) (*w_resp.AppError) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return resp
	}

	// Cập nhật field
	// existing.field = req.field
	// ...

	err = u.repo.Update(ctx, existing, id)
	if err != nil {
		resp := w_resp.NewAppError(500, err, time.Now())
		return  resp
	}

	return nil
}