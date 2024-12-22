package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipServiceIntf defines a 1:1 mapping of all relationship-related
// resolver methods (Friendship and Coachship queries & mutations) to service methods.
type RelationshipServiceIntf interface {
	// ---------------------------------------------------------------------------
	// Friendship Queries
	// ---------------------------------------------------------------------------

	// friendship(id: ObjectID!): Friendship
	Friendship(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)

	// myFriends(limit: Int = 10, offset: Int = 0): [Friendship!]!
	MyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)

	// friends(ofUserID: ObjectID!, limit: Int = 10, offset: Int = 0): [Friendship!]!
	Friends(ctx context.Context, ofUserId primitive.ObjectID, limit *int, offset *int) ([]*model.Friendship, error)

	// myFriendRequests: [Friendship!]!
	MyFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)

	// sentFriendRequests: [Friendship!]!
	SentFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error)

	// friendshipStatus(otherUserId: ObjectID!): RelationshipStatus!
	FriendshipStatus(ctx context.Context, otherUserId primitive.ObjectID) (model.RelationshipStatus, error)

	// ---------------------------------------------------------------------------
	// Coachship Queries
	// ---------------------------------------------------------------------------

	// coach(id: ObjectID!): Coachship
	Coach(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)

	// student(id: ObjectID!): Coachship
	Student(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)

	// myCoaches(limit: Int = 10, offset: Int = 0): [Coachship!]!
	MyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// myStudents(limit: Int = 10, offset: Int = 0): [Coachship!]!
	MyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// myStudentRequests: [Coachship!]!
	MyStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// myCoachRequests: [Coachship!]!
	MyCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// sentCoachRequests: [Coachship!]!
	SentCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// sentStudentRequests: [Coachship!]!
	SentStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error)

	// isStudent(studentId: ObjectID!): RelationshipStatus!
	IsStudent(ctx context.Context, studentId primitive.ObjectID) (model.RelationshipStatus, error)

	// isCoach(coachId: ObjectID!): RelationshipStatus!
	IsCoach(ctx context.Context, coachId primitive.ObjectID) (model.RelationshipStatus, error)

	// ---------------------------------------------------------------------------
	// Friendship Mutations
	// ---------------------------------------------------------------------------

	// sendFriendRequest(receiverId: ObjectID!): Friendship
	SendFriendRequest(ctx context.Context, receiverId primitive.ObjectID) (*model.Friendship, error)

	// acceptFriendRequest(friendshipId: ObjectID!): Friendship
	AcceptFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error)

	// rejectFriendRequest(friendshipId: ObjectID!): Friendship
	RejectFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error)

	// cancelFriendRequest(friendshipId: ObjectID!): Friendship
	CancelFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error)

	// endFriendship(friendshipId: ObjectID!): Friendship
	EndFriendship(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error)

	// ---------------------------------------------------------------------------
	// Coachship Mutations
	// ---------------------------------------------------------------------------

	// requestToBeStudent(ofUserId: ObjectID!): Coachship
	RequestToBeStudent(ctx context.Context, ofUserId primitive.ObjectID) (*model.Coachship, error)

	// acceptStudentRequest(coachshipId: ObjectID!): Coachship
	AcceptStudentRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// rejectStudentRequest(coachshipId: ObjectID!): Coachship
	RejectStudentRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// cancelRequestToBeStudent(coachshipId: ObjectID!): Coachship
	CancelRequestToBeStudent(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// removeStudent(coachshipId: ObjectID!): Coachship
	RemoveStudent(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// requestToBeCoach(ofUserId: ObjectID!): Coachship
	RequestToBeCoach(ctx context.Context, ofUserId primitive.ObjectID) (*model.Coachship, error)

	// acceptCoachRequest(coachshipId: ObjectID!): Coachship
	AcceptCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// rejectCoachRequest(coachshipId: ObjectID!): Coachship
	RejectCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// cancelCoachRequest(coachshipId: ObjectID!): Coachship
	CancelCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// removeCoach(coachshipId: ObjectID!): Coachship
	RemoveCoach(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error)

	// ---------------------------------------------------------------------------
	// Entity Resolvers (used internally by GraphQL resolvers if needed)
	// ---------------------------------------------------------------------------

	FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)

	// FindRelationshipByID can dynamically load either a Coachship or Friendship,
	// returning it as the interface type (model.Relationship).
	FindRelationshipByID(ctx context.Context, id primitive.ObjectID) (model.Relationship, error)
}
