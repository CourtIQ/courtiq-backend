package validation

import (
	"context"
)

// Validator defines a generic interface for validating inputs
type Validator interface {
	// Validate validates the input and returns an error if validation fails
	Validate(ctx context.Context, input interface{}) error
}

// ValidatorFunc is a function type that implements the Validator interface
type ValidatorFunc func(ctx context.Context, input interface{}) error

// Validate implements the Validator interface for ValidatorFunc
func (f ValidatorFunc) Validate(ctx context.Context, input interface{}) error {
	return f(ctx, input)
}

// CompositeValidator combines multiple validators into one
type CompositeValidator struct {
	validators []Validator
}

// NewCompositeValidator creates a new composite validator from multiple validators
func NewCompositeValidator(validators ...Validator) *CompositeValidator {
	return &CompositeValidator{
		validators: validators,
	}
}

// Validate runs all validators in sequence, returning on the first error
func (cv *CompositeValidator) Validate(ctx context.Context, input interface{}) error {
	for _, validator := range cv.validators {
		if err := validator.Validate(ctx, input); err != nil {
			return err
		}
	}
	return nil
}