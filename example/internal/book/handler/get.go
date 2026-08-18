package handler

import (
	"time"
	mapping "example/internal/book/mapping"
	"strconv"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

func (h *BookHandler) GetByID() func(*gin.Context) {
	return func(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := w_resp.NewAppError(400, err, time.Now())
		w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, resp)
		return
	}

	data, err_get := h.usecase.GetByID(ctx.Request.Context(), id)
	if err_get != nil {
		w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, err_get)
		return
	}
	dto := mapping.FromDomainToDtoBook(data)
	w_resp.WriteSuccess(ctx, dto)
}
}


func (h *BookHandler) GetAll() func(*gin.Context) {
	return func(ctx *gin.Context) {
	req := w_utils.PaginationRequest{}
	data, total, err := h.usecase.GetAll(ctx.Request.Context(), &req)
	if err != nil {
		w_resp.WriteErrorWithLog(ctx, h.Log, h.Cfg, err)
		return
	}
	dto := mapping.FromDomainToDtoBookList(data)
	paginationResp := req.ToPaginationResponse(total)
	w_resp.WriteSuccessWithPagination(ctx, dto, &paginationResp)
}
}