// Package response contains HTTP response DTO types.
package response

import "time"

type ErrorDetail struct {
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}
