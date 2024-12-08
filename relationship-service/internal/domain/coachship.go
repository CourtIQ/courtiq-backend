// internal/domain/coachship.go
package domain

import "time"

type Coachship struct {
	ID             string             `bson:"_id,omitempty" json:"id"`
	ParticipantIDs []string           `bson:"participantIds" json:"participantIds"`
	Type           RelationshipType   `bson:"type" json:"type"`
	Status         RelationshipStatus `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      *time.Time         `bson:"updatedAt,omitempty" json:"-"`
	CoachID        string             `bson:"coachId" json:"coachId"`
	CoacheeID      string             `bson:"coacheeId" json:"coacheeId"`
}

var _ Relationship = (*Coachship)(nil)

// Relationship interface methods
func (c *Coachship) GetID() string                 { return c.ID }
func (c *Coachship) GetType() RelationshipType     { return c.Type }
func (c *Coachship) GetStatus() RelationshipStatus { return c.Status }
func (c *Coachship) GetCreatedAt() time.Time       { return c.CreatedAt }
func (c *Coachship) GetUpdatedAt() *time.Time      { return c.UpdatedAt }
func (c *Coachship) GetParticipantIDs() []string   { return c.ParticipantIDs }

// NewCoachship creates a new Coachship instance with default status and timestamps.
// coachID and coacheeID define the participants in the coachship.
func NewCoachship(coachID, coacheeID string) *Coachship {
	now := time.Now().UTC()
	return &Coachship{
		ParticipantIDs: []string{coachID, coacheeID},
		Type:           RelationshipTypeCoachship,
		Status:         RelationshipStatusPending,
		CreatedAt:      now,
		UpdatedAt:      &now,
		CoachID:        coachID,
		CoacheeID:      coacheeID,
	}
}
