package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// RelationshipsCollection is the name of the MongoDB collection for relationships
	RelationshipsCollection = db.RelationshipsCollection
)

// FriendshipRepository defines the interface for friendship repository operations
type FriendshipRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	FindBetweenUsers(ctx context.Context, userID1, userID2 primitive.ObjectID) (*model.Friendship, error)
	GetFriendships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Friendship, error)
	GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error)
	GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error)
	Create(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// CoachshipRepository defines the interface for coachship repository operations
type CoachshipRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	FindBetweenUsers(ctx context.Context, userID1, userID2 primitive.ObjectID) (*model.Coachship, error)
	GetCoachships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Coachship, error)
	GetCoaches(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetStudents(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetSentCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetSentStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	GetReceivedStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error)
	Create(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// WrapRepositoryError standardizes repository errors to UIErrors
func WrapRepositoryError(err error, operation string, resourceType string) error {
	if err == nil {
		return nil
	}

	// Create a more specific error message
	message := ""
	code := ""

	switch resourceType {
	case "friendship":
		code = "FRIENDSHIP_ERROR"
		message = "Error in friendship operation: " + operation
	case "coachship":
		code = "COACHSHIP_ERROR"
		message = "Error in coaching relationship operation: " + operation
	default:
		code = "REPOSITORY_ERROR"
		message = "Error in repository operation: " + operation
	}

	return utils.StandardizeError(utils.NewUIError(code, message, err))
}
