package repository

import (
	"context"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchUpRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error)
	Insert(ctx context.Context, matchUp *model.MatchUp) error
	Update(ctx context.Context, matchUp *model.MatchUp) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Find(ctx context.Context, filter interface{}) ([]*model.MatchUp, error)
}
