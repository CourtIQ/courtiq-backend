// internal/repository/relationship_repository_intf.go
package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
)

// RelationshipRepository defines methods for interacting with relationships in the data store.
// By passing `context.Context` to each method, callers can propagate timeouts, cancellations,
// and request-scoped metadata down to the database operations.
type RelationshipRepository interface {
	// Create inserts a new relationship into the database.
	// The implementation can inspect r.GetType() to determine storage details.
	Create(ctx context.Context, r domain.Relationship) error

	// GetByID retrieves a relationship by its ID.
	// Returns a domain.Relationship interface, which could be a Friendship or Coachship.
	GetByID(ctx context.Context, id string) (domain.Relationship, error)

	// Update modifies specific fields of an existing relationship identified by ID.
	Update(ctx context.Context, id string, fields map[string]interface{}) error

	// Delete removes a relationship from the database by its ID.
	Delete(ctx context.Context, id string) error

	// ListByStatus retrieves relationships by a given status, with pagination.
	// Returns a slice of domain.Relationship, allowing the caller to handle them generically.
	ListByStatus(ctx context.Context, status domain.RelationshipStatus, limit int, offset int) ([]domain.Relationship, error)

	// Count returns the number of documents matching the given filter.
	// 'filter' is a map of conditions that must match for documents to be counted.
	// For example: filter could be bson.M{"status": "PENDING", "participantIds": userID}
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)

	GetFriendshipByID(ctx context.Context, id string) (*domain.Friendship, error)
}
