package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardErrors(t *testing.T) {
	// Test that standard errors are defined
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrInvalidInput)
	assert.NotNil(t, ErrUnauthorized)
	assert.NotNil(t, ErrForbidden)
	assert.NotNil(t, ErrDatabaseOperation)
	assert.NotNil(t, ErrInternalServer)
	assert.NotNil(t, ErrAlreadyExists)
	assert.NotNil(t, ErrConflict)
}

func TestValidationError(t *testing.T) {
	// Test creating and formatting validation errors
	err := NewValidationError("email", "invalid format")
	
	// Check error message
	assert.Contains(t, err.Error(), "email")
	assert.Contains(t, err.Error(), "invalid format")
	
	// Check type assertions
	assert.True(t, IsValidationError(err))
	assert.False(t, IsValidationError(errors.New("generic error")))
	
	// Test casting to ValidationError
	valErr, ok := err.(ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "email", valErr.Field)
	assert.Equal(t, "invalid format", valErr.Message)
}

func TestWrapError(t *testing.T) {
	// Test wrapping errors with context
	original := errors.New("original error")
	wrapped := WrapError(original, "context message")
	
	assert.Contains(t, wrapped.Error(), "context message")
	assert.Contains(t, wrapped.Error(), "original error")
	assert.True(t, errors.Is(wrapped, original))
}

func TestErrorCheckers(t *testing.T) {
	// Test IsNotFoundError
	wrappedNotFound := WrapError(ErrNotFound, "user not found")
	assert.True(t, IsNotFoundError(ErrNotFound))
	assert.True(t, IsNotFoundError(wrappedNotFound))
	assert.False(t, IsNotFoundError(errors.New("random error")))
	
	// Test IsForbiddenError
	wrappedForbidden := WrapError(ErrForbidden, "access denied")
	assert.True(t, IsForbiddenError(ErrForbidden))
	assert.True(t, IsForbiddenError(wrappedForbidden))
	assert.False(t, IsForbiddenError(errors.New("random error")))
	
	// Test remaining error checkers
	assert.True(t, IsUnauthorizedError(ErrUnauthorized))
	assert.True(t, IsInvalidInputError(ErrInvalidInput))
	assert.True(t, IsAlreadyExistsError(ErrAlreadyExists))
	assert.True(t, IsConflictError(ErrConflict))
	assert.True(t, IsDatabaseOperationError(ErrDatabaseOperation))
}