package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	Err        error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeResourceNotFound  = "RESOURCE_NOT_FOUND"
	ErrCodeDuplicateResource = "DUPLICATE_RESOURCE"
	ErrCodeDatabase          = "DATABASE_ERROR"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
)

// New creates a new AppError
func New(statusCode int, code, message string, details any, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
		Err:        err,
	}
}

// NewInvalidInputError creates a new invalid input error
func NewInvalidInputError(message string, details any, err error) *AppError {
	return New(http.StatusBadRequest, ErrCodeInvalidInput, message, details, err)
}

// NewResourceNotFoundError creates a new resource not found error
func NewResourceNotFoundError(message string, details any, err error) *AppError {
	return New(http.StatusNotFound, ErrCodeResourceNotFound, message, details, err)
}

// NewDuplicateResourceError creates a new duplicate resource error
func NewDuplicateResourceError(message string, details any, err error) *AppError {
	return New(http.StatusConflict, ErrCodeDuplicateResource, message, details, err)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, ErrCodeDatabase, message, nil, err)
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, ErrCodeInternal, message, nil, err)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, err error) *AppError {
	return New(http.StatusUnauthorized, ErrCodeUnauthorized, message, nil, err)
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, err error) *AppError {
	return New(http.StatusForbidden, ErrCodeForbidden, message, nil, err)
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == http.StatusNotFound
	}
	return false
}

// GetStatusCode returns the status code from an error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}
