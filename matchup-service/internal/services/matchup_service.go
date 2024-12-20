package service

import (
	"context"
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpServiceIntf defines the interface for matchup operations
type MatchUpServiceIntf interface {
	GetMatchUp(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error)
	CreateMatchUp(ctx context.Context, matchUpFormatInput model.MatchUpFormatInput, matchUpType model.MatchUpType, participants []primitive.ObjectID) (*model.MatchUp, error)
	UpdateMatchUpStatus(ctx context.Context, status model.MatchUpStatus, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	AddPointToMatchUp(ctx context.Context, matchUpFormat model.MatchUpFormatInput, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	UndoShotFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	UndoPointFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	DeleteMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error)
}

// MatchUpService implements the MatchUpServiceIntf
type MatchUpService struct {
	matchUpRepo repository.MatchUpRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService(matchUpRepo repository.MatchUpRepository) MatchUpServiceIntf {
	return &MatchUpService{matchUpRepo: matchUpRepo}
}

func (s *MatchUpService) GetMatchUp(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error) {
	return s.matchUpRepo.FindByID(ctx, id)
}

func (s *MatchUpService) CreateMatchUp(
	ctx context.Context,
	matchUpFormatInput model.MatchUpFormatInput,
	matchUpType model.MatchUpType,
	participants []primitive.ObjectID) (*model.MatchUp, error) {

	// Converting TiebreakFormatInput -> TiebreakFormat
	//tiebreakFormat := model.TiebreakFormat{
	//	Points:       matchUpFormatInput.SetFormat.TiebreakFormat.Points,
	//	MustWinByTwo: matchUpFormatInput.SetFormat.TiebreakFormat.MustWinByTwo,
	//}

	// Converting SetFormatInput -> SetFormat
	setFormat := model.SetFormat{
		NumberOfGames: matchUpFormatInput.SetFormat.NumberOfGames,
		DeuceType:     matchUpFormatInput.SetFormat.DeuceType,
		MustWinByTwo:  matchUpFormatInput.SetFormat.MustWinByTwo,
		//TiebreakFormat: &tiebreakFormat,
		TiebreakFormat: nil,
		TiebreakAt:     matchUpFormatInput.SetFormat.TiebreakAt,
	}

	// Converting MatchUpFormatInput -> MatchUpFormat
	matchupFormat := model.MatchUpFormat{
		Tracker:        matchUpFormatInput.Tracker,
		NumberOfSets:   matchUpFormatInput.NumberOfSets,
		SetFormat:      &setFormat,
		FinalSetFormat: nil,
		InitialServer:  "",
	}

	// Setting initial score object
	score := model.Score{
		A: &model.SideScore{
			Player:               participants[0],
			CurrentPointScore:    model.GameScoreLove,
			CurrentGameScore:     0,
			CurrentSetScore:      0,
			CurrentTiebreakScore: nil,
		},
		B: &model.SideScore{
			Player:               participants[1],
			CurrentPointScore:    model.GameScoreLove,
			CurrentGameScore:     0,
			CurrentSetScore:      0,
			CurrentTiebreakScore: nil,
		},
	}

	timeNow := time.Now()

	matchUp := &model.MatchUp{
		ID:             primitive.NewObjectID(),
		MatchUpFormat:  &matchupFormat,
		MatchUpStatus:  model.MatchUpStatusRequested,
		MatchUpType:    matchUpType,
		ParticipantIds: participants,
		Participants: &model.ParticipantsMap{
			A: participants[0],
			B: participants[1],
		},
		CurrentScore:   &score,
		CurrentServer:  participants[0], // TODO assuming 0
		PointsSequence: nil,
		StartTime:      timeNow,
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
	}

	if err := s.matchUpRepo.Insert(ctx, matchUp); err != nil {
		return nil, err
	}

	return matchUp, nil
}

func (s *MatchUpService) UpdateMatchUpStatus(ctx context.Context, status model.MatchUpStatus, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if err == nil {
		matchUp.MatchUpStatus = status
		err = s.matchUpRepo.Update(ctx, matchUp)
	}
	return matchUp, err
}

func (s *MatchUpService) AddPointToMatchUp(ctx context.Context, matchUpFormatInput model.MatchUpFormatInput, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if err != nil {
		return nil, err
	}
	// Converting SetFormatInput -> SetFormat of SetFormat
	setFormat := model.SetFormat{
		NumberOfGames: matchUpFormatInput.SetFormat.NumberOfGames,
		DeuceType:     matchUpFormatInput.SetFormat.DeuceType,
		MustWinByTwo:  matchUpFormatInput.SetFormat.MustWinByTwo,
	}

	var tiebreakFormat *model.TiebreakFormat
	var finalSetFormat *model.SetFormat
	if matchUpFormatInput.FinalSetFormat != nil {
		// Converting TiebreakFormatInput -> TiebreakFormat
		tiebreakFormat = &model.TiebreakFormat{
			Points:       matchUpFormatInput.FinalSetFormat.TiebreakFormat.Points,
			MustWinByTwo: matchUpFormatInput.FinalSetFormat.TiebreakFormat.MustWinByTwo,
		}
		// Converting SetFormatInput -> SetFormat of FinalSetFormat
		finalSetFormat = &model.SetFormat{
			NumberOfGames:  matchUpFormatInput.FinalSetFormat.NumberOfGames,
			DeuceType:      matchUpFormatInput.FinalSetFormat.DeuceType,
			MustWinByTwo:   matchUpFormatInput.FinalSetFormat.MustWinByTwo,
			TiebreakFormat: tiebreakFormat,
			TiebreakAt:     matchUpFormatInput.FinalSetFormat.TiebreakAt,
		}
	}
	// Converting matchUpFormatInput -> matchUpFormat
	matchUpFormat := model.MatchUpFormat{
		ID:             primitive.NewObjectID(),
		Tracker:        matchUpFormatInput.Tracker,
		NumberOfSets:   matchUpFormatInput.NumberOfSets,
		SetFormat:      &setFormat,
		FinalSetFormat: finalSetFormat,
		InitialServer:  "", // TODO where does this come from?
	}
	matchUp.MatchUpFormat = &matchUpFormat
	err = s.matchUpRepo.Update(ctx, matchUp)

	return matchUp, err
}
func (s *MatchUpService) UndoShotFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if err == nil {
		// TODO finish business logic
		err = s.matchUpRepo.Update(ctx, matchUp)
		return matchUp, err
	}
	return nil, err
}
func (s *MatchUpService) UndoPointFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if err == nil {
		// TODO finish business logic
		err = s.matchUpRepo.Update(ctx, matchUp)
		return matchUp, err
	}
	return nil, err
}
func (s *MatchUpService) DeleteMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	err := s.matchUpRepo.Delete(ctx, matchUpID)

	return nil, err // TODO should this return matchup?
}
