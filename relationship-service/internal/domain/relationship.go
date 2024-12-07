// internal/domain/relationship.go
package domain

import "time"

type RelationshipType string
type RelationshipStatus string

const (
	RelationshipTypeFriendship RelationshipType = "FRIENDSHIP"
	RelationshipTypeCoachship  RelationshipType = "COACHSHIP"

	RelationshipStatusPending RelationshipStatus = "PENDING"
	RelationshipStatusActive  RelationshipStatus = "ACTIVE"
	// Add other statuses like REJECTED, ENDED if needed
)

// Relationship interface defines common behavior that both Friendship and Coachship must implement.
type Relationship interface {
	GetID() string
	GetType() RelationshipType
	GetStatus() RelationshipStatus
	GetCreatedAt() time.Time
	GetUpdatedAt() *time.Time
	GetParticipantIDs() []string
}
