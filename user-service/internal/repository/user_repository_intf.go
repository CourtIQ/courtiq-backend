package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, user *model.UpdateUserInput) (*model.User, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
}
