package service

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
	ownerId, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if input == nil {
		return nil, errors.New("InitiateMatchUpInput cannot be nil")
	}

	now := time.Now()

	// Build the core MatchUp doc with defaults
	mu := &model.MatchUp{
		ID:             primitive.NewObjectID(),
		Owner:          ownerId,
		MatchUpTracker: input.MatchUpTracker,
		MatchUpType:    input.MatchUpType,
		MatchUpStatus:  model.MatchUpStatusScheduled,
		InitialServer:  input.InitialServer,
		CreatedAt:      now,
		LastUpdated:    now,
	}

	// If visibility is not provided, default to PRIVATE
	if input.Visibility == nil {
		mu.Visibility = model.VisibilityPrivate
	} else {
		mu.Visibility = *input.Visibility
	}

	// Convert the MatchUpFormat (assuming input.MatchUpFormat is non-nil)
	mu.MatchUpFormat = &model.MatchUpFormat{
		NumberOfSets: input.MatchUpFormat.NumberOfSets,
		SetFormat: &model.SetFormat{
			NumberOfGames: input.MatchUpFormat.SetFormat.NumberOfGames,
			DeuceType:     input.MatchUpFormat.SetFormat.DeuceType,
			MustWinByTwo:  input.MatchUpFormat.SetFormat.MustWinByTwo,
			TiebreakFormat: &model.TiebreakFormat{
				Points:       input.MatchUpFormat.SetFormat.TiebreakFormat.Points,
				MustWinByTwo: input.MatchUpFormat.SetFormat.TiebreakFormat.MustWinByTwo,
			},
			TiebreakAt: input.MatchUpFormat.SetFormat.TiebreakAt,
		},
	}
	// If there's a finalSetFormat, apply it
	if input.MatchUpFormat.FinalSetFormat != nil {
		mu.MatchUpFormat.FinalSetFormat = &model.SetFormat{
			NumberOfGames: input.MatchUpFormat.FinalSetFormat.NumberOfGames,
			DeuceType:     input.MatchUpFormat.FinalSetFormat.DeuceType,
			MustWinByTwo:  input.MatchUpFormat.FinalSetFormat.MustWinByTwo,
			TiebreakFormat: &model.TiebreakFormat{
				Points:       input.MatchUpFormat.FinalSetFormat.TiebreakFormat.Points,
				MustWinByTwo: input.MatchUpFormat.FinalSetFormat.TiebreakFormat.MustWinByTwo,
			},
			TiebreakAt: input.MatchUpFormat.FinalSetFormat.TiebreakAt,
		}
	}

	// Convert participants (simple approach—no guest logic).
	var participants []*model.Participant
	for _, partInput := range input.Participants {
		if partInput == nil {
			return nil, errors.New("participant input cannot be nil")
		}

		var participantID primitive.ObjectID
		if partInput.ID.IsZero() {
			participantID = primitive.NewObjectID()
		} else {
			participantID = partInput.ID
		}

		participants = append(participants, &model.Participant{
			ID:          participantID.Hex(),
			DisplayName: partInput.DisplayedName,
			TeamSide:    partInput.TeamSide,
		})
	}
	mu.Participants = participants

	return mu, nil
}
