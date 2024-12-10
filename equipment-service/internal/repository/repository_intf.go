// internal/repository/repository_intf.go
package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TennisRacketRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error)
	Insert(ctx context.Context, racket *model.TennisRacket) error
	Update(ctx context.Context, racket *model.TennisRacket) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Find(ctx context.Context, filter interface{}) ([]*model.TennisRacket, error)
}

type TennisStringRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error)
	Insert(ctx context.Context, s *model.TennisString) error
	Update(ctx context.Context, s *model.TennisString) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Find(ctx context.Context, filter interface{}) ([]*model.TennisString, error)
}
