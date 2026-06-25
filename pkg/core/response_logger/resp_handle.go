package response_logger

import (
	"core-v/pkg/core/response_logger/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func WriteSuccessWithLog(c *gin.Context, log *logger.EntityLogger, data interface{}) {
	r := ResponseHandle{
		Status:    http.StatusOK,
		Data:      data,
		Message:   "Success",
		Timestamp: time.Now(),
	}
	// Log thành công (level Info)
	log.InfoWithFields("API response success",
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", r.Status),
	)
	c.JSON(r.Status, r)
}

// WriteError gửi response lỗi + ghi log error
func WriteError(c *gin.Context, log *logger.EntityLogger, re *ErrorResponse) {
	r := ResponseHandle{
		Status:    re.Status,
		Message:   re.Error.Error(),
		Timestamp: time.Now(),
	}

	// Ghi log lỗi với đầy đủ context
	log.ErrorWithFields("API error",
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", re.Status),
		zap.Error(re.Error),
		zap.Time("timestamp", re.TimeStamp),
	)

	c.JSON(r.Status, r)
}

// WriteErrorSimple dùng nếu không cần log chi tiết
func WriteErrorSimple(c *gin.Context, log *logger.EntityLogger, status int, err error) {
	re := NewErrorResponse(status, err, time.Now())
	WriteError(c, log, re)
}
