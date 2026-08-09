package handler

import (
	"errors"
	dtoAuth "example/internal/auth-service/dto/request"
	mapping "example/internal/auth-service/mapping"
	utils "example/internal/shared/utils"
	"time"

	w_easyjson "github.com/DVV-15324/witches/pkg/core/easyjson"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandle) HandleRegister() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data dtoAuth.Register

		// Bind JSON first
		if err := w_easyjson.BindJSON(c, &data); err != nil {
			w_resp.WriteErrorWithLog(c, h.Log, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()), utils.ReqKey)
			return
		}

		// Then validate
		err_valid := data.Validate()
		if err_valid != nil {
			resq := w_resp.NewAppError(401, err_valid, time.Now())
			w_resp.WriteErrorWithLog(c, h.Log, resq, utils.ReqKey)
			return // Added missing return
		}

		// Fixed: Use the data variable, not an empty struct
		entity := mapping.FromRegisterToModelAuth(&data)
		nameUser := data.Name

		// Register
		resq := h.UsecaseAuth.Register(c.Request.Context(), entity, nameUser)
		if resq != nil {
			w_resp.WriteErrorWithLog(c, h.Log, resq, utils.ReqKey)
			return
		}
		w_resp.WriteSuccess(c, "Dang ki thanh cong")
	}
}
