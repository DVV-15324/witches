package user

import (
	dtoUser "example/internal/dto/user/response"
	mapping "example/internal/mapping"

	u_ctx "github.com/DVV-15324/witches/pkg/core/utils"

	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"
	u_uid "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
)

func (h *HandleUser) HandleGetAllUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		user, resp := h.UsecaseUser.UsecaseGetAllUser(c)
		if resp != nil {
			u_response.WriteError(c, h.Log, resp)
			return
		}
		//Chuyển user từ entity -> dto
		dto := []*dtoUser.User{}
		for i := 0; i < len(user); i++ {
			dto = append(dto, mapping.FromEntityToDtoUser(user[i]))
		}
		u_response.WriteSuccess(c, dto)
	}
}

func (h *HandleUser) HandleGetUserById() func(c *gin.Context) {
	return func(c *gin.Context) {
		//uid := shared.DecodeFromBase58(id)
		id := u_ctx.GetRequestContext(c.Request.Context()).Sub
		uid := u_uid.DecodeFromBase58(id)
		user, resp := h.UsecaseUser.UsecaseGetUserById(c, int(uid.LocalID))
		if resp != nil {
			u_response.WriteError(c, h.Log, resp)
			return
		}
		dto := &dtoUser.User{}
		dto = mapping.FromEntityToDtoUser(user)
		u_response.WriteSuccess(c, dto)
	}
}
