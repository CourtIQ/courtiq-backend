package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/db"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// friendshipRepository is the concrete implementation of the FriendshipRepository interface.
type friendshipRepository struct {
	coll *mongo.Collection
}

// NewFriendshipRepository creates a new NewFriendshipRepository implementation
func NewFriendshipRepository(mdb *db.MongoDB) FriendshipRepository {
	return &friendshipRepository{
		coll: mdb.GetCollection(db.FriendshipsCollection),
	}
}

func (r *friendshipRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	var friendship model.Friendship
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&friendship)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *friendshipRepository) GetMyFriends(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Friendship, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"senderId": userID},
			{"receiverId": userID},
		},
		"status": "ACTIVE",
	}
	findOpts := utils.BuildFindOptions(limit, offset)
	cursor, err := r.coll.Find(ctx, filter, findOpts)
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

func (r *friendshipRepository) GetFriendsOfUser(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Friendship, error) {
	// This is essentially the same logic as GetMyFriends, just used to fetch another user's friends.
	filter := bson.M{
		"$or": []bson.M{
			{"senderId": userID},
			{"receiverId": userID},
		},
		"status": "ACTIVE",
	}
	findOpts := utils.BuildFindOptions(limit, offset)
	cursor, err := r.coll.Find(ctx, filter, findOpts)
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

func (r *friendshipRepository) GetMyFriendRequests(ctx context.Context, userID string) ([]*model.Friendship, error) {
	// Requests received by this user
	filter := bson.M{
		"receiverId": userID,
		"status":     "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *friendshipRepository) GetSentFriendRequests(ctx context.Context, userID string) ([]*model.Friendship, error) {
	// Requests sent by this user
	filter := bson.M{
		"senderId": userID,
		"status":   "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *friendshipRepository) Insert(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	friendship.ID = primitive.NewObjectID()
	_, err := r.coll.InsertOne(ctx, friendship)
	if err != nil {
		return nil, err
	}
	return friendship, nil
}

func (r *friendshipRepository) Update(ctx context.Context, friendship *model.Friendship) (*model.Friendship, error) {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": friendship.ID}, bson.M{"$set": friendship})
	if err != nil {
		return nil, err
	}
	return friendship, nil
}

func (r *friendshipRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
