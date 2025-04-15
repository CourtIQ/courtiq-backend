package resolvers

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindCoachshipByID is the resolver for the findCoachshipByID field.
func (r *entityResolver) FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	// Fetch the coachship, ensuring user IDs (CoachID, StudentID etc.) are loaded
	coachship, err := r.RelationshipService.FindCoachshipByID(ctx, id)
	if err != nil {
		// Consider wrapping error using utils.StandardizeError or similar
		return nil, fmt.Errorf("failed to find coachship by ID %s: %w", id.Hex(), err)
	}
	// The resolvers for Coach, Student, Initiator, Receiver will handle populating the User objects.
	return coachship, nil
}

// FindFriendshipByID is the resolver for the findFriendshipByID field.
func (r *entityResolver) FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	// Fetch the friendship, ensuring user IDs (InitiatorID, ReceiverID) are loaded
	friendship, err := r.RelationshipService.FindFriendshipByID(ctx, id)
	if err != nil {
		// Consider wrapping error using utils.StandardizeError or similar
		return nil, fmt.Errorf("failed to find friendship by ID %s: %w", id.Hex(), err)
	}
	// The resolvers for Initiator, Receiver will handle populating the User objects.
	return friendship, nil
}

// FindUserByID is the resolver for the findUserByID field.
// This is crucial for federation. It tells the gateway how to reconstitute
// a User entity when its origin is this service (e.g., as part of a Friendship).
// It only needs to return a User object with the ID field populated.
func (r *entityResolver) FindUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	// For federation, we just need to return a User stub with the ID.
	// The gateway will use this ID to fetch the full details from the user-service.
	return &model.User{ID: id}, nil
}

// Entity returns graph.EntityResolver implementation.
func (r *Resolver) Entity() graph.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
