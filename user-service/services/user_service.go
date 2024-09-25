package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // For ObjectId conversion
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// strPtr is a utility function to return a pointer to a string
func strPtr(s string) *string {
	return &s
}

// UserService implements the UserServiceProvider interface
type UserService struct {
	Collection *mongo.Collection
}

// NewUserService initializes UserService with the users collection
func NewUserService() *UserService {
	collection := MongoClient.Database("prod-db").Collection("users")
	return &UserService{
		Collection: collection,
	}
}

// FindUserByID retrieves a user by ObjectId from MongoDB
func (s *UserService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	// Convert string id to MongoDB ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectId format: %v", err)
	}

	var user model.User
	err = s.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no user found with id %s", id)
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user in MongoDB
func (s *UserService) UpdateUser(ctx context.Context, input model.UserUpdateInput) (*model.User, error) {
	// Convert string id to MongoDB ObjectId
	objectId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectId format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"username":    input.Username,
			"displayName": input.DisplayName,
			"email":       input.Email,
			"gender":      input.Gender,
			"nationality": input.Nationality,
			"dob":         input.Dob,
			"lastUpdated": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser model.User
	err = s.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no user found with id %s", input.ID)
	} else if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser deletes a user from MongoDB
func (s *UserService) DeleteUser(ctx context.Context, id string) (bool, error) {
	// Convert string id to MongoDB ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid ObjectId format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	result, err := s.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

// IsUsernameAvailable checks if a username is available in MongoDB
func (s *UserService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := s.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// FetchCurrentUser retrieves a mock or authenticated current user
func (s *UserService) FetchCurrentUser(ctx context.Context) (*model.User, error) {
	// For now, returning a mock user for current user
	return &model.User{
		ID:          "current-user-id",
		Username:    strPtr("current_user"),
		DisplayName: strPtr("Current User"),
		Email:       "currentuser@example.com",
		Gender:      strPtr("Male"),
		Nationality: strPtr("Testland"),
		Dob:         strPtr("1990-01-01"),
		ProfileImage: &model.ProfileImage{
			Small:  "small-url",
			Medium: "medium-url",
			Large:  "large-url",
		},
		CreatedAt:   time.Now().Format(time.RFC3339),         // Keep as string
		LastUpdated: strPtr(time.Now().Format(time.RFC3339)), // Use strPtr here
	}, nil
}
