// services/user_service.go
package services

import (
	"context"
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
	updates.UpdatedAt = time.Now()

	update := bson.M{
		"$set": updates,
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(ctx, id)
}
