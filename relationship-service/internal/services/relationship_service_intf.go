package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipServiceIntf defines the service interface for relationship management
type RelationshipServiceIntf interface {
	// ---------------------------------------------------------------------------
	// Entity Resolvers
	// ---------------------------------------------------------------------------
	FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)

	// ---------------------------------------------------------------------------
	// Friendship Queries
	// ---------------------------------------------------------------------------
	CheckFriendshipStatus(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error)
	GetMyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)
	GetSentFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)
	GetReceivedFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)

	// ---------------------------------------------------------------------------
	// Friendship Mutations
	// ---------------------------------------------------------------------------
	SendFriendRequest(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error)
	AcceptFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error)
	RejectFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error)
	CancelFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error)
	RemoveFriend(ctx context.Context, friendID primitive.ObjectID) (bool, error)
	BlockUser(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error)
	UnblockUser(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error)

	// ---------------------------------------------------------------------------
	// Coachship Queries
	// ---------------------------------------------------------------------------
	IsCoachOf(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error)
	IsStudentOf(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error)
	GetCoachships(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetMyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetMyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetSentCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetReceivedCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetSentStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	GetReceivedStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// ---------------------------------------------------------------------------
	// Coachship Mutations
	// ---------------------------------------------------------------------------
	RequestToBeCoachOf(ctx context.Context, userID primitive.ObjectID) (*model.Coachship, error)
	RequestToBeCoachedBy(ctx context.Context, userID primitive.ObjectID) (*model.Coachship, error)
	AcceptToBeCoachOf(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	RejectToBeCoachOf(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	AcceptToBeCoachedBy(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	RejectToBeCoachedBy(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	CancelCoachRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	CancelStudentRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error)
	EndCoachingAsCoach(ctx context.Context, coachshipID primitive.ObjectID) (bool, error)
	EndCoachingAsStudent(ctx context.Context, coachshipID primitive.ObjectID) (bool, error)
}
