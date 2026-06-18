package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status    int
	Error     string
	TimeStamp time.Time
}

func (e *ErrorResponse) PrintErrorResponse() {
	fmt.Printf("[Reponse] Status: %v\n", e.Status)
	fmt.Printf("[Reponse] Error: %v\n", e.Error)
	fmt.Printf("[Reponse] TimeStamp: %v\n", e.TimeStamp)
}

func NewErrorResponse(status int, Error error, timeStamp time.Time) *ErrorResponse {
	return &ErrorResponse{
		Status:    status,
		Error:     Error.Error(),
		TimeStamp: timeStamp,
	}
}

type ResponseHandle struct {
	Status    int         `json:"status"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func WriteSuccess(c *gin.Context, data interface{}) {
	var r ResponseHandle
	r.Status = 200
	r.Data = data
	r.Timestamp = time.Now()
	c.JSON(r.Status, r)
}

func WriteError(c *gin.Context, re *ErrorResponse) {
	var r ResponseHandle
	r.Status = re.Status
	r.Error = re.Error
	r.Timestamp = time.Now()
	c.JSON(r.Status, r)
}
