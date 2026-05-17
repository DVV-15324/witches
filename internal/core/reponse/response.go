package response

import (
	"fmt"
	"time"
)

type EntityResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	TimeStamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func NewEntityResponse(status int, message string, timeStamp time.Time, data interface{}) *EntityResponse {
	return &EntityResponse{
		Status:    status,
		Message:   message,
		TimeStamp: timeStamp,
		Data:      data,
	}
}

func (e *EntityResponse) PrintReponse() {
	fmt.Printf("[Reponse] Status: %v\n", e.Status)
	fmt.Printf("[Reponse] Message: %v\n", e.Message)
	fmt.Printf("[Reponse] TimeStamp: %v\n", e.TimeStamp)
	fmt.Printf("[Reponse] Data: %v\n", e.Data)
}
