package services

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/factory"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/validation"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

// MatchUpService implements the MatchUpServiceIntf interface
type MatchUpService struct {
	matchupsRepo repository.MatchupsRepository
	shotsRepo    repository.ShotsRepository
}

// NewMatchUpService creates a new instance of MatchUpService
func NewMatchUpService(
	matchupsRepo repository.MatchupsRepository,
	shotsRepo repository.ShotsRepository,
) *MatchUpService {
	return &MatchUpService{
		matchupsRepo: matchupsRepo,
		shotsRepo:    shotsRepo,
	}
}

// GetMatchUpFormats retrieves all match up formats with pagination
func (s *MatchUpService) GetMatchUpFormats(ctx context.Context, limit *int, offset *int) ([]*model.MatchUpFormat, error) {
	return nil, ErrNotImplemented
}

// CreateMatchUpFormat creates a new match up format
func (s *MatchUpService) CreateMatchUpFormat(ctx context.Context, input model.MatchUpFormatInput) (*model.MatchUpFormat, error) {
	return nil, ErrNotImplemented
}

// InitiateMatchUp starts a new match up
func (s *MatchUpService) InitiateMatchUp(ctx context.Context, input model.InitiateMatchUpInput) (*model.MatchUp, error) {
	// Get current user from context
	ownerID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Validate input using the matchup validator
	matchupValidator := validation.NewMatchUpValidator()
	if err := matchupValidator.ValidateInitiateMatchUpInput(ctx, input); err != nil {
		return nil, err
	}

	// If matchUpFormat is provided, validate it using the format validator
	if input.MatchUpFormat != nil {
		formatValidator := validation.NewFormatValidator()
		if err := formatValidator.ValidateMatchUpFormatInput(ctx, *input.MatchUpFormat); err != nil {
			return nil, err
		}
	}

	// Create a match factory and create a new match from the input
	factory := factory.NewMatchUpFactory()
	matchUp := factory.CreateMatchUpFromInitiateMatchUpInput(ownerID, input)

	// Save the match to the database
	createdMatchUp, err := s.matchupsRepo.Insert(ctx, matchUp)
	if err != nil {
		return nil, err
	}

	return createdMatchUp, nil
}

// GetMatchUps retrieves all match ups with pagination
func (s *MatchUpService) GetMatchUps(ctx context.Context, limit *int, offset *int) ([]*model.MatchUp, error) {
	return nil, ErrNotImplemented
}

// GetMatchUpById retrieves a match up by its ID
func (s *MatchUpService) GetMatchUpById(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error) {
	return nil, ErrNotImplemented
}

// UpdateMatchUp updates a match up's status
func (s *MatchUpService) UpdateMatchUp(ctx context.Context, id primitive.ObjectID, status model.MatchUpStatus) (*model.MatchUp, error) {
	return nil, ErrNotImplemented
}

// AddShot adds a new shot to a match up
func (s *MatchUpService) AddShot(ctx context.Context, input model.AddShotInput) (*model.MatchUpShot, error) {
	return nil, ErrNotImplemented
}

// GetMatchUpShots retrieves all shots for a match up with pagination
func (s *MatchUpService) GetMatchUpShots(ctx context.Context, matchUpId primitive.ObjectID, limit *int, offset *int) ([]*model.MatchUpShot, error) {
	return nil, ErrNotImplemented
}

// GetShotsByGame retrieves all shots for a specific game in a match up
func (s *MatchUpService) GetShotsByGame(ctx context.Context, matchUpId primitive.ObjectID, setNumber int, gameNumber int) ([]*model.MatchUpShot, error) {
	return nil, ErrNotImplemented
}
