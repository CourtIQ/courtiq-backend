package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FriendshipRepositoryImpl implements FriendshipRepository
type FriendshipRepositoryImpl struct {
	baseRepo *repository.BaseRepository[model.Friendship]
}

// NewFriendshipRepository creates a new instance of FriendshipRepository
func NewFriendshipRepository(factory *repository.RepositoryFactory) FriendshipRepository {
	baseRepo := repository.NewRepository[model.Friendship](factory, RelationshipsCollection)
	return &FriendshipRepositoryImpl{
		baseRepo: baseRepo,
	}
}

// FindByID finds a friendship by its ID
func (r *FriendshipRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	filter := bson.M{
		"type": model.RelationshipTypeFriendship,
	}

	friendship, err := r.baseRepo.FindByIDWithFilters(ctx, id, filter)
	if err != nil {
		return nil, WrapRepositoryError(err, "find by ID", "friendship")
	}
	return friendship, nil
}

// FindBetweenUsers finds a friendship between two users
func (r *FriendshipRepositoryImpl) FindBetweenUsers(ctx context.Context, userID1, userID2 primitive.ObjectID) (*model.Friendship, error) {
	filter := bson.M{
		"$or": []bson.M{
			{
				"initiator._id": userID1,
				"receiver._id":  userID2,
			},
			{
				"initiator._id": userID2,
				"receiver._id":  userID1,
			},
		},
		"type": model.RelationshipTypeFriendship,
	}

	friendship, err := r.baseRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, WrapRepositoryError(err, "find between users", "friendship")
	}
	return friendship, nil
}

// GetFriendships gets all friendships for a user with a specific status
func (r *FriendshipRepositoryImpl) GetFriendships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Friendship, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"initiator._id": userID},
			{"receiver._id": userID},
		},
		"status": status,
		"type":   model.RelationshipTypeFriendship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetSentRequests gets all friend requests sent by a user
func (r *FriendshipRepositoryImpl) GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error) {
	filter := bson.M{
		"initiator._id": userID,
		"status":        model.RelationshipStatusPending,
		"type":          model.RelationshipTypeFriendship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetReceivedRequests gets all friend requests received by a user
func (r *FriendshipRepositoryImpl) GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Friendship, error) {
	filter := bson.M{
		"receiver._id": userID,
		"status":       model.RelationshipStatusPending,
		"type":         model.RelationshipTypeFriendship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// Create creates a new friendship
func (r *FriendshipRepositoryImpl) Create(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	// Set _id if not already set
	if friendship.ID == primitive.NilObjectID {
		friendship.ID = primitive.NewObjectID()
	}

	// Ensure type is set to FRIENDSHIP
	friendship.Type = model.RelationshipTypeFriendship

	created, err := r.baseRepo.Insert(ctx, friendship)
	if err != nil {
		return nil, WrapRepositoryError(err, "create", "friendship")
	}
	return created.(*model.Friendship), nil
}

// Update updates an existing friendship
func (r *FriendshipRepositoryImpl) Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	updated, err := r.baseRepo.Update(ctx, friendship.ID, friendship)
	if err != nil {
		return nil, WrapRepositoryError(err, "update", "friendship")
	}
	return updated, nil
}

// Delete deletes a friendship
func (r *FriendshipRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	err := r.baseRepo.Delete(ctx, id)
	if err != nil {
		return false, WrapRepositoryError(err, "delete", "friendship")
	}
	return true, nil
}
