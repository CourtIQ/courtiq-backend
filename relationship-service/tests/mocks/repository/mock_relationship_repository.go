// tests/mocks/repository/mock_relationship_repository.go
package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRelationshipRepository struct {
	mock.Mock
}

func (m *MockRelationshipRepository) Create(ctx context.Context, r domain.Relationship) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRelationshipRepository) GetByID(ctx context.Context, id string) (domain.Relationship, error) {
	args := m.Called(ctx, id)
	var rel domain.Relationship
	if v, ok := args.Get(0).(domain.Relationship); ok {
		rel = v
	}
	return rel, args.Error(1)
}

func (m *MockRelationshipRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	args := m.Called(ctx, id, fields)
	return args.Error(0)
}

func (m *MockRelationshipRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRelationshipRepository) ListByStatus(ctx context.Context, status domain.RelationshipStatus, limit int, offset int) ([]domain.Relationship, error) {
	args := m.Called(ctx, status, limit, offset)
	var rels []domain.Relationship
	if v, ok := args.Get(0).([]domain.Relationship); ok {
		rels = v
	}
	return rels, args.Error(1)
}

func (m *MockRelationshipRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	count, _ := args.Get(0).(int64)
	return count, args.Error(1)
}
