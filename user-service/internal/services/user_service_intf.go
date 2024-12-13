package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model" // Adjust import path based on your project structure
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService defines the operations for managing users.
type UserServiceIntf interface {
	// Me retrieves the profile of the currently authenticated user.
	// The user context is assumed to have the necessary information.
	Me(ctx context.Context) (*model.User, error)

	// GetUser retrieves a user's profile by their unique ID.
	// Returns an error if the user is not found.
	GetUser(ctx context.Context, id primitive.ObjectID) (*model.User, error)

	// UpdateUser updates a user's profile based on the input data.
	// Returns the updated User object or an error if the update fails.
	UpdateUser(ctx context.Context, input *model.UpdateUserInput) (*model.User, error)

	// IsUsernameAvailable checks if the given username is available.
	// Returns true if the username is available, false otherwise.
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
}
