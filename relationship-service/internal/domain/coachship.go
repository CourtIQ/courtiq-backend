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

// Implement the Relationship interface
func (c *Coachship) GetID() string                 { return c.ID }
func (c *Coachship) GetType() RelationshipType     { return c.Type }
func (c *Coachship) GetStatus() RelationshipStatus { return c.Status }
func (f *Coachship) GetCreatedAt() time.Time {
	return f.CreatedAt
}
func (f *Coachship) GetUpdatedAt() *time.Time {
	return f.UpdatedAt
}
func (c *Coachship) GetParticipantIDs() []string { return c.ParticipantIDs }
