package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection names should leverage the shared constants
const (
	RelationshipsCollection = db.RelationshipsCollection
)

// FriendshipRepository defines methods for friendship data access
type FriendshipRepository interface {
	// Find methods
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	FindBetweenUsers(ctx context.Context, userID1, userID2 primitive.ObjectID) (*model.Friendship, error)

	// Get collections
	GetFriendships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Friendship, error)
	GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error)
	GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error)

	// Mutations
	Create(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// CoachshipRepository defines methods for coachship data access
type CoachshipRepository interface {
	// Find methods
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	FindBetweenUsers(ctx context.Context, userID1, userID2 primitive.ObjectID) (*model.Coachship, error)

	// Get collections
	GetCoachships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Coachship, error)
	GetCoaches(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetStudents(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)

	// Get requests
	GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetSentCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetSentStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)

	// Mutations
	Create(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
}
