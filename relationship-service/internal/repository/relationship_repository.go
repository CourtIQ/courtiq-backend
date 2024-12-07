// internal/repository/relationship_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type relationshipRepository struct {
	coll *mongo.Collection
}

func NewRelationshipRepository(coll *mongo.Collection) RelationshipRepository {
	return &relationshipRepository{coll: coll}
}

func (r *relationshipRepository) Create(rel domain.Relationship) error {
	ctx := context.TODO() // In real code, pass context down from callers instead of using TODO()

	_, err := r.coll.InsertOne(ctx, rel)
	if err != nil {
		return fmt.Errorf("failed to insert relationship: %w", err)
	}
	return nil
}

func (r *relationshipRepository) GetByID(id string) (domain.Relationship, error) {
	return nil, errors.New("GetByID not implemented")
}

func (r *relationshipRepository) Update(id string, fields map[string]interface{}) error {
	return errors.New("Update not implemented")
}

func (r *relationshipRepository) Delete(id string) error {
	return errors.New("Delete not implemented")
}

func (r *relationshipRepository) ListByStatus(status domain.RelationshipStatus, limit int, offset int) ([]domain.Relationship, error) {
	return nil, errors.New("ListByStatus not implemented")
}
