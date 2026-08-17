package response

import (
	"time"
)

type AppError struct {
	Status    int
	AError    error
	TimeStamp time.Time
}

func (e *AppError) Error() string {
	if e.AError != nil {
		return e.AError.Error()
	}
	return "unknown error"
}
func NewAppError(status int, Error error, timeStamp time.Time) *AppError {
	return &AppError{
		Status:    status,
		AError:    Error,
		TimeStamp: timeStamp,
	}
}
