package validation

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	internalErrors "github.com/CourtIQ/courtiq-backend/matchup-service/internal/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpValidator validates matchup-related inputs
type MatchUpValidator struct {
	participantValidator *ParticipantValidator
}

// NewMatchUpValidator creates a new matchup validator
func NewMatchUpValidator() *MatchUpValidator {
	return &MatchUpValidator{
		participantValidator: NewParticipantValidator(),
	}
}

// Validate implements the Validator interface
func (v *MatchUpValidator) Validate(ctx context.Context, input interface{}) error {
	switch typedInput := input.(type) {
	case model.InitiateMatchUpInput:
		return v.ValidateInitiateMatchUpInput(ctx, typedInput)
	default:
		return errors.ErrUnsupported
	}
}

// ValidateInitiateMatchUpInput validates the input for initiating a match
func (v *MatchUpValidator) ValidateInitiateMatchUpInput(ctx context.Context, input model.InitiateMatchUpInput) error {
	// Validate match type and participant count
	if err := v.validateMatchTypeAndParticipants(input.MatchUpType, input.Participants); err != nil {
		return err
	}

	// Validate that the initial server is one of the participants
	if err := v.validateInitialServer(input.InitialServer, input.Participants); err != nil {
		return err
	}

	// Validate participant team distribution
	if err := v.participantValidator.ValidateParticipantTeamDistribution(input.MatchUpType, input.Participants); err != nil {
		return err
	}

	return nil
}

// validateMatchTypeAndParticipants validates that the number of participants matches the match type
func (v *MatchUpValidator) validateMatchTypeAndParticipants(matchType model.MatchUpType, participants []*model.ParticipantInput) error {
	if len(participants) == 0 {
		return internalErrors.NewRequiredFieldError("participants")
	}

	switch matchType {
	case model.MatchUpTypeSingles:
		if len(participants) != 2 {
			return internalErrors.NewInvalidParticipantCountError("singles", 2, len(participants))
		}
	case model.MatchUpTypeDoubles:
		if len(participants) != 4 {
			return internalErrors.NewInvalidParticipantCountError("doubles", 4, len(participants))
		}
	default:
		return internalErrors.NewInvalidMatchFormatError("invalid match type")
	}

	return nil
}

// validateInitialServer validates that the initial server is one of the participants
func (v *MatchUpValidator) validateInitialServer(serverID primitive.ObjectID, participants []*model.ParticipantInput) error {
	for _, participant := range participants {
		if participant.ID != nil && participant.ID.Hex() == serverID.Hex() {
			return nil
		}
	}

	return internalErrors.NewInitialServerNotParticipantError()
}

// validateMatchUpFormat validates the match format settings
func (v *MatchUpValidator) validateMatchUpFormat(format model.MatchUpFormatInput) error {
	// Basic validation - more complex rules can be added
	if format.NumberOfSets <= 0 {
		return internalErrors.NewInvalidMatchFormatError("number of sets must be positive")
	}

	if format.SetFormat == nil {
		return internalErrors.NewInvalidMatchFormatError("set format is required")
	}

	// Additional format validation can be added here

	return nil
}
