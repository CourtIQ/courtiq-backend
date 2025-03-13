package validation

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	internalErrors "github.com/CourtIQ/courtiq-backend/matchup-service/internal/errors"
)

// FormatValidator validates matchup format-related inputs
type FormatValidator struct{}

// NewFormatValidator creates a new format validator
func NewFormatValidator() *FormatValidator {
	return &FormatValidator{}
}

// Validate implements the Validator interface
func (v *FormatValidator) Validate(ctx context.Context, input interface{}) error {
	switch typedInput := input.(type) {
	case model.MatchUpFormatInput:
		return v.ValidateMatchUpFormatInput(ctx, typedInput)
	case *model.MatchUpFormatInput:
		if typedInput == nil {
			return internalErrors.NewRequiredFieldError("matchUpFormat")
		}
		return v.ValidateMatchUpFormatInput(ctx, *typedInput)
	default:
		return fmt.Errorf("unsupported input type for FormatValidator: %T", input)
	}
}

// ValidateMatchUpFormatInput validates a matchup format input
func (v *FormatValidator) ValidateMatchUpFormatInput(ctx context.Context, input model.MatchUpFormatInput) error {
	// Validate setFormat
	// Check if SetFormat is a pointer type in your model definition
	if err := v.ValidateSetFormat(input.SetFormat); err != nil {
		return err
	}

	// Validate finalSetFormat if numberOfSets is not 1
	if input.NumberOfSets != 1 {
		if input.FinalSetFormat != nil {
			if err := v.ValidateSetFormat(input.FinalSetFormat); err != nil {
				return err
			}
		}
	} else if input.FinalSetFormat != nil {
		return internalErrors.NewInvalidMatchFormatError("finalSetFormat should be nil when numberOfSets is 1")
	}

	return nil
}

// ValidateSetFormat validates a set format input
// Changed to accept a pointer to handle both pointer and non-pointer cases
func (v *FormatValidator) ValidateSetFormat(format *model.SetFormatInput) error {
	if format == nil {
		return internalErrors.NewRequiredFieldError("setFormat")
	}

	// Check if mustWinByTwo and tiebreakFormat are configured correctly
	if format.MustWinByTwo {
		if format.TiebreakFormat != nil {
			return internalErrors.NewInvalidMatchFormatError("tiebreakFormat should be nil when mustWinByTwo is true")
		}
	} else {
		// When mustWinByTwo is false, tiebreakFormat must be specified
		if format.TiebreakFormat == nil {
			return internalErrors.NewInvalidMatchFormatError("tiebreakFormat is required when mustWinByTwo is false")
		}

		// Validate the tiebreakFormat
		if err := v.ValidateTiebreakFormat(format.TiebreakFormat, int(format.NumberOfGames)); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTiebreakFormat validates a tiebreak format input
// Changed to accept a pointer to handle both pointer and non-pointer cases
func (v *FormatValidator) ValidateTiebreakFormat(format *model.TiebreakFormatInput, numberOfGames int) error {
	if format == nil {
		return internalErrors.NewRequiredFieldError("tiebreakFormat")
	}

	// Handle the case where TiebreakAt is a pointer
	tiebreakAt := format.TiebreakAt

	// Validate that tiebreakAt is >= numberOfGames
	if tiebreakAt < numberOfGames {
		return internalErrors.NewInvalidMatchFormatError(
			fmt.Sprintf("tiebreakAt (%d) must be greater than or equal to numberOfGames (%d)",
				tiebreakAt, numberOfGames))
	}

	return nil
}
