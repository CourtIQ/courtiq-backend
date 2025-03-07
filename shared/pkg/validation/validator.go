package validation

import (
	"regexp"
	"strings"

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
)

var (
	// EmailRegex validates email format
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// UsernameRegex validates username format (alphanumeric, underscore, 3-30 chars)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	
	// NameRegex validates name format (letters, spaces, hyphens, apostrophes)
	NameRegex = regexp.MustCompile(`^[a-zA-Z' -]{2,50}$`)
)

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return sharedErrors.NewValidationError("email", "email is required")
	}
	
	if !EmailRegex.MatchString(email) {
		return sharedErrors.NewValidationError("email", "invalid email format")
	}
	
	return nil
}

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if username == "" {
		return sharedErrors.NewValidationError("username", "username is required")
	}
	
	if !UsernameRegex.MatchString(username) {
		return sharedErrors.NewValidationError("username", "username must be 3-30 characters and contain only letters, numbers, and underscores")
	}
	
	return nil
}

// ValidateName validates a name field
func ValidateName(name string, fieldName string) error {
	if name == "" {
		return sharedErrors.NewValidationError(fieldName, fieldName+" is required")
	}
	
	if !NameRegex.MatchString(name) {
		return sharedErrors.NewValidationError(fieldName, fieldName+" must be 2-50 characters and contain only letters, spaces, hyphens, and apostrophes")
	}
	
	return nil
}

// ValidateRequired validates that a string is not empty
func ValidateRequired(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return sharedErrors.NewValidationError(fieldName, fieldName+" is required")
	}
	
	return nil
}