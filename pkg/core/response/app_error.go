package response

import (
	"time"
)

type AppError struct {
	Status    int
	Error     error
	TimeStamp time.Time
}

func NewErrorResponse(status int, Error error, timeStamp time.Time) *AppError {
	return &AppError{
		Status:    status,
		Error:     Error,
		TimeStamp: timeStamp,
	}
}
