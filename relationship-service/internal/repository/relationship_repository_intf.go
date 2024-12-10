package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FriendshipRepository defines the data persistence layer behavior for friendships.
type FriendshipRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	GetMyFriends(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Friendship, error)
	GetFriendsOfUser(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Friendship, error)
	GetMyFriendRequests(ctx context.Context, userID string) ([]*model.Friendship, error)
	GetSentFriendRequests(ctx context.Context, userID string) ([]*model.Friendship, error)
	Insert(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// CoachshipRepository defines the data persistence layer behavior for coachships.
type CoachshipRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	GetMyCoaches(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Coachship, error)
	GetMyStudents(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Coachship, error)
	GetMyStudentRequests(ctx context.Context, userID string) ([]*model.Coachship, error)
	GetMyCoachRequests(ctx context.Context, userID string) ([]*model.Coachship, error)
	GetSentCoachRequests(ctx context.Context, userID string) ([]*model.Coachship, error)
	GetSentStudentRequests(ctx context.Context, userID string) ([]*model.Coachship, error)
	Insert(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
