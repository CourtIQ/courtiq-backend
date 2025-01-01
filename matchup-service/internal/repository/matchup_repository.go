package repository

import (
	"context"
	"errors"
	"log"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchUpRepositoryIntf interface {
	CreateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	UpdateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	DeleteMatchUp(ctx context.Context, id string) (*model.MatchUp, error)
	FindMatchUpByID(ctx context.Context, id string) (*model.MatchUp, error)
}

type matchupRepository struct {
	matchupColl *mongo.Collection
}

func NewMatchUpRepository(mdb *db.MongoDB) MatchUpRepositoryIntf {
	return &matchupRepository{
		matchupColl: mdb.GetCollection(db.TennisMatchupsCollection),
	}
}

func (r *matchupRepository) CreateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error) {
	matchup.ID = primitive.NewObjectID()
	log.Printf("repo=%+v, repo.matchupColl=%+v", r, r.matchupColl)

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

func (r *matchupRepository) DeleteMatchUp(ctx context.Context, id string) (*model.MatchUp, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ObjectID string")
	}

	var deleted model.MatchUp
	err = r.matchupColl.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&deleted)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("point not found")
		}
		return nil, err
	}

	return &deleted, nil
}

// This approach does a full replacement of the document.
func (r *matchupRepository) UpdateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error) {
	if matchup == nil {
		return nil, errors.New("matchup cannot be nil")
	}
	if matchup.ID.IsZero() {
		return nil, errors.New("matchup must have a non-zero ID for update")
	}

	filter := bson.M{"_id": matchup.ID}

	replacement := matchup

	_, err := r.matchupColl.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return nil, err
	}

	return matchup, nil
}
