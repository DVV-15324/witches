package auth

import (
	dtoAuth "arc-golang/internal/dto/auth/request"
	mapping "arc-golang/internal/mapping"
	c_response "arc-golang/internal/utils"

	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandleAuth) HandleLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data dtoAuth.Login
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
		entity := mapping.FromLoginToEntityAuth(&dtoAuth.Login{})
		claims, resq := h.UsecaseAuth.UsecaseLogin(c, entity)
		if resq != nil {
			c_response.WriteError(c, resq)
			return
		}
		c_response.WriteSuccess(c, claims)

	}
}
