package service

import (
	"context"
	"time"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpServiceInterface defines the interface for matchup operations
type MatchUpServiceInterface interface {
	GetMatchUp(ctx context.Context, id string) (*model.MatchUp, error)
}

// MatchUpService implements the MatchUpServiceInterface
type MatchUpService struct {
	// You can add repository dependencies here
	// Example: matchupRepo repository.MatchUpRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService() MatchUpServiceInterface {
	return &MatchUpService{}
}

// GetMatchUp retrieves a matchup by ID (dummy implementation)
func (s *MatchUpService) GetMatchUp(ctx context.Context, id string) (*model.MatchUp, error) {
	// This is a dummy implementation
	// In a real implementation, you would fetch this from your database

	return &model.MatchUp{
		ID:          primitive.NewObjectID(),
		MatchUpType: model.MatchUpTypeSingles,
		StartTime:   time.Now(),
		CreatedAt:   time.Now(), // Created 1 hour before start
	}, nil
}
