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

	side := GetParticipantSide(matchUp, point.PointWinner)
	IncrementScore(matchUp.CurrentScore, side)
	if IsGameFinished(matchUp.CurrentScore) {
		*matchUp.CurrentGameIndexWithinSet++
	}

	matchUp.PointsSequence = append(matchUp.PointsSequence, &point)
	if ShouldChangeServer(matchUp) {
		matchUp.CurrentServer = *GetNextServer(matchUp)
	}

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

func GetNextServer(matchUp *model.MatchUp) *primitive.ObjectID {
	allParticipants := append(matchUp.ParticipantsMap.A, matchUp.ParticipantsMap.B...)
	for i := 0; i < len(allParticipants); i++ {
		if i == len(allParticipants)-1 && *allParticipants[i] == matchUp.CurrentServer {
			return allParticipants[0]
		} else if *allParticipants[i] == matchUp.CurrentServer {
			return allParticipants[i+1]
		}
	}
	return nil
}

func ShouldChangeServer(matchUp *model.MatchUp) bool {
	//TODO implement logic
	return false
}

func GetParticipantSide(matchUp *model.MatchUp, playerID primitive.ObjectID) string {
	for _, participant := range matchUp.ParticipantsMap.A {
		if *participant == playerID {
			return "A"
		}
	}
	return "B"
}

// IncrementScore increments the score for the given side according to Tennis score logic
func IncrementScore(score *model.Score, side string) {
	if score.A.CurrentPointScore == model.GameScoreForty && score.B.CurrentPointScore == model.GameScoreForty {
		if side == "A" {
			score.A.CurrentPointScore = model.GameScoreAdvantage
		} else {
			score.B.CurrentPointScore = model.GameScoreAdvantage
		}
		return
	}

	if side == "A" {
		if score.A.CurrentPointScore == model.GameScoreAdvantage {
			score.A.CurrentPointScore = model.GameScoreGame
		} else if score.B.CurrentPointScore == model.GameScoreAdvantage {
			score.A.CurrentPointScore = model.GameScoreForty
			score.B.CurrentPointScore = model.GameScoreForty
		} else {
			score.A.CurrentPointScore = incrementScore(score.A.CurrentPointScore)
		}
	} else {
		if score.B.CurrentPointScore == model.GameScoreAdvantage {
			score.B.CurrentPointScore = model.GameScoreGame
		} else if score.A.CurrentPointScore == model.GameScoreAdvantage {
			score.A.CurrentPointScore = model.GameScoreForty
			score.B.CurrentPointScore = model.GameScoreForty
		} else {
			score.B.CurrentPointScore = incrementScore(score.B.CurrentPointScore)
		}
	}
}

func incrementScore(gameScore model.GameScore) model.GameScore {
	switch gameScore {
	case model.GameScoreLove:
		return model.GameScoreFifteen
	case model.GameScoreFifteen:
		return model.GameScoreThirty
	case model.GameScoreThirty:
		return model.GameScoreForty
	default:
		return model.GameScoreGame
	}

}

func IsGameFinished(score *model.Score) bool {
	if score.A.CurrentPointScore == model.GameScoreGame || score.B.CurrentPointScore == model.GameScoreGame {
		return true
	} else {
		return false
	}

}
