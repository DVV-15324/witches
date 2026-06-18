package user

import (
	dtoUser "arc-golang/internal/dto/user/response"
	mapping "arc-golang/internal/mapping"
	c_ctx "arc-golang/internal/utils"
	c_response "arc-golang/internal/utils"
	c_uid "arc-golang/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *HandleUser) HandleGetAllUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		user, resp := h.UsecaseUser.UsecaseGetAllUser(c)
		if resp != nil {
			c_response.WriteError(c, resp)
			return
		}
		//Chuyển user từ entity -> dto
		dto := []*dtoUser.User{}
		for i := 0; i < len(user); i++ {
			dto = append(dto, mapping.ToDtoUser(user[i]))
		}
		c_response.WriteSuccess(c, dto)
	}
}

func (h *HandleUser) HandleGetUserById() func(c *gin.Context) {
	return func(c *gin.Context) {
		//uid := shared.DecodeFromBase58(id)
		id := c_ctx.GetRequestContext(c.Request.Context()).Sub
		uid := c_uid.DecodeFromBase58(id)
		user, resp := h.UsecaseUser.UsecaseGetUserById(c, int(uid.LocalID))
		if resp != nil {
			c_response.WriteError(c, resp)
			return
		}
		dto := &dtoUser.User{}
		dto = mapping.ToDtoUser(user)
		c_response.WriteSuccess(c, dto)
	}
}
