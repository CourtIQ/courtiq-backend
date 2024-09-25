// user_service_provider.go

package providers

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
)

// UserServiceProvider defines the contract for user-related operations.
type UserServiceProvider interface {
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, input model.UserUpdateInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
	FetchCurrentUser(ctx context.Context) (*model.User, error)
}
