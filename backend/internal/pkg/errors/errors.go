// Package errors provides domain error types and HTTP status mapping.
package errors

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeNotFound     Code = "NOT_FOUND"
	CodeValidation   Code = "VALIDATION_ERROR"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeConflict     Code = "CONFLICT"
	CodeInternal     Code = "INTERNAL_ERROR"
	CodeRateLimited  Code = "TOO_MANY_REQUESTS"
)

type DomainError struct {
	Code     Code
	Message  string
	Err      error
	Metadata map[string]interface{}
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error { return e.Err }

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewNotFound(entity string, id interface{}) *DomainError {
	return &DomainError{Code: CodeNotFound, Message: fmt.Sprintf("%s not found: %v", entity, id)}
}

func NewValidationFailed(errs []ValidationError) *DomainError {
	metadata := make(map[string]interface{})
	details := make([]map[string]string, len(errs))
	for i, e := range errs {
		details[i] = map[string]string{"field": e.Field, "message": e.Message}
	}
	metadata["details"] = details
	return &DomainError{Code: CodeValidation, Message: "Validation failed", Metadata: metadata}
}

func NewUnauthorized(msg string) *DomainError {
	return &DomainError{Code: CodeUnauthorized, Message: msg}
}

func NewForbidden(msg string) *DomainError {
	return &DomainError{Code: CodeForbidden, Message: msg}
}

func NewConflict(msg string) *DomainError {
	return &DomainError{Code: CodeConflict, Message: msg}
}

func NewInternal(err error) *DomainError {
	return &DomainError{Code: CodeInternal, Message: "Internal server error", Err: err}
}

func NewRateLimited() *DomainError {
	return &DomainError{Code: CodeRateLimited, Message: "Too many requests"}
}

func HTTPStatus(err error) int {
	var de *DomainError
	if !As(err, &de) {
		return 500
	}
	switch de.Code {
	case CodeNotFound:
		return 404
	case CodeValidation:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeConflict:
		return 409
	case CodeRateLimited:
		return 429
	default:
		return 500
	}
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
