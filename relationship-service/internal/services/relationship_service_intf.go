package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipServiceIntf defines a 1:1 mapping of all relationship-related
// resolver methods (Friendship and Coachship queries & mutations) to service methods.
type RelationshipServiceIntf interface {
	// Friendship Queries
	Friendship(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	MyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)
	Friends(ctx context.Context, ofUserID string, limit *int, offset *int) ([]*model.Friendship, error)
	MyFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	SentFriendRequests(ctx context.Context) ([]*model.Friendship, error)
	FriendshipStatus(ctx context.Context, otherUserID primitive.ObjectID) (*model.RelationshipStatus, error)

	// Coachship Queries
	Coachship(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	MyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	MyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)
	MyStudentRequests(ctx context.Context) ([]*model.Coachship, error)
	MyCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	SentCoachRequests(ctx context.Context) ([]*model.Coachship, error)
	SentStudentRequests(ctx context.Context) ([]*model.Coachship, error)
	IsStudent(ctx context.Context, studentID string) (*model.RelationshipStatus, error)
	IsCoach(ctx context.Context, coachID string) (*model.RelationshipStatus, error)

	// Friendship Mutations
	SendFriendRequest(ctx context.Context, receiverID string) (*model.Friendship, error)
	AcceptFriendRequest(ctx context.Context, friendshipID primitive.ObjectID) (*model.Friendship, error)
	RejectFriendRequest(ctx context.Context, friendshipID primitive.ObjectID) (*model.Friendship, error)
	CancelFriendRequest(ctx context.Context, friendshipID primitive.ObjectID) (*model.Friendship, error)
	EndFriendship(ctx context.Context, friendshipID primitive.ObjectID) (*model.Friendship, error)

	// Coachship Mutations
	RequestToBeStudent(ctx context.Context, ofUserID string) (*model.Coachship, error)
	AcceptStudentRequest(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	RejectStudentRequest(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	CancelRequestToBeStudent(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	RemoveStudent(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)

	RequestToBeCoach(ctx context.Context, ofUserID string) (*model.Coachship, error)
	AcceptCoachRequest(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	RejectCoachRequest(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	CancelCoachRequest(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)
	RemoveCoach(ctx context.Context, coachshipID primitive.ObjectID) (*model.Coachship, error)

	// Entity Resolvers
	FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	FindRelationshipByID(ctx context.Context, id primitive.ObjectID) (model.Relationship, error)
}
