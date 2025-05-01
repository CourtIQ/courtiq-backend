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

// ShotsRepository defines the interface for matchup shots repository operations
type ShotsRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUpShot, error)
	FindByMatchUpID(ctx context.Context, matchUpID primitive.ObjectID) ([]*model.MatchUpShot, error)
	FindShotsByGame(ctx context.Context, matchUpID primitive.ObjectID, setNumber, gameNumber int) ([]*model.MatchUpShot, error)
	FindLastShot(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUpShot, error)
	Insert(ctx context.Context, shot *model.MatchUpShot) (*model.MatchUpShot, error)
	Update(ctx context.Context, shot *model.MatchUpShot) (*model.MatchUpShot, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// ShotsRepositoryImpl implements ShotsRepository
type ShotsRepositoryImpl struct {
	baseRepo *repository.BaseRepository[model.MatchUpShot]
	factory  *repository.RepositoryFactory
}

// NewShotsRepository creates a new instance of ShotsRepository
func NewShotsRepository(factory *repository.RepositoryFactory) ShotsRepository {
	baseRepo := repository.NewRepository[model.MatchUpShot](factory, db.TennisMatchupsShotsCollection)
	return &ShotsRepositoryImpl{
		baseRepo: baseRepo,
		factory:  factory,
	}
}

// FindByID finds a shot by its ID
func (r *ShotsRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUpShot, error) {
	shot, err := r.baseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return shot, nil
}

// FindByMatchUpID finds all shots for a specific matchup
func (r *ShotsRepositoryImpl) FindByMatchUpID(ctx context.Context, matchUpID primitive.ObjectID) ([]*model.MatchUpShot, error) {
	filter := bson.M{
		"matchUpId": matchUpID,
	}

	// Sort by timestamp to get shots in order
	opts := options.Find().SetSort(bson.M{"timestamp": 1})

	shots, err := r.baseRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return shots, nil
}

// FindShotsByGame finds all shots for a specific game within a match
func (r *ShotsRepositoryImpl) FindShotsByGame(ctx context.Context, matchUpID primitive.ObjectID, setNumber, gameNumber int) ([]*model.MatchUpShot, error) {
	filter := bson.M{
		"matchUpId":  matchUpID,
		"setNumber":  setNumber,
		"gameNumber": gameNumber,
	}

	// Sort by timestamp to get shots in order
	opts := options.Find().SetSort(bson.M{"timestamp": 1})

	shots, err := r.baseRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return shots, nil
}

// FindLastShot finds the most recent shot for a matchup
func (r *ShotsRepositoryImpl) FindLastShot(ctx context.Context, matchUpID primitive.ObjectID) (*model.MatchUpShot, error) {
	filter := bson.M{
		"matchUpId": matchUpID,
	}

	// Sort by timestamp in descending order and limit to 1
	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})

	var shot model.MatchUpShot
	collection := r.factory.GetCollection(db.TennisMatchupsShotsCollection)
	err := collection.FindOne(ctx, filter, opts).Decode(&shot)
	if err != nil {
		return nil, err
	}
	return &shot, nil
}

// Insert creates a new shot
func (r *ShotsRepositoryImpl) Insert(ctx context.Context, shot *model.MatchUpShot) (*model.MatchUpShot, error) {
	if shot.ID == primitive.NilObjectID {
		shot.ID = primitive.NewObjectID()
	}

	created, err := r.baseRepo.Insert(ctx, shot)
	if err != nil {
		return nil, err
	}
	return created.(*model.MatchUpShot), nil
}

// Update updates an existing shot
func (r *ShotsRepositoryImpl) Update(ctx context.Context, shot *model.MatchUpShot) (*model.MatchUpShot, error) {
	updated, err := r.baseRepo.Update(ctx, shot.ID, shot)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete deletes a shot
func (r *ShotsRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	err := r.baseRepo.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
