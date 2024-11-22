// interfaces/user_service.go
package interfaces

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, user *models.User) (*models.User, error)
}
