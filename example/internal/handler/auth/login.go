package auth

import (
	dtoAuth "example/internal/dto/auth/request"
	mapping "example/internal/mapping"

	u_response "github.com/DVV-15324/witches/pkg/core/response_logger"

	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandleAuth) HandleLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data dtoAuth.Login
		// xac thuc
		err_valid := data.Validate()
		if err_valid != nil {
			resq := u_response.NewErrorResponse(401, err_valid, time.Now())
			u_response.WriteError(c, h.Log, resq)
		}
		//lấy dữ liệu
		err := c.ShouldBindJSON(&data)
		if err != nil {
			resq := u_response.NewErrorResponse(401, err, time.Now())
			u_response.WriteError(c, h.Log, resq)
			return
		}
		entity := mapping.FromLoginToEntityAuth(&dtoAuth.Login{})
		claims, resq := h.UsecaseAuth.UsecaseLogin(c, entity)
		if resq != nil {
			u_response.WriteError(c, h.Log, resq)
			return
		}
		u_response.WriteSuccess(c, claims)

	}
}
