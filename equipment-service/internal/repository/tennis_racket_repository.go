package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tennisRacketMongoRepo struct {
	collection *mongo.Collection
}

// NewTennisRacketMongoRepo creates a new TennisRacketRepository implementation
func NewTennisRacketMongoRepo(mdb *db.MongoDB) TennisRacketRepository {
	return &tennisRacketMongoRepo{
		collection: mdb.GetCollection(db.TennisRacketsCollection),
	}
}

func (r *tennisRacketMongoRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error) {
	var racket model.TennisRacket
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&racket)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &racket, err
}

func (r *tennisRacketMongoRepo) Insert(ctx context.Context, racket *model.TennisRacket) error {
	_, err := r.collection.InsertOne(ctx, racket)
	return err
}

func (r *tennisRacketMongoRepo) Update(ctx context.Context, racket *model.TennisRacket) error {
	filter := bson.M{"_id": racket.ID}
	update := bson.M{"$set": racket}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *tennisRacketMongoRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *tennisRacketMongoRepo) Find(ctx context.Context, filter interface{}) ([]*model.TennisRacket, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rackets []*model.TennisRacket
	if err := cursor.All(ctx, &rackets); err != nil {
		return nil, err
	}
	return rackets, nil
}
