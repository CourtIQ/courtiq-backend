// services/user_service.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/user-service/db"
	"github.com/CourtIQ/courtiq-backend/user-service/interfaces"
	"github.com/CourtIQ/courtiq-backend/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(mongodb *db.MongoDB) interfaces.UserService {
	return &UserService{
		collection: mongodb.GetCollection(db.UsersCollection),
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id primitive.ObjectID, updates *models.User) (*models.User, error) {
	// Create a new context with test user ID
	testID := "673fcc1444a88ee43696e40b"
	ctx = context.WithValue(ctx, "user_id", testID)

	// Get ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Convert context ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	updates.UpdatedAt = time.Now()

	update := bson.M{
		"$set": updates,
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(ctx, objectID)
}
