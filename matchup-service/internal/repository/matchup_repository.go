package repository

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type matchUpMongoRepo struct {
	collection *mongo.Collection
}

// NewMatchUpMongoRepo creates a new NewMatchUpMongoRepository implementation
func NewMatchUpMongoRepo(mdb *db.MongoDB) MatchUpRepository {
	return &matchUpMongoRepo{
		collection: mdb.GetCollection(db.MatchupsCollection),
	}
}

// FindByID retrieves a matchup from the database by its ID
func (r *matchUpMongoRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error) {
	var matchUp model.MatchUp
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&matchUp)
	if err == mongo.ErrNoDocuments {
		return nil, nil // No document found
	}
	return &matchUp, err
}

// Insert adds a new matchup to the database
func (r *matchUpMongoRepo) Insert(ctx context.Context, matchUp *model.MatchUp) error {
	_, err := r.collection.InsertOne(ctx, matchUp)
	return err
}

// Update modifies an existing matchup in the database
func (r *matchUpMongoRepo) Update(ctx context.Context, matchUp *model.MatchUp) error {
	filter := bson.M{"_id": matchUp.ID}
	update := bson.M{"$set": matchUp}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a matchup from the database by ID
func (r *matchUpMongoRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Find retrieves matchups from the database based on a filter
func (r *matchUpMongoRepo) Find(ctx context.Context, filter interface{}) ([]*model.MatchUp, error) {
	// Logging for debugging purposes
	fmt.Printf("Executing Find with filter: %v\n", filter)

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matchUps []*model.MatchUp
	if err := cursor.All(ctx, &matchUps); err != nil {
		return nil, err
	}

	// Logging for debugging purposes
	fmt.Printf("Found %d matchups\n", len(matchUps))
	return matchUps, nil
}
