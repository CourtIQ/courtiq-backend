package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

type RelationshipService interface {
	// Friendships
	SendFriendRequest(ctx context.Context, receiverID string) (bool, error)
	AcceptFriendRequest(ctx context.Context, friendshipID string) (bool, error)
	RejectFriendRequest(ctx context.Context, friendshipID string) (bool, error)
	CancelFriendRequest(ctx context.Context, friendshipID string) (bool, error)
	EndFriendship(ctx context.Context, friendshipID string) (bool, error)

	// Coachships
	SendCoachRequest(ctx context.Context, userID string) (*model.Coachship, error)
	SendCoacheeRequest(ctx context.Context, userID string) (*model.Coachship, error)
	AcceptCoachRequest(ctx context.Context, coachshipID string) (*model.Coachship, error)
	AcceptCoacheeRequest(ctx context.Context, coachshipID string) (*model.Coachship, error)
	DeclineCoachRequest(ctx context.Context, coachshipID string) (bool, error)
	DeclineCoacheeRequest(ctx context.Context, coachshipID string) (bool, error)
	CancelCoachRequest(ctx context.Context, coachshipID string) (bool, error)
	CancelCoacheeRequest(ctx context.Context, coachshipID string) (bool, error)
	EndCoachship(ctx context.Context, coachshipID string) (bool, error)

	// Queries - Coachships
	GetCoachship(ctx context.Context, id string) (*model.Coachship, error)
	ListCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	ListStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	ListSentCoacheeRequests(ctx context.Context) ([]*model.Coachship, error)
	ListReceivedCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	ListSentCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	ListReceivedCoacheeRequests(ctx context.Context) ([]*model.Coachship, error)
	CheckCoachshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error)

	// Queries - Friendships
	GetFriendship(ctx context.Context, id string) (*model.Friendship, error)
	ListFriends(ctx context.Context, limit int, offset int) ([]*model.Friendship, error)
	ListPendingFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	ListSentFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	CheckFriendshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error)
}
