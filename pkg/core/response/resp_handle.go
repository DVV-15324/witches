package response_logger

import (
	"github.com/DVV-15324/witches/pkg/core/response_logger/logger"
	"github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ResponseHandle struct {
	Status    int         `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// WriteSuccess gửi response thành công (có thể log nếu muốn)
func WriteSuccess(c *gin.Context, data interface{}) {
	r := ResponseHandle{
		Status:    http.StatusOK,
		Data:      data,
		Message:   "Success",
		Timestamp: time.Now(),
	}
	c.JSON(r.Status, r)
}

// WriteSuccessWithLog gửi response thành công + log (tuỳ chọn)
func WriteSuccessWithLog(c *gin.Context, log *logger.EntityLogger, data interface{}, keyReq string) {
	r := ResponseHandle{
		Status:    http.StatusOK,
		Data:      data,
		Message:   "Success",
		Timestamp: time.Now(),
	}
	tid := utils.GetTid(c, keyReq)
	sub := utils.GetSub(c, keyReq)
	// Log thành công (level Info)
	log.InfoWithFields("API response success",
		zap.String("trace_id(tid)", tid),
		zap.String("subject(sub)", sub),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", r.Status),
		zap.Time("timestamp", time.Now()),
	)
	c.JSON(r.Status, r)
}

// WriteError gửi response lỗi + nhớ limmit
func WriteError(c *gin.Context, log *logger.EntityLogger, re *ErrorResponse) {
	r := ResponseHandle{
		Status:    re.Status,
		Message:   re.Error.Error(),
		Timestamp: time.Now(),
	}
	c.JSON(r.Status, r)
}

// WriteError gửi response lỗi + ghi log error
func WriteErrorWithLog(c *gin.Context, log *logger.EntityLogger, re *ErrorResponse, keyReq string) {
	r := ResponseHandle{
		Status:    re.Status,
		Message:   re.Error.Error(),
		Timestamp: time.Now(),
	}
	tid := utils.GetTid(c, keyReq)
	sub := utils.GetSub(c, keyReq)

	// Ghi log lỗi với đầy đủ context
	log.ErrorWithFields("API error",
		zap.String("trace_id(tid)", tid),
		zap.String("subject(sub)", sub),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", re.Status),
		zap.Error(re.Error),
		zap.Time("timestamp", re.TimeStamp), // Thời gian bắt tại layer xảy ra lỗi
	)

	c.JSON(r.Status, r)
}
