package response

import (
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	"github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AppResponse struct {
	Status    int         `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type PaginationResponseWrapper struct {
	Status     int                       `json:"status"`
	Data       interface{}               `json:"data,omitempty"`
	Pagination *utils.PaginationResponse `json:"pagination,omitempty"`
	Message    string                    `json:"message,omitempty"`
	Timestamp  time.Time                 `json:"timestamp"`
}

// WriteSuccess gửi response thành công => không log
func WriteSuccess(c *gin.Context, data interface{}) {
	r := AppResponse{
		Status:    http.StatusOK,
		Data:      data,
		Message:   "Success",
		Timestamp: time.Now(),
	}
	c.JSON(r.Status, r)
}

// WriteSuccessWithPagination - Response thành công có kèm pagination  => không log
func WriteSuccessWithPagination(c *gin.Context, data interface{}, pagination *utils.PaginationResponse) {
	r := PaginationResponseWrapper{
		Status:     http.StatusOK,
		Data:       data,
		Pagination: pagination,
		Message:    "Success",
		Timestamp:  time.Now(),
	}
	c.JSON(r.Status, r)
}

// WriteSuccessWithLog gửi response thành công => log
func WriteSuccessWithLog(c *gin.Context, log *logger.EntityLogger, data interface{}, keyReq string) {
	r := AppResponse{
		Status:    http.StatusOK,
		Data:      data,
		Message:   "Success",
		Timestamp: time.Now(),
	}
	tid := utils.GetTid(c, keyReq)
	sub := utils.GetSub(c, keyReq)
	log.InfoWithFields("API response success",
		zap.String("trace_id(tid)", tid),
		zap.String("subject(sub)", sub),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", r.Status),
		zap.Time("timestamp", time.Now()),
	)
	c.JSON(r.Status, r)
}

// WriteSuccessWithPaginationAndLog - Response thành công có pagination => log
func WriteSuccessWithPaginationAndLog(c *gin.Context, log *logger.EntityLogger, data interface{}, pagination *utils.PaginationResponse, keyReq string) {
	r := PaginationResponseWrapper{
		Status:     http.StatusOK,
		Data:       data,
		Pagination: pagination,
		Message:    "Success",
		Timestamp:  time.Now(),
	}
	tid := utils.GetTid(c, keyReq)
	sub := utils.GetSub(c, keyReq)

	log.InfoWithFields("API response success with pagination",
		zap.String("trace_id(tid)", tid),
		zap.String("subject(sub)", sub),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", r.Status),
		zap.Time("timestamp", time.Now()),
	)
	c.JSON(r.Status, r)
}

// WriteError gửi response lỗi => không log
func WriteError(c *gin.Context, re *AppError) {
	r := AppResponse{
		Status:    re.Status,
		Message:   re.Error.Error(),
		Timestamp: re.TimeStamp,
	}
	c.JSON(r.Status, r)
}

// WriteError gửi response lỗi => log
func WriteErrorWithLog(c *gin.Context, log *logger.EntityLogger, re *AppError, keyReq string) {
	r := AppResponse{
		Status:    re.Status,
		Message:   re.Error.Error(),
		Timestamp: re.TimeStamp,
	}
	tid := utils.GetTid(c, keyReq)
	sub := utils.GetSub(c, keyReq)

	// Ghi log lỗi với đầy đủ context
	log.ErrorWithFields("API error",
		zap.String("trace_id(tid)", tid),
		zap.String("subject(sub)", sub),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", re.Status),
		zap.Error(re.Error),
		zap.Time("timestamp", re.TimeStamp),
	)

	c.JSON(r.Status, r)
}
