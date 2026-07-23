package response_logger

import (
	"time"
)

type ErrorResponse struct {
	Status    int
	Error     error
	TimeStamp time.Time
}

func NewErrorResponse(status int, Error error, timeStamp time.Time) *ErrorResponse {
	return &ErrorResponse{
		Status:    status,
		Error:     Error,
		TimeStamp: timeStamp,
	}
}
