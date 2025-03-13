package factory

import (
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpFactory handles creation and conversion of MatchUp entities
type MatchUpFactory struct{}

// NewMatchUpFactory creates a new MatchUpFactory
func NewMatchUpFactory() *MatchUpFactory {
	return &MatchUpFactory{}
}

// CreateFromInput creates a new MatchUp from InitiateMatchUpInput
func (f *MatchUpFactory) CreateMatchUpFromInitiateMatchUpInput(ownerID primitive.ObjectID, input model.InitiateMatchUpInput) *model.MatchUp {
	now := time.Now()

	// Create base matchup
	matchUp := &model.MatchUp{
		ID:                 primitive.NewObjectID(),
		Owner:              ownerID,
		MatchUpTracker:     input.MatchUpTracker,
		MatchUpType:        input.MatchUpType,
		MatchUpStatus:      model.MatchUpStatusScheduled,
		InitialServer:      input.InitialServer,
		CurrentServer:      input.InitialServer,
		FirstShot:          nil,
		LastShot:           nil,
		Winner:             nil,
		Loser:              nil,
		ScheduledStartTime: nil,
		StartTime:          nil,
		EndTime:            nil,
		CreatedAt:          now,
		LastUpdated:        now,
	}

	// Set format from input or use default if not provided
	if input.MatchUpFormat != nil {
		matchUp.MatchUpFormat = f.convertMatchUpFormat(input.MatchUpFormat)
	}

	// Set participants from input
	matchUp.Participants = f.convertParticipants(input.Participants)

	// Initialize score based on format
	matchUp.CurrentScore = f.initializeScore()

	return matchUp
}

// convertMatchUpFormat maps format from input to domain model
func (f *MatchUpFactory) convertMatchUpFormat(formatInput *model.MatchUpFormatInput) *model.MatchUpFormat {
	format := &model.MatchUpFormat{
		NumberOfSets: formatInput.NumberOfSets,
		SetFormat: &model.SetFormat{
			NumberOfGames: formatInput.SetFormat.NumberOfGames,
			DeuceType:     formatInput.SetFormat.DeuceType,
			MustWinByTwo:  formatInput.SetFormat.MustWinByTwo,
		},
	}

	// Map tiebreak format if exists
	if formatInput.SetFormat.TiebreakFormat != nil {
		format.SetFormat.TiebreakFormat = &model.TiebreakFormat{
			Points:       formatInput.SetFormat.TiebreakFormat.Points,
			MustWinByTwo: formatInput.SetFormat.TiebreakFormat.MustWinByTwo,
			TiebreakAt:   &formatInput.SetFormat.TiebreakFormat.TiebreakAt,
		}
	}

	// Map final set format if exists
	if formatInput.FinalSetFormat != nil {
		format.FinalSetFormat = &model.SetFormat{
			NumberOfGames: formatInput.FinalSetFormat.NumberOfGames,
			DeuceType:     formatInput.FinalSetFormat.DeuceType,
			MustWinByTwo:  formatInput.FinalSetFormat.MustWinByTwo,
		}

		// Map tiebreak format for final set if exists
		if formatInput.FinalSetFormat.TiebreakFormat != nil {
			format.FinalSetFormat.TiebreakFormat = &model.TiebreakFormat{
				Points:       formatInput.FinalSetFormat.TiebreakFormat.Points,
				MustWinByTwo: formatInput.FinalSetFormat.TiebreakFormat.MustWinByTwo,
				TiebreakAt:   &formatInput.FinalSetFormat.TiebreakFormat.TiebreakAt,
			}
		}
	}

	return format
}

// initializeScore creates an initial score state based on the match format
func (f *MatchUpFactory) initializeScore() *model.MatchUpScore {
	// Create empty score with correct number of sets
	score := &model.MatchUpScore{
		Sets:            make([]*model.SetScore, 0),
		IsMatchComplete: false,
	}

	// Initialize first set
	firstSet := &model.SetScore{
		SetIndex:         1,
		IsCompleted:      false,
		IsTiebreakActive: false,
		Sides:            make([]*model.SideSetScore, 2),
	}

	// Initialize sides for Team A and Team B
	firstSet.Sides[0] = &model.SideSetScore{
		Side:           model.TeamSideTeamA,
		GamesWon:       0,
		InGameScore:    model.InGameScoreZero,
		TiebreakPoints: nil,
	}

	firstSet.Sides[1] = &model.SideSetScore{
		Side:           model.TeamSideTeamB,
		GamesWon:       0,
		InGameScore:    model.InGameScoreZero,
		TiebreakPoints: nil,
	}

	// Add first set to score
	score.Sets = append(score.Sets, firstSet)

	return score
}

// convertParticipants maps participants from input to domain model
func (f *MatchUpFactory) convertParticipants(participantsInput []*model.ParticipantInput) []*model.Participant {

	participants := make([]*model.Participant, len(participantsInput))

	for i, p := range participantsInput {
		participant := &model.Participant{
			DisplayName: p.DisplayedName,
			TeamSide:    p.TeamSide,
			IsGuest:     p.ID == nil, // Mark as guest if no ID provided
		}

		// If ID is provided, use it
		if p.ID != nil {
			participant.ID = *p.ID
		} else {
			newID := primitive.NewObjectID()
			participant.ID = newID
		}

		participants[i] = participant
	}

	return participants
}
