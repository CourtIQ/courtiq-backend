package services

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
)

// MatchUpServiceInterface defines the interface for matchup operations
type MatchUpServiceInterface interface {
	InitiateMatchUp(ctx context.Context, input model.InitiateMatchUpInput) (*model.MatchUp, error)
}

// MatchUpService implements the MatchUpServiceInterface
type MatchUpService struct {
	matchupRepo repository.MatchUpRepositoryIntf
	pointsRepo  repository.PointsRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService(matchupRepo repository.MatchUpRepositoryIntf, pointsRepo repository.PointsRepository) MatchUpServiceInterface {
	return &MatchUpService{
		matchupRepo: matchupRepo,
		pointsRepo:  pointsRepo,
	}
}

// InitiateMatchUp creates a new MatchUp document in the database
func (s *MatchUpService) InitiateMatchUp(ctx context.Context, input model.InitiateMatchUpInput) (*model.MatchUp, error) {
	mu, err := NewMatchUpFromInitiateInput(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to build MatchUp from input: %w", err)
	}

	// 2) Insert the newly built MatchUp into the database
	created, err := s.matchupRepo.CreateMatchUp(ctx, mu)
	if err != nil {
		return nil, fmt.Errorf("failed to create MatchUp in repository: %w", err)
	}

	// 3) Return the created document
	return created, nil
}
