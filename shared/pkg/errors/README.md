# Shared Errors Package

This package provides standardized error types and error handling utilities for all microservices.

## Standard Error Types

```go
// Core error types
var (
    ErrNotFound          = errors.New("resource not found")
    ErrInvalidInput      = errors.New("invalid input")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrForbidden         = errors.New("forbidden")
    ErrDatabaseOperation = errors.New("database operation failed")
    ErrInternalServer    = errors.New("internal server error")
    ErrAlreadyExists     = errors.New("resource already exists")
    ErrConflict          = errors.New("conflict with current state")
)
```

## Validation Errors

The package includes a `ValidationError` type for handling field-level validation errors:

```go
// Create a validation error
err := NewValidationError("email", "invalid email format")

// Check if an error is a validation error
if IsValidationError(err) {
    // Handle validation error
}
```

## Error Helpers

The package provides utility functions for working with errors:

```go
// Wrap an error with additional context
err = WrapError(err, "failed to process user registration")

// Check error types
if IsNotFoundError(err) {
    // Handle not found error
}

if IsUnauthorizedError(err) {
    // Handle unauthorized error
}
```

## Extending with Service-Specific Errors

In your microservice, you can extend these error types with domain-specific errors:

```go
// In your service's errors package
package errors

import (
    "errors"
    sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
)

// Service-specific errors
var (
    ErrUserNotVerified = errors.New("user is not verified")
    ErrAccountLocked   = errors.New("account is locked")
)

// Check if an error is a user not verified error
func IsUserNotVerifiedError(err error) bool {
    return errors.Is(err, ErrUserNotVerified)
}

// Create domain-specific validation error
func NewUserValidationError(field, message string) error {
    return sharedErrors.NewValidationError(field, message)
}
```

## Error Translation for GraphQL

To convert errors to GraphQL-compatible formats:

```go
func ConvertToGraphQLError(err error) *graphql.Error {
    if sharedErrors.IsNotFoundError(err) {
        return &graphql.Error{
            Message: err.Error(),
            Extensions: map[string]interface{}{
                "code": "NOT_FOUND",
            },
        }
    }
    
    if sharedErrors.IsValidationError(err) {
        validErr := err.(sharedErrors.ValidationError)
        return &graphql.Error{
            Message: err.Error(),
            Extensions: map[string]interface{}{
                "code": "VALIDATION_ERROR",
                "field": validErr.Field,
            },
        }
    }
    
    // Handle other error types...
    
    return &graphql.Error{
        Message: "Internal server error",
        Extensions: map[string]interface{}{
            "code": "INTERNAL_SERVER_ERROR",
        },
    }
}
```