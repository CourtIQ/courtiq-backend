package interfaces

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
)

type UserService interface {
	// Find user by ID
	FindUserByID(ctx context.Context, id string) (*model.User, error)

	// Update user details
	UpdateUser(ctx context.Context, input model.UserUpdateInput) (*model.User, error)

	// Delete user
	DeleteUser(ctx context.Context, id string) (bool, error)

	// Check username availability
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)

	// Fetch current user profile
	FetchCurrentUser(ctx context.Context) (*model.User, error)
}
