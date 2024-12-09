package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

type RelationshipService interface {
	// Mutations - Friendships
	SendFriendRequest(ctx context.Context, receiverID string) (*bool, error)
	AcceptFriendRequest(ctx context.Context, friendshipID string) (*bool, error)
	RejectFriendRequest(ctx context.Context, friendshipID string) (*bool, error)
	CancelFriendRequest(ctx context.Context, friendshipID string) (*bool, error)
	EndFriendship(ctx context.Context, friendshipID string) (*bool, error)

	// Queries - Friendships
	GetFriendship(ctx context.Context, id string) (*model.Friendship, error)
	ListMyFriends(ctx context.Context, limit int, offset int) ([]*model.Friendship, error)
	ListFriends(ctx context.Context, ofUserID string, limit int, offset int) ([]*model.Friendship, error)
	ListFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	ListSentFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	CheckFriendshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error)

	// Mutations - Coachships
	RequestToBeStudent(ctx context.Context, ofUserID string) (*bool, error)
	AcceptStudentRequest(ctx context.Context, coachshipID string) (*bool, error)
	RejectStudentRequest(ctx context.Context, coachshipID string) (*bool, error)
	CancelStudentRequest(ctx context.Context, coachshipID string) (*bool, error)
	RemoveStudent(ctx context.Context, coachshipID string) (*bool, error)

	RequestToBeCoach(ctx context.Context, ofUserID string) (*bool, error)
	AcceptCoachRequest(ctx context.Context, coachshipID string) (*bool, error)
	RejectCoachRequest(ctx context.Context, coachshipID string) (*bool, error)
	CancelCoachRequest(ctx context.Context, coachshipID string) (*bool, error)
	RemoveCoach(ctx context.Context, coachshipID string) (*bool, error)

	// Queries - Coachships
	GetCoachship(ctx context.Context, id string) (*model.Coachship, error)
	ListCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	ListStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	ListSentStudentRequests(ctx context.Context) ([]*model.Coachship, error)
	ListReceivedCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	ListSentCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	ListReceivedStudentRequests(ctx context.Context) ([]*model.Coachship, error)
	IsStudent(ctx context.Context, studentID string) (*model.RelationshipStatus, error)
	IsCoach(ctx context.Context, coachID string) (*model.RelationshipStatus, error)
}
