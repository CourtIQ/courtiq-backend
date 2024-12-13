package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	coll *mongo.Collection
}

// NewUserRepository creates a new UserRepository backed by MongoDB.
func NewUserRepository(mdb *db.MongoDB) UserRepository {
	return &userRepository{
		coll: mdb.GetCollection(db.UsersCollection),
	}
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	filter := bson.M{"_id": id}

	var user model.User
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user found with that ID
			return nil, nil
		}
		// An unexpected error occurred
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, input *model.UpdateUserInput) (*model.User, error) {
	// Create a map for fields to update
	updateFields := bson.M{}

	// Conditionally add fields if they are provided
	if input.FirstName != nil {
		updateFields["firstName"] = *input.FirstName
	}
	if input.LastName != nil {
		updateFields["lastName"] = *input.LastName
	}
	if input.DateOfBirth != nil {
		// Assuming `DateOfBirth` is a scalar Time type that you can store as a Date in Mongo
		updateFields["dateOfBirth"] = *input.DateOfBirth
	}
	if input.Bio != nil {
		updateFields["bio"] = *input.Bio
	}
	lastUpdated := primitive.NewDateTimeFromTime(time.Now())

	updateFields["lastUpdated"] = lastUpdated

	// Handle location if provided
	if input.Location != nil {
		locUpdate := bson.M{}
		if input.Location.City != nil {
			locUpdate["city"] = *input.Location.City
		}
		if input.Location.State != nil {
			locUpdate["state"] = *input.Location.State
		}
		if input.Location.Country != nil {
			locUpdate["country"] = *input.Location.Country
		}
		if input.Location.Latitude != nil {
			locUpdate["latitude"] = *input.Location.Latitude
		}
		if input.Location.Longitude != nil {
			locUpdate["longitude"] = *input.Location.Longitude
		}

		// Only set 'location' if at least one subfield was provided
		if len(locUpdate) > 0 {
			updateFields["location"] = locUpdate
		}
	}

	if len(updateFields) == 0 {
		return r.GetByID(ctx, id)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser model.User

	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err == mongo.ErrNoDocuments {
		// No user found to update
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &updatedUser, nil
}

func (r *userRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}
