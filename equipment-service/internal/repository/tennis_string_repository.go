package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tennisStringMongoRepo struct {
	collection *mongo.Collection
}

// NewTennisStringMongoRepo creates a new TennisStringRepository implementation
func NewTennisStringMongoRepo(mdb *db.MongoDB) TennisStringRepository {
	return &tennisStringMongoRepo{
		collection: mdb.GetCollection(db.TennisStringsCollection),
	}
}

func (r *tennisStringMongoRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error) {
	var s model.TennisString
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

func (r *tennisStringMongoRepo) Insert(ctx context.Context, s *model.TennisString) error {
	_, err := r.collection.InsertOne(ctx, s)
	return err
}

func (r *tennisStringMongoRepo) Update(ctx context.Context, s *model.TennisString) error {
	filter := bson.M{"_id": s.ID}
	update := bson.M{"$set": s}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *tennisStringMongoRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *tennisStringMongoRepo) Find(ctx context.Context, filter interface{}) ([]*model.TennisString, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var strings []*model.TennisString
	if err := cursor.All(ctx, &strings); err != nil {
		return nil, err
	}
	return strings, nil
}
