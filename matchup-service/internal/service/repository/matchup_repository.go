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

type MatchUpRepository interface {
	CreateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	FindMatchUpByID(ctx context.Context, id string) (*model.MatchUp, error)
}

type matchupRepository struct {
	pointsColl  *mongo.Collection
	matchupColl *mongo.Collection
}

func NewMatchUpRepository(mdb *db.MongoDB) MatchUpRepository {
	return &matchupRepository{
		pointsColl:  mdb.GetCollection(db.TennisMatchupsPointsCollection),
		matchupColl: mdb.GetCollection(db.TennisMatchupsCollection),
	}
}

func (r *matchupRepository) CreateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error) {
	matchup.ID = primitive.NewObjectID()

	_, err := r.matchupColl.InsertOne(ctx, matchup)
	if err != nil {
		return nil, err
	}

	return matchup, nil
}

func (r *matchupRepository) FindMatchUpByID(ctx context.Context, id string) (*model.MatchUp, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ObjectID string")
	}

	var found model.MatchUp
	err = r.matchupColl.FindOne(ctx, bson.M{"_id": objID}).Decode(&found)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("matchUp not found")
		}
		return nil, err
	}

	return &found, nil
}
