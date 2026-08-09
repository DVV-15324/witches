package handler

import (
	"errors"
	dtoRefresh "example/internal/refresh-service/dto/request"
	"example/internal/shared/utils"
	"time"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
)

func (h *RefreshHandle) HandleRefreshToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtoRefresh.RefreshTokenRequest

		// Bind request
		if err := c.ShouldBindJSON(&req); err != nil {
			w_resp.WriteErrorWithLog(c, h.Log, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()), utils.ReqKey)
			return
		}

		// Get device ID from header
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = "unknown" // or handle as needed
		}

		// Fixed: Added missing deviceID parameter
		tokenResp, err := h.RefreshUsecase.Refresh(c.Request.Context(), req.RefreshToken, deviceID)
		if err != nil {
			w_resp.WriteErrorWithLog(c, h.Log, err, utils.ReqKey)
			return
		}

		w_resp.WriteSuccess(c, tokenResp)
	}
}
