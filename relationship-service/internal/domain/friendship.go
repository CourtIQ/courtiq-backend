// internal/domain/friendship.go
package domain

import "time"

type Friendship struct {
	ID             string             `bson:"_id,omitempty" json:"id"`
	ParticipantIDs []string           `bson:"participantIds" json:"participantIds"`
	Type           RelationshipType   `bson:"type" json:"type"`
	Status         RelationshipStatus `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      *time.Time         `bson:"updatedAt,omitempty" json:"-"`
	SenderID       string             `bson:"senderId" json:"senderId"`
	ReceiverID     string             `bson:"receiverId" json:"receiverId"`
}

var _ Relationship = (*Friendship)(nil)

// Relationship interface methods
func (f *Friendship) GetID() string                 { return f.ID }
func (f *Friendship) GetType() RelationshipType     { return f.Type }
func (f *Friendship) GetStatus() RelationshipStatus { return f.Status }
func (f *Friendship) GetCreatedAt() time.Time {
	return f.CreatedAt
}
func (f *Friendship) GetUpdatedAt() *time.Time {
	return f.UpdatedAt
}
func (f *Friendship) GetParticipantIDs() []string {
	return f.ParticipantIDs
}

func NewFriendship(senderId, receiverID string) *Friendship {
	now := time.Now().UTC()
	return &Friendship{
		ParticipantIDs: []string{senderId, receiverID},
		Type:           RelationshipTypeFriendship,
		Status:         RelationshipStatusPending,
		CreatedAt:      now,
		UpdatedAt:      &now,
		SenderID:    senderId,
		ReceiverID:     receiverID,
	}
}
