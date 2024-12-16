package service

import (
	"context"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpServiceIntf defines the interface for matchup operations
type MatchUpServiceIntf interface {
	CreateMatchUp(ctx context.Context, matchUpFormat model.MatchUpFormatInput) (*model.MatchUp, error)
	GetMatchUp(ctx context.Context, id string) (*model.MatchUp, error)
}

// MatchUpService implements the MatchUpServiceIntf
type MatchUpService struct {
	matchUpRepo repository.MatchUpRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService(matchUpRepo repository.MatchUpRepository) MatchUpServiceIntf {
	return &MatchUpService{matchUpRepo: matchUpRepo}
}

func (s *MatchUpService) CreateMatchUp(ctx context.Context, matchUpFormat model.MatchUpFormatInput) (*model.MatchUp, error) {
	matchUp := &model.MatchUp{
		ID:                          primitive.NewObjectID(),
		MatchUpFormat:               nil,
		MatchUpStatus:               "",
		MatchUpType:                 "",
		ParticipantIds:              nil,
		Participants:                nil,
		CurrentSetIndex:             nil,
		CurrentGameIndexWithinSet:   nil,
		CurrentPointIndexWithinGame: nil,
		CurrentScore:                nil,
		CurrentServer:               primitive.ObjectID{},
		PointsSequence:              nil,
		StartTime:                   time.Now(),
		EndTime:                     nil,
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	if err := s.matchUpRepo.Insert(ctx, matchUp); err != nil {
		return nil, err
	}

	return matchUp, nil
}

// GetMatchUp retrieves a matchup by ID (dummy implementation)
func (s *MatchUpService) GetMatchUp(ctx context.Context, id string) (*model.MatchUp, error) {
	// This is a dummy implementation
	// In a real implementation, you would fetch this from your database

	return &model.MatchUp{
		ID:            primitive.NewObjectID(),
		MatchUpStatus: model.MatchUpStatusInProgress,
		MatchUpType:   model.MatchUpTypeSingles,
		StartTime:     time.Now(),
		CreatedAt:     time.Now(), // Created 1 hour before start
		UpdatedAt:     time.Now(),
	}, nil
}
