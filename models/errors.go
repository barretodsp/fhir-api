package models

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code       string `json:"code"`    // Machine-readable error code
	Message    string `json:"message"` // Human-readable message
	StatusCode int    `json:"-"`       // HTTP status code
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s (status: %d)", e.Code, e.Message, e.StatusCode)
}

func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrNotFound       = NewAppError("NOT_FOUND", "resource not found", http.StatusNotFound)
	ErrInvalidInput   = NewAppError("INVALID_INPUT", "invalid input data", http.StatusBadRequest)
	ErrUnauthorized   = NewAppError("UNAUTHORIZED", "unauthorized access", http.StatusUnauthorized)
	ErrForbidden      = NewAppError("FORBIDDEN", "access forbidden", http.StatusForbidden)
	ErrDatabase       = NewAppError("DATABASE_ERROR", "database operation failed", http.StatusInternalServerError)
	ErrInternalServer = NewAppError("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)
