// internal/repository/relationship_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type relationshipRepository struct {
	coll *mongo.Collection
}

func NewRelationshipRepository(coll *mongo.Collection) RelationshipRepository {
	return &relationshipRepository{coll: coll}
}

func (r *relationshipRepository) Create(ctx context.Context, rel domain.Relationship) error {
	_, err := r.coll.InsertOne(ctx, rel)
	if err != nil {
		return fmt.Errorf("failed to insert relationship: %w", err)
	}
	return nil
}

func (r *relationshipRepository) GetByID(ctx context.Context, id string) (domain.Relationship, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var relationship domain.Relationship
	err = r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&relationship)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no relationship found with ID: %s", id)
		}
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	return relationship, nil
}

func (r *relationshipRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": fields}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update relationship: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no relationship found with ID: %s", id)
	}

	return nil
}

func (r *relationshipRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no relationship found with ID: %s", id)
	}

	return nil
}

func (r *relationshipRepository) ListByStatus(ctx context.Context, status domain.RelationshipStatus, limit int, offset int) ([]domain.Relationship, error) {
	return nil, errors.New("ListByStatus not implemented")
}

func (r *relationshipRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	// Convert the map[string]interface{} filter to bson.M for MongoDB
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	count, err := r.coll.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}
