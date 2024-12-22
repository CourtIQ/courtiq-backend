package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendshipRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	Insert(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error)
	Find(ctx context.Context, filter interface{}, limit *int, offset *int) ([]*model.Friendship, error)
	FindOne(ctx context.Context, filter interface{}) (*model.Friendship, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*model.Friendship, error)
	FindOneAndDelete(ctx context.Context, filter interface{}) (*model.Friendship, error)
}

type CoachshipRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	Insert(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error)
	Delete(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error)
	Find(ctx context.Context, filter interface{}, limit *int, offset *int) ([]*model.Coachship, error)
	FindOne(ctx context.Context, filter interface{}) (*model.Coachship, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*model.Coachship, error)
	FindOneAndDelete(ctx context.Context, filter interface{}) (*model.Coachship, error)
}
