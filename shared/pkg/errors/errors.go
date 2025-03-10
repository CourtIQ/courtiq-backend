package errors

import (
	"errors"
	"fmt"
)

// Standard errors that can be used across all services
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when the input to a function is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when a user is not authorized to perform an action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a user is forbidden from accessing a resource
	ErrForbidden = errors.New("forbidden")

	// ErrDatabaseOperation is returned when a database operation fails
	ErrDatabaseOperation = errors.New("database operation failed")

	// ErrInternalServer is returned when an unexpected internal error occurs
	ErrInternalServer = errors.New("internal server error")

	// ErrAlreadyExists is returned when trying to create a resource that already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrConflict is returned when there's a conflict with the current state
	ErrConflict = errors.New("conflict with current state")
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the error message
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

// WrapError wraps an error with a context message
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsForbiddenError checks if an error is a forbidden error
func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsUnauthorizedError checks if an error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsInvalidInputError checks if an error is an invalid input error
func IsInvalidInputError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsAlreadyExistsError checks if an error is an already exists error
func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsConflictError checks if an error is a conflict error
func IsConflictError(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsDatabaseOperationError checks if an error is a database operation error
func IsDatabaseOperationError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}

// NewBadRequestError creates a new invalid input error with a message
func NewBadRequestError(message string) error {
	return WrapError(ErrInvalidInput, message)
}

// NewNotFoundError creates a new not found error with a message
func NewNotFoundError(message string) error {
	return WrapError(ErrNotFound, message)
}

// NewForbiddenError creates a new forbidden error with a message
func NewForbiddenError(message string) error {
	return WrapError(ErrForbidden, message)
}

// NewUnauthorizedError creates a new unauthorized error with a message
func NewUnauthorizedError(message string) error {
	return WrapError(ErrUnauthorized, message)
}

// NewInternalError creates a new internal server error with a message
func NewInternalError(message string) error {
	return WrapError(ErrInternalServer, message)
}

// NewAlreadyExistsError creates a new already exists error with a message
func NewAlreadyExistsError(message string) error {
	return WrapError(ErrAlreadyExists, message)
}

// NewConflictError creates a new conflict error with a message
func NewConflictError(message string) error {
	return WrapError(ErrConflict, message)
}