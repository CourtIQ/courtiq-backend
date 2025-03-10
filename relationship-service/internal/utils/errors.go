package utils

import (
	"fmt"

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
)

// UIError is a custom error type for domain-specific errors with UI-friendly messages.
// This complements the shared error package by adding domain context.
type UIError struct {
	Code    string // e.g. "COACHSHIP_NOT_FOUND"
	Message string // e.g. "Coachship does not exist."
	err     error  // Underlying error, often from shared package
}

func (e UIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrapper interface
func (e UIError) Unwrap() error {
	return e.err
}

// NewUIError creates a new UI error with an underlying error
func NewUIError(code, message string, err error) UIError {
	return UIError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

// Domain-specific error wrappers that use shared errors underneath
func NewCoachshipNotFoundError() error {
	return UIError{
		Code:    "COACHSHIP_NOT_FOUND",
		Message: "Coaching relationship does not exist.",
		err:     sharedErrors.ErrNotFound,
	}
}

func NewFriendshipNotFoundError() error {
	return UIError{
		Code:    "FRIENDSHIP_NOT_FOUND",
		Message: "Friendship does not exist.",
		err:     sharedErrors.ErrNotFound,
	}
}

func NewRelationshipForbiddenError(message string) error {
	return UIError{
		Code:    "RELATIONSHIP_FORBIDDEN",
		Message: message,
		err:     sharedErrors.ErrForbidden,
	}
}

func NewSelfRequestError(action string) error {
	return UIError{
		Code:    "SELF_REQUEST_ERROR",
		Message: fmt.Sprintf("Cannot %s yourself", action),
		err:     sharedErrors.ErrInvalidInput,
	}
}

func NewRelationshipAlreadyExistsError(relationshipType string) error {
	return UIError{
		Code:    fmt.Sprintf("%s_ALREADY_EXISTS", relationshipType),
		Message: fmt.Sprintf("A %s relationship already exists or is pending", relationshipType),
		err:     sharedErrors.ErrAlreadyExists,
	}
}

func NewMissingIDError(relationshipType string) error {
	return UIError{
		Code:    fmt.Sprintf("MISSING_%s_ID", relationshipType),
		Message: fmt.Sprintf("Cannot update a %s without an ID.", relationshipType),
		err:     sharedErrors.ErrInvalidInput,
	}
}

// IsUIError checks if an error is a UI error and optionally with a specific code
func IsUIError(err error, code ...string) bool {
	uiErr, ok := err.(UIError)
	if !ok {
		return false
	}

	if len(code) == 0 {
		return true
	}

	for _, c := range code {
		if uiErr.Code == c {
			return true
		}
	}

	return false
}

// GetUIErrorCode extracts the code from a UIError
func GetUIErrorCode(err error) string {
	if uiErr, ok := err.(UIError); ok {
		return uiErr.Code
	}
	return ""
}

// GetUIErrorMessage extracts the message from a UIError
func GetUIErrorMessage(err error) string {
	if uiErr, ok := err.(UIError); ok {
		return uiErr.Message
	}
	return err.Error()
}

// StandardizeError converts any error to a UIError with appropriate code
func StandardizeError(err error) error {
	// If already a UIError, return as is
	if IsUIError(err) {
		return err
	}

	// Map standard errors to domain-specific UI errors
	switch {
	case sharedErrors.IsNotFoundError(err):
		return NewUIError("NOT_FOUND", "Resource not found", err)
	case sharedErrors.IsForbiddenError(err):
		return NewUIError("FORBIDDEN", "Permission denied", err)
	case sharedErrors.IsUnauthorizedError(err):
		return NewUIError("UNAUTHORIZED", "Authentication required", err)
	case sharedErrors.IsInvalidInputError(err):
		return NewUIError("INVALID_INPUT", "Invalid input provided", err)
	case sharedErrors.IsAlreadyExistsError(err):
		return NewUIError("ALREADY_EXISTS", "Resource already exists", err)
	case sharedErrors.IsConflictError(err):
		return NewUIError("CONFLICT", "Operation would conflict with current state", err)
	case sharedErrors.IsDatabaseOperationError(err):
		return NewUIError("DATABASE_ERROR", "Database operation failed", err)
	default:
		return NewUIError("INTERNAL_ERROR", "An unexpected error occurred", err)
	}
}