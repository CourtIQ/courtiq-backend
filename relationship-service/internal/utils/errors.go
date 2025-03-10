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
		Message: "Coachship does not exist.",
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