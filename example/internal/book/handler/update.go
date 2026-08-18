package handler

import (
	"example/internal/book/dto/request"
	domainBook "example/internal/shared/domain"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func (h *BookHandler) Update() func(*gin.Context) {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			resp := w_resp.NewAppError(400, err, time.Now())
			w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
			return
		}

		var req request.UpdateBookRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			resp := w_resp.NewAppError(400, err, time.Now())
			w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
			return
		}

		domain := domainBook.Book{}
		resp := h.usecase.Update(ctx.Request.Context(), id, &domain)

		if resp != nil {
			w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
			return
		}
		w_resp.WriteSuccess(ctx, "cap nhat thanh cong")
	}
}
