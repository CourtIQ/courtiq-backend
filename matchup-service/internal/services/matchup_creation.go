package services

import (
	"context"
	"errors"
	"fmt"
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

	ownerId, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 1) Check the top-level input (InitiateMatchUpInput!)
	if input == nil {
		return nil, errors.New("InitiateMatchUpInput cannot be nil")
	}

	now := time.Now()

	// 2) Construct the core MatchUp doc
	mu := &model.MatchUp{
		ID:             primitive.NewObjectID(),
		Owner:          ownerId,
		MatchUpTracker: input.MatchUpTracker,
		MatchUpType:    input.MatchUpType,
		MatchUpStatus:  model.MatchUpStatusScheduled,
		InitialServer:  input.InitialServer,
		TrackingStyle:  *input.TrackingStyle,
		CreatedAt:      now,
		LastUpdated:    now,
	}

	// Visibility defaults to PRIVATE if not provided
	if input.Visibility == nil {
		mu.Visibility = model.VisibilityPrivate
	} else {
		mu.Visibility = *input.Visibility
	}

	// 3) Validate required fields within MatchUpFormat
	if input.MatchUpFormat == nil {
		return nil, errors.New("matchUpFormat is required (cannot be null)")
	}
	if input.MatchUpFormat.SetFormat == nil {
		return nil, errors.New("matchUpFormat.setFormat is required (cannot be null)")
	}

	// 4) Build the SetFormat (required)
	setFormatInput := input.MatchUpFormat.SetFormat
	var tbf *model.TiebreakFormat
	// TiebreakFormat can be nil if no tiebreak is used
	if setFormatInput.TiebreakFormat != nil {
		tbf = &model.TiebreakFormat{
			Points:       setFormatInput.TiebreakFormat.Points,
			MustWinByTwo: setFormatInput.TiebreakFormat.MustWinByTwo,
			TiebreakAt:   setFormatInput.TiebreakFormat.TiebreakAt, // <--- NEW
		}
	}

	mainSetFormat := &model.SetFormat{
		NumberOfGames:  setFormatInput.NumberOfGames,
		DeuceType:      setFormatInput.DeuceType,
		MustWinByTwo:   setFormatInput.MustWinByTwo,
		TiebreakFormat: tbf,
		// Notice: we no longer set TiebreakAt directly on SetFormat
	}

	// 5) Build the finalSetFormat if provided (optional)
	var finalSetFormat *model.SetFormat
	if input.MatchUpFormat.FinalSetFormat != nil {
		fsf := input.MatchUpFormat.FinalSetFormat

		// Tiebreak for final set can also be nil
		var finalTbf *model.TiebreakFormat
		if fsf.TiebreakFormat != nil {
			finalTbf = &model.TiebreakFormat{
				Points:       fsf.TiebreakFormat.Points,
				MustWinByTwo: fsf.TiebreakFormat.MustWinByTwo,
				TiebreakAt:   fsf.TiebreakFormat.TiebreakAt, // <--- NEW
			}
		}
		finalSetFormat = &model.SetFormat{
			NumberOfGames:  fsf.NumberOfGames,
			DeuceType:      fsf.DeuceType,
			MustWinByTwo:   fsf.MustWinByTwo,
			TiebreakFormat: finalTbf,
			// Again, TiebreakAt is inside finalTbf, not here
		}
	}

	// 6) Build MatchUpFormat
	mu.MatchUpFormat = &model.MatchUpFormat{
		NumberOfSets:   input.MatchUpFormat.NumberOfSets,
		SetFormat:      mainSetFormat,
		FinalSetFormat: finalSetFormat,
	}

	fmt.Println("Participants: ", len(input.Participants))
	// 7) Participants (id must always be provided)
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
			ID:          partInput.ID.Hex(),
			DisplayName: partInput.DisplayedName,
			TeamSide:    partInput.TeamSide,
			IsGuest:     &isGuest,
		})
	}
	mu.Participants = participants

	return mu, nil
}
