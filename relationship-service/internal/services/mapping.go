// internal/services/mapping.go
package services

import (
	"time"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
)

func domainFriendshipToModel(f *domain.Friendship) *model.Friendship {
	if f == nil {
		return nil
	}
	var updatedAt *string
	if f.GetUpdatedAt() != nil {
		u := f.GetUpdatedAt().UTC().Format(time.RFC3339)
		updatedAt = &u
	}
	createdAt := f.GetCreatedAt().UTC().Format(time.RFC3339)

	return &model.Friendship{
		ID:             f.GetID(),
		ParticipantIds: f.GetParticipantIDs(),
		Type:           model.RelationshipType(f.GetType()),
		Status:         model.RelationshipStatus(f.GetStatus()),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		SenderID:       f.SenderID,
		ReceiverID:     f.ReceiverID,
	}
}

func domainCoachshipToModel(c *domain.Coachship) *model.Coachship {
	if c == nil {
		return nil
	}
	var updatedAt *string
	if c.GetUpdatedAt() != nil {
		u := c.GetUpdatedAt().UTC().Format(time.RFC3339)
		updatedAt = &u
	}
	createdAt := c.GetCreatedAt().UTC().Format(time.RFC3339)

	return &model.Coachship{
		ID:             c.GetID(),
		ParticipantIds: c.GetParticipantIDs(),
		Type:           model.RelationshipType(c.GetType()),
		Status:         model.RelationshipStatus(c.GetStatus()),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		CoachID:        c.CoachID,
		StudentID:      c.StudentID,
	}
}

func relationshipStatusToModel(status domain.RelationshipStatus) *model.RelationshipStatus {
	s := model.RelationshipStatus(status)
	return &s
}
