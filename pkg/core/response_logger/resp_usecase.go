package response_logger

import (
	"fmt"
	"time"
)

type ErrorResponse struct {
	Status    int
	Error     error
	TimeStamp time.Time
}

func (e *ErrorResponse) PrintErrorResponse() {
	fmt.Printf("[Reponse] Status: %v\n", e.Status)
	fmt.Printf("[Reponse] Error: %v\n", e.Error.Error())
	fmt.Printf("[Reponse] TimeStamp: %v\n", e.TimeStamp)
}

func NewErrorResponse(status int, Error error, timeStamp time.Time) *ErrorResponse {
	return &ErrorResponse{
		Status:    status,
		Error:     Error,
		TimeStamp: timeStamp,
	}
}
