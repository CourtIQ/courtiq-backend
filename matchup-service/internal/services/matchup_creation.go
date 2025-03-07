package services

import (
	"context"
	"errors"
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewMatchUpFromInitiateInput builds a MatchUp document from InitiateMatchUpInput.
// It assigns defaults (like visibility=PRIVATE) if not provided and
// sets MatchUpStatus to SCHEDULED by default.
func NewMatchUpFromInitiateInput(
	ctx context.Context,
	input *model.InitiateMatchUpInput,
) (*model.MatchUp, error) {
	// 1) Validate the context and retrieve the owner's ID
	ownerId, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 2) Validate the input
	if input == nil {
		return nil, errors.New("InitiateMatchUpInput cannot be nil")
	}

	// 3) Initialize the MatchUp document with core fields
	now := time.Now()
	mu := &model.MatchUp{
		ID:             primitive.NewObjectID(),      // Generate a new unique ID for the MatchUp
		Owner:          ownerId,                      // Set the owner from the context
		MatchUpTracker: input.MatchUpTracker,         // Set the MatchUp tracker
		MatchUpType:    input.MatchUpType,            // Set the MatchUp type (e.g., SINGLES, DOUBLES)
		MatchUpStatus:  model.MatchUpStatusScheduled, // Default status is SCHEDULED
		InitialServer:  input.InitialServer,          // Set the initial server
		CreatedAt:      now,                          // Set the creation timestamp
		LastUpdated:    now,                          // Set the last updated timestamp
	}

	// 4) Handle optional fields with defaults
	// Set TrackingStyle (default to empty string if nil)
	if input.TrackingStyle == nil {
		mu.TrackingStyle = nil // Default to empty string
	} else {
		mu.TrackingStyle = input.TrackingStyle
	}

	// Set Visibility (default to PRIVATE if nil)
	if input.Visibility == nil {
		mu.Visibility = model.VisibilityPrivate // Default to PRIVATE
	} else {
		mu.Visibility = *input.Visibility
	}

	// 5) Validate and build MatchUpFormat
	if input.MatchUpFormat == nil {
		return nil, errors.New("matchUpFormat is required (cannot be null)")
	}
	if input.MatchUpFormat.SetFormat == nil {
		return nil, errors.New("matchUpFormat.setFormat is required (cannot be null)")
	}

	// Build the main SetFormat
	setFormatInput := input.MatchUpFormat.SetFormat
	var tbf *model.TiebreakFormat
	if setFormatInput.TiebreakFormat != nil {
		tbf = &model.TiebreakFormat{
			Points:       setFormatInput.TiebreakFormat.Points,
			MustWinByTwo: setFormatInput.TiebreakFormat.MustWinByTwo,
			TiebreakAt:   setFormatInput.TiebreakFormat.TiebreakAt,
		}
	}

	mainSetFormat := &model.SetFormat{
		NumberOfGames:  setFormatInput.NumberOfGames,
		DeuceType:      setFormatInput.DeuceType,
		MustWinByTwo:   setFormatInput.MustWinByTwo,
		TiebreakFormat: tbf,
	}

	// Build the finalSetFormat (optional)
	var finalSetFormat *model.SetFormat
	if input.MatchUpFormat.FinalSetFormat != nil {
		fsf := input.MatchUpFormat.FinalSetFormat
		var finalTbf *model.TiebreakFormat
		if fsf.TiebreakFormat != nil {
			finalTbf = &model.TiebreakFormat{
				Points:       fsf.TiebreakFormat.Points,
				MustWinByTwo: fsf.TiebreakFormat.MustWinByTwo,
				TiebreakAt:   fsf.TiebreakFormat.TiebreakAt,
			}
		}
		finalSetFormat = &model.SetFormat{
			NumberOfGames:  fsf.NumberOfGames,
			DeuceType:      fsf.DeuceType,
			MustWinByTwo:   fsf.MustWinByTwo,
			TiebreakFormat: finalTbf,
		}
	}

	// Set the MatchUpFormat
	mu.MatchUpFormat = &model.MatchUpFormat{
		NumberOfSets:   input.MatchUpFormat.NumberOfSets,
		SetFormat:      mainSetFormat,
		FinalSetFormat: finalSetFormat,
	}

	// 6) Validate and build Participants
	if len(input.Participants) != 2 && len(input.Participants) != 4 {
		return nil, errors.New("either 2 or 4 participants are required")
	}

	var participants []*model.Participant
	for _, partInput := range input.Participants {
		if partInput == nil {
			return nil, errors.New("participant input cannot be nil")
		}

		var isGuest bool
		// If ID is nil, assign a new one and mark as guest
		if partInput.ID == nil {
			newID := primitive.NewObjectID()
			partInput.ID = &newID
			isGuest = true
		} else {
			// If an ID is present, this is an existing user
			isGuest = false
		}

		participants = append(participants, &model.Participant{
			ID:          partInput.ID.Hex(), // Convert ObjectID to string
			DisplayName: partInput.DisplayedName,
			TeamSide:    partInput.TeamSide,
			IsGuest:     isGuest, // Set the guest flag
		})
	}
	mu.Participants = participants

	// 7) Set the current server to the initial server
	mu.CurrentServer = input.InitialServer

	return mu, nil
}
