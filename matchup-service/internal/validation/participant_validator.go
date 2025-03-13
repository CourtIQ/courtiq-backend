package validation

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	internalErrors "github.com/CourtIQ/courtiq-backend/matchup-service/internal/errors"
	sharedValidation "github.com/CourtIQ/courtiq-backend/shared/pkg/validation"
)

// ParticipantValidator validates participant-related inputs
type ParticipantValidator struct{}

// NewParticipantValidator creates a new participant validator
func NewParticipantValidator() *ParticipantValidator {
	return &ParticipantValidator{}
}

// ValidateParticipant validates an individual participant
func (v *ParticipantValidator) ValidateParticipant(participant *model.ParticipantInput) error {
	// Validate displayed name
	if err := sharedValidation.ValidateName(participant.DisplayedName, "displayedName"); err != nil {
		return err
	}

	// Validate team side
	if !participant.TeamSide.IsValid() {
		return internalErrors.NewInvalidMatchFormatError("invalid team side value")
	}

	return nil
}

// ValidateParticipantTeamDistribution validates that participants are distributed correctly between teams
func (v *ParticipantValidator) ValidateParticipantTeamDistribution(matchType model.MatchUpType, participants []*model.ParticipantInput) error {
	// Count participants per team
	teamACounts := 0
	teamBCounts := 0

	for _, participant := range participants {
		// Validate individual participant
		if err := v.ValidateParticipant(participant); err != nil {
			return err
		}

		// Count team distribution
		switch participant.TeamSide {
		case model.TeamSideTeamA:
			teamACounts++
		case model.TeamSideTeamB:
			teamBCounts++
		}
	}

	// Check distribution based on match type
	switch matchType {
	case model.MatchUpTypeSingles:
		if teamACounts != 1 || teamBCounts != 1 {
			return internalErrors.NewInvalidTeamDistributionError("singles")
		}
	case model.MatchUpTypeDoubles:
		if teamACounts != 2 || teamBCounts != 2 {
			return internalErrors.NewInvalidTeamDistributionError("doubles")
		}
	}

	return nil
}

// Validate implements the Validator interface
func (v *ParticipantValidator) Validate(ctx context.Context, input interface{}) error {
	switch typedInput := input.(type) {
	case []*model.ParticipantInput:
		for _, participant := range typedInput {
			if err := v.ValidateParticipant(participant); err != nil {
				return err
			}
		}
		return nil
	case *model.ParticipantInput:
		return v.ValidateParticipant(typedInput)
	default:
		return fmt.Errorf("unsupported input type for ParticipantValidator: %T", input)
	}
}
