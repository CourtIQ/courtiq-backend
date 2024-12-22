package service

import (
	"context"
	"github.com/CourtIQ/courtiq-backend/shared/services/utils"
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpServiceIntf defines the interface for matchup operations
type MatchUpServiceIntf interface {
	GetMatchUp(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error)
	CreateMatchUp(ctx context.Context, matchUpFormatInput model.MatchUpFormatInput, matchUpType model.MatchUpType, participantsMapInput model.ParticipantsMapInput) (*model.MatchUp, error)
	UpdateMatchUpStatus(ctx context.Context, status model.MatchUpStatus, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	AddPointToMatchUp(ctx context.Context, pointInput model.PointInput, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	UndoLastShotFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error)
	UndoLastPointFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error)
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
	participantsMapInput model.ParticipantsMapInput) (*model.MatchUp, error) {

	//Converting TiebreakFormatInput -> TiebreakFormat
	var tiebreakFormat *model.TiebreakFormat
	if matchUpFormatInput.SetFormat.TiebreakFormat != nil {
		tiebreakFormat = &model.TiebreakFormat{
			Points:       matchUpFormatInput.SetFormat.TiebreakFormat.Points,
			MustWinByTwo: matchUpFormatInput.SetFormat.TiebreakFormat.MustWinByTwo,
		}
	}

	// Converting SetFormatInput -> SetFormat
	setFormat := model.SetFormat{
		NumberOfGames:  matchUpFormatInput.SetFormat.NumberOfGames,
		DeuceType:      matchUpFormatInput.SetFormat.DeuceType,
		MustWinByTwo:   matchUpFormatInput.SetFormat.MustWinByTwo,
		TiebreakFormat: tiebreakFormat,
		TiebreakAt:     matchUpFormatInput.SetFormat.TiebreakAt,
	}

	// Converting SetFormatInput -> SetFormat of FinalSetFormat
	var finalSetFormat *model.SetFormat
	//if matchUpFormatInput.FinalSetFormat != nil {
	//	finalSetFormat = &model.SetFormat{
	//		NumberOfGames:  matchUpFormatInput.FinalSetFormat.NumberOfGames,
	//		DeuceType:      matchUpFormatInput.FinalSetFormat.DeuceType,
	//		MustWinByTwo:   matchUpFormatInput.FinalSetFormat.MustWinByTwo,
	//		TiebreakFormat: tiebreakFormat,
	//		TiebreakAt:     matchUpFormatInput.FinalSetFormat.TiebreakAt,
	//	}
	//}

	// Converting MatchUpFormatInput -> MatchUpFormat
	matchUpFormat := model.MatchUpFormat{
		Tracker:        matchUpFormatInput.Tracker,
		NumberOfSets:   matchUpFormatInput.NumberOfSets,
		SetFormat:      &setFormat,
		FinalSetFormat: finalSetFormat,
		InitialServer:  "",
	}

	// Setting initial score object
	score := model.Score{
		A: &model.SideScore{
			Player:               participantsMapInput.A[0],
			CurrentPointScore:    model.GameScoreLove,
			CurrentGameScore:     0,
			CurrentSetScore:      0,
			CurrentTiebreakScore: nil,
		},
		B: &model.SideScore{
			Player:               participantsMapInput.B[0],
			CurrentPointScore:    model.GameScoreLove,
			CurrentGameScore:     0,
			CurrentSetScore:      0,
			CurrentTiebreakScore: nil,
		},
	}
	// Converting ParticipantsMapInput -> ParticipantsMap
	participantsMap := model.ParticipantsMap{
		A: utils.ConvertListOfObjToListOfPtr(&participantsMapInput.A),
		B: utils.ConvertListOfObjToListOfPtr(&participantsMapInput.B),
	}

	timeNow := time.Now()

	matchUp := &model.MatchUp{
		ID:              primitive.NewObjectID(),
		MatchUpFormat:   &matchUpFormat,
		MatchUpStatus:   model.MatchUpStatusRequested,
		MatchUpType:     matchUpType,
		ParticipantsMap: &participantsMap,
		CurrentScore:    &score,
		CurrentServer:   participantsMapInput.A[0], // TODO assuming 0
		PointsSequence:  make([]*model.Point, 0),
		StartTime:       timeNow,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
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

func (s *MatchUpService) AddPointToMatchUp(ctx context.Context, pointInput model.PointInput, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if matchUp == nil || err != nil {
		return nil, err
	}

	var shots []*model.Shot
	// Converting ShotsInput -> Shots
	if pointInput.ShotsInput != nil {
		for _, shotInput := range pointInput.ShotsInput {
			shot := model.Shot{
				PlayerID:          shotInput.PlayerID,
				ShotType:          shotInput.ShotType,
				ServeStyle:        shotInput.ServeStyle,
				GroundStrokeType:  shotInput.GroundStrokeType,
				GroundStrokeStyle: shotInput.GroundStrokeStyle,
				PlayedAt:          shotInput.PlayedAt,
			}
			shots = append(shots, &shot)
		}
	}
	// Converting PointInput -> Point
	point := model.Point{
		SetIndex:            pointInput.SetIndex,
		GameIndexWithinSet:  pointInput.GameIndexWithinSet,
		IsTiebreak:          pointInput.IsTiebreak,
		TiebreakPointNumber: pointInput.TiebreakPointNumber,
		PointWinner:         pointInput.PointWinner,
		PointServer:         pointInput.PointServer,
		PointWinReason:      pointInput.PointWinReason,
		PlayingSide:         pointInput.PlayingSide,
		CourtSide:           pointInput.CourtSide,
		Shots:               shots,
		IsBreakPoint:        pointInput.IsBreakPoint,
		IsGamePoint:         pointInput.IsGamePoint,
		IsSetPoint:          pointInput.IsSetPoint,
		IsMatchPoint:        pointInput.IsMatchPoint,
		PlayedAt:            pointInput.PlayedAt,
	}

	matchUp.PointsSequence = append(matchUp.PointsSequence, &point)

	err = s.matchUpRepo.Update(ctx, matchUp)

	return matchUp, err
}
func (s *MatchUpService) UndoLastShotFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if matchUp != nil && err == nil {
		// TODO finish business logic
		err = s.matchUpRepo.Update(ctx, matchUp)
		return matchUp, err
	}
	return nil, err
}

func (s *MatchUpService) UndoLastPointFromMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if matchUp != nil && err == nil {
		if len(matchUp.PointsSequence) > 0 {
			matchUp.PointsSequence = matchUp.PointsSequence[:len(matchUp.PointsSequence)-1]
			err = s.matchUpRepo.Update(ctx, matchUp)
		}
		return matchUp, err
	}
	return nil, err
}

func (s *MatchUpService) DeleteMatchUp(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUp, error) {
	matchUp, err := s.matchUpRepo.FindByID(ctx, matchUpID)
	if matchUp != nil && err == nil {
		err = s.matchUpRepo.Delete(ctx, matchUpID)
	}
	return matchUp, err
}
