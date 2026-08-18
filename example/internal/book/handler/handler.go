package handler

import (
	"context"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	domainBook "example/internal/shared/domain"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"

)
type IBookUsecase interface {
	Create(ctx context.Context, req *domainBook.Book) (*w_resp.AppError)
	GetByID(ctx context.Context, id int) (*domainBook.Book, *w_resp.AppError)
	GetAll(ctx context.Context, req *w_utils.PaginationRequest) ([]*domainBook.Book, int64, *w_resp.AppError)
	Update(ctx context.Context, id int, req *domainBook.Book) (*w_resp.AppError)
	Delete(ctx context.Context, id int) *w_resp.AppError
}
type BookHandler struct {
	Log         *logger.ModelLogger
	usecase IBookUsecase
	Cfg     *wcmd_utils.Config
}

// NewBookHandler - Constructor
func NewBookHandler(usecase IBookUsecase, log *logger.ModelLogger, Cfg *wcmd_utils.Config) *BookHandler {
	return &BookHandler{
	usecase: usecase,
	Log:         log,
	Cfg         :Cfg,
	}
}