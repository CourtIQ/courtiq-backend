package repository

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type friendshipRepository struct {
	collection *mongo.Collection
}

func NewFriendshipRepository(mdb *db.MongoDB) FriendshipRepository {
	return &friendshipRepository{
		collection: mdb.GetCollection(db.FriendshipsCollection),
	}
}

// FindByID finds a Friendship by its ObjectID.
func (r *friendshipRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	filter := bson.M{"_id": id}

	var friendship model.Friendship
	err := r.collection.FindOne(ctx, filter).Decode(&friendship)
	if err != nil {
		// Check if the document was not found in Mongo
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFriendshipNotFound
		}
		return nil, err
	}
	return &friendship, nil
}

// Insert creates a new Friendship document and returns the inserted record.
func (r *friendshipRepository) Insert(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	if friendship.ID.IsZero() {
		friendship.ID = primitive.NewObjectID()
	}

	res, err := r.collection.InsertOne(ctx, friendship)
	if err != nil {
		return nil, err
	}

	var insertedFriendship model.Friendship
	if err := r.collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&insertedFriendship); err != nil {
		return nil, err
	}

	return &insertedFriendship, nil
}

// Update updates an existing Friendship by its ObjectID and returns the updated record.
func (r *friendshipRepository) Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	if friendship.ID.IsZero() {
		return nil, utils.ErrMissingFriendshipID
	}

	filter := bson.M{"_id": friendship.ID}
	update := bson.M{"$set": friendship}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedFriendship model.Friendship
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedFriendship)
	if err != nil {
		// Again, check if the doc was not found
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFriendshipNotFound
		}
		return nil, err
	}

	return &updatedFriendship, nil
}

// Delete deletes a Friendship by its ObjectID and returns the deleted document.
func (r *friendshipRepository) Delete(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	filter := bson.M{"_id": id}

	var deletedFriendship model.Friendship
	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&deletedFriendship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFriendshipNotFound
		}
		return nil, err
	}

	return &deletedFriendship, nil
}

// Find finds all Friendships matching the given filter, with optional limit and offset.
func (r *friendshipRepository) Find(ctx context.Context, filter interface{}, limit *int, offset *int) ([]*model.Friendship, error) {
	findOpts := utils.BuildFindOptions(limit, offset)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendships []*model.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	return friendships, nil
}

// FindOne finds a single Friendship matching the given filter.
func (r *friendshipRepository) FindOne(ctx context.Context, filter interface{}) (*model.Friendship, error) {
	var friendship model.Friendship
	err := r.collection.FindOne(ctx, filter).Decode(&friendship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &friendship, nil
}

// FindOneAndUpdate finds a single Friendship that matches `filter`, applies `update`,
func (r *friendshipRepository) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*model.Friendship, error) {

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedFriendship model.Friendship
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedFriendship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFriendshipNotFound
		}
		return nil, err
	}

	return &updatedFriendship, nil
}

// FindOneAndDelete finds a single Friendship matching `filter`, deletes it,
func (r *friendshipRepository) FindOneAndDelete(ctx context.Context, filter interface{}) (*model.Friendship, error) {
	var deletedFriendship model.Friendship
	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&deletedFriendship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFriendshipNotFound
		}
		return nil, err
	}
	return &deletedFriendship, nil
}
