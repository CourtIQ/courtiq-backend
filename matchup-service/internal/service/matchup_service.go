package service

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
)

// MatchUpServiceInterface defines the interface for matchup operations
type MatchUpServiceInterface interface {
	InitiateMatchUp(ctx context.Context, input model.CreateMatchUpInput) (*model.MatchUp, error)
}

// MatchUpService implements the MatchUpServiceInterface
type MatchUpService struct {
	matchupRepo repository.MatchUpRepository
	pointsReop  repository.PointsRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService() MatchUpServiceInterface {
	return &MatchUpService{}
}

// InitiateMatchUp creates a new MatchUp document in the database
func (s *MatchUpService) InitiateMatchUp(ctx context.Context, input model.CreateMatchUpInput) (*model.MatchUp, error) {
	return nil, errors.New("not implemented")
}
