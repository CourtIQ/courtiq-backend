// internal/repository/relationship_repository_intf.go
package repository

import (
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
)

type RelationshipRepository interface {
	// Create inserts a new relationship (Friendship, Coachship, etc.) into the database.
	// The implementation will likely inspect r.GetType() to determine storage details.
	Create(r domain.Relationship) error

	// GetByID retrieves a relationship by its ID.
	// Returns a domain.Relationship interface, which could be a Friendship or Coachship.
	GetByID(id string) (domain.Relationship, error)

	// Update modifies specific fields of an existing relationship identified by ID.
	Update(id string, fields map[string]interface{}) error

	// Delete removes a relationship from the database by its ID.
	Delete(id string) error

	// ListByStatus retrieves relationships of any type by a given status, with pagination.
	// Returns a slice of domain.Relationship, allowing calling code to handle them generically.
	ListByStatus(status domain.RelationshipStatus, limit int, offset int) ([]domain.Relationship, error)
}
