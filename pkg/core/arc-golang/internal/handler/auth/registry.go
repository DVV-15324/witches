package auth

import (
	dtoAuth "arc-golang/internal/dto/auth/request"
	entityAuth "arc-golang/internal/entity/auth"
	mapping "arc-golang/internal/mapping"
	c_response "arc-golang/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandleAuth) HandleRegister() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data dtoAuth.Register
		// xac thuc
		err_valid := data.Validate()
		if err_valid != nil {
			resq := c_response.NewErrorResponse(401, err_valid, time.Now())
			c_response.WriteError(c, resq)
		}
		//lấy dữ liệu
		err := c.ShouldBindJSON(&data)
		if err != nil {
			resq := c_response.NewErrorResponse(401, err, time.Now())
			c_response.WriteError(c, resq)
			return
		}
		entity := &entityAuth.Auth{}

		//Chuyển dto sang entity
		entity = mapping.FromRegisterToEntityAuth(&dtoAuth.Register{})
		nameUser := data.Name

		// dang ki
		resq := h.UsecaseAuth.UsecaseRegister(c, entity, nameUser)
		if resq != nil {
			c_response.WriteError(c, resq)
			return
		}
		c_response.WriteSuccess(c, "Dang ki thanh cong")
	}
}
