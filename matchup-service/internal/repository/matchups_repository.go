package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MatchupsRepository defines the interface for matchup repository operations
type MatchupsRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error)
	Insert(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	Update(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
	GetMatchups(ctx context.Context, limit, offset *int) ([]*model.MatchUp, error)
	GetMatchupsByTeam(ctx context.Context, teamID primitive.ObjectID, limit, offset *int) ([]*model.MatchUp, error)
	GetMyMatchupsByStatus(ctx context.Context, status model.MatchUpStatus, limit, offset *int) ([]*model.MatchUp, error)
}

// MatchupsRepositoryImpl implements MatchupsRepository
type MatchupsRepositoryImpl struct {
	baseRepo *repository.BaseRepository[model.MatchUp]
}

// NewMatchupsRepository creates a new instance of MatchupsRepository
func NewMatchupsRepository(factory *repository.RepositoryFactory) MatchupsRepository {
	baseRepo := repository.NewRepository[model.MatchUp](factory, db.TennisMatchupsCollection)
	return &MatchupsRepositoryImpl{
		baseRepo: baseRepo,
	}
}

// FindByID finds a matchup by its ID
func (r *MatchupsRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error) {
	matchup, err := r.baseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return matchup, nil
}

// Insert creates a new matchup
func (r *MatchupsRepositoryImpl) Insert(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error) {
	if matchup.ID == primitive.NilObjectID {
		matchup.ID = primitive.NewObjectID()
	}

	_, err := r.baseRepo.Insert(ctx, matchup)
	if err != nil {
		return nil, err
	}
	return matchup, nil
}

// Update updates an existing matchup
func (r *MatchupsRepositoryImpl) Update(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error) {
	updated, err := r.baseRepo.Update(ctx, matchup.ID, matchup)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete deletes a matchup
func (r *MatchupsRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	err := r.baseRepo.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetMatchups retrieves matchups with optional pagination
func (r *MatchupsRepositoryImpl) GetMatchups(ctx context.Context, limit, offset *int) ([]*model.MatchUp, error) {
	filter := bson.M{}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	matchups, err := r.baseRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return matchups, nil
}

// GetMatchupsByTeam retrieves matchups for a specific team with optional pagination
func (r *MatchupsRepositoryImpl) GetMatchupsByTeam(ctx context.Context, teamID primitive.ObjectID, limit, offset *int) ([]*model.MatchUp, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"participants.teamId": teamID},
		},
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	matchups, err := r.baseRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return matchups, nil
}

// GetMyMatchupsByStatus retrieves matchups for the current user with a specific status
func (r *MatchupsRepositoryImpl) GetMyMatchupsByStatus(ctx context.Context, status model.MatchUpStatus, limit, offset *int) ([]*model.MatchUp, error) {
	filter := bson.M{
		"status": status,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	matchups, err := r.baseRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return matchups, nil
}
