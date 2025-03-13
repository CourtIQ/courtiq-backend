package errors

import (
	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
)

// Validation error constants
const (
	ErrInvalidParticipantCount    = "invalid participant count for match type"
	ErrInvalidTeamDistribution    = "invalid team distribution for participants"
	ErrInitialServerNotParticipant = "initial server must be one of the participants"
	ErrInvalidMatchFormat         = "invalid match format configuration"
	ErrRequiredField              = "required field is missing"
)

// NewInvalidParticipantCountError returns an error for invalid participant count
func NewInvalidParticipantCountError(matchType string, expected, actual int) error {
	return sharedErrors.NewValidationError(
		"participants",
		ErrInvalidParticipantCount+": "+matchType+" requires "+string(expected)+" participants, but got "+string(actual),
	)
}

// NewInvalidTeamDistributionError returns an error for invalid team distribution
func NewInvalidTeamDistributionError(matchType string) error {
	return sharedErrors.NewValidationError(
		"teamDistribution",
		ErrInvalidTeamDistribution+": "+matchType+" requires equal distribution across teams",
	)
}

// NewInitialServerNotParticipantError returns an error when initial server is not a participant
func NewInitialServerNotParticipantError() error {
	return sharedErrors.NewValidationError(
		"initialServer",
		ErrInitialServerNotParticipant,
	)
}

// NewInvalidMatchFormatError returns an error for invalid match format
func NewInvalidMatchFormatError(reason string) error {
	return sharedErrors.NewValidationError(
		"matchUpFormat",
		ErrInvalidMatchFormat+": "+reason,
	)
}

// NewRequiredFieldError returns an error for a missing required field
func NewRequiredFieldError(fieldName string) error {
	return sharedErrors.NewValidationError(
		fieldName,
		ErrRequiredField,
	)
}