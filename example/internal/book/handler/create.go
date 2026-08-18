package handler

import (
	"example/internal/book/dto/request"
	domainBook "example/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
	"time"
)

func (h *BookHandler) Create() func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req request.CreateBookRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			resp := w_resp.NewAppError(400, err, time.Now())
			w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
			return
		}
		domain := &domainBook.Book{
			// field of request
		}
		resp := h.usecase.Create(ctx.Request.Context(), domain)
		if resp != nil {
			w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
			return
		}
		w_resp.WriteSuccess(ctx, "tao thang cong")
	}
}
