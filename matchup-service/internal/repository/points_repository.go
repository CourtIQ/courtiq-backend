package repository

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PointsRepository interface {
	CreatePoint(ctx context.Context, point *model.MatchUpPoint) (*model.MatchUpPoint, error)
	FindPointByID(ctx context.Context, id string) (*model.MatchUpPoint, error)
	UpdatePoint(ctx context.Context, point *model.MatchUpPoint) (*model.MatchUpPoint, error)
	DeletePoint(ctx context.Context, id string) (*model.MatchUpPoint, error)
}

type pointsRepository struct {
	pointsColl *mongo.Collection
}

func NewPointsRepositoru(mdb *db.MongoDB) PointsRepository {
	return &pointsRepository{
		pointsColl: mdb.GetCollection(db.TennisMatchupsPointsCollection),
	}
}

func (r *pointsRepository) CreatePoint(ctx context.Context, matchup *model.MatchUpPoint) (*model.MatchUpPoint, error) {
	matchup.ID = primitive.NewObjectID()

	_, err := r.pointsColl.InsertOne(ctx, matchup)
	if err != nil {
		return nil, err
	}

	return matchup, nil
}

func (r *pointsRepository) FindPointByID(ctx context.Context, id string) (*model.MatchUpPoint, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ObjectID string")
	}

	var found model.MatchUpPoint
	err = r.pointsColl.FindOne(ctx, bson.M{"_id": objID}).Decode(&found)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("matchUp not found")
		}
		return nil, err
	}

	return &found, nil
}

func (r *pointsRepository) DeletePoint(ctx context.Context, id string) (*model.MatchUpPoint, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ObjectID string")
	}

	var deleted model.MatchUpPoint
	err = r.pointsColl.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&deleted)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("point not found")
		}
		return nil, err
	}

	return &deleted, nil
}

func (r *pointsRepository) UpdatePoint(ctx context.Context, point *model.MatchUpPoint) (*model.MatchUpPoint, error) {
	if point == nil {
		return nil, errors.New("point cannot be nil")
	}
	if point.ID.IsZero() {
		return nil, errors.New("point must have a non-zero ID for update")
	}

	filter := bson.M{"_id": point.ID}

	// `replacement` is the new version of the document that will replace the old one.
	replacement := point

	_, err := r.pointsColl.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return nil, err
	}

	return point, nil
}
