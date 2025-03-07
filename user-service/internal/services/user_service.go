package services

import (
	"context"
	"fmt"
	"time"

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userService struct {
	userRepo repository.Repository[model.User]
}

// NewUserService constructs a userService with the shared repository.
func NewUserService(userRepo repository.Repository[model.User]) UserServiceIntf {
	return &userService{
		userRepo: userRepo,
	}
}

// Me retrieves the profile of the currently authenticated user.
func (s *userService) Me(ctx context.Context) (*model.User, error) {
	// Extract mongoId from context
	mongoID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, sharedErrors.WrapError(sharedErrors.ErrUnauthorized, "could not retrieve user id from context")
	}

	user, err := s.userRepo.FindByID(ctx, mongoID)
	if err != nil {
		if sharedErrors.IsNotFoundError(err) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "error retrieving user")
	}

	return user, nil
}

// GetUser retrieves a user's profile by their unique ID.
func (s *userService) GetUser(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if sharedErrors.IsNotFoundError(err) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "error retrieving user by ID")
	}

	return user, nil
}

// UpdateUser updates a user's profile based on the input data.
func (s *userService) UpdateUser(ctx context.Context, input *model.UpdateUserInput) (*model.User, error) {
	// Extract user id from context
	mongoID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, sharedErrors.WrapError(sharedErrors.ErrUnauthorized, "could not retrieve user id from context")
	}

	// Create a map for fields to update
	updateFields := bson.M{}

	// Conditionally add fields if they are provided
	if input.FirstName != nil && input.LastName != nil {
		updateFields["firstName"] = *input.FirstName
		updateFields["lastName"] = *input.LastName
		updateFields["displayName"] = fmt.Sprintf("%s %s", *input.FirstName, *input.LastName)
	}

	if input.DateOfBirth != nil {
		updateFields["dateOfBirth"] = *input.DateOfBirth
	}
	if input.Bio != nil {
		updateFields["bio"] = *input.Bio
	}
	if input.Username != nil {
		updateFields["username"] = *input.Username
	}
	if input.Gender != nil {
		updateFields["gender"] = *input.Gender
	}

	if input.FcmTokens != nil {
		var tokens []string
		for _, tokenPtr := range input.FcmTokens {
			if tokenPtr != nil {
				tokens = append(tokens, *tokenPtr)
			}
		}
		updateFields["fcmTokens"] = tokens
	}

	updateFields["lastUpdated"] = primitive.NewDateTimeFromTime(time.Now())

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

	// If no fields to update, just retrieve and return the current user
	if len(updateFields) == 0 {
		return s.GetUser(ctx, mongoID)
	}

	// Use the repository's FindOneAndUpdate method
	filter := bson.M{"_id": mongoID}
	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updatedUser, err := s.userRepo.FindOneAndUpdate(ctx, filter, update, opts)
	if err != nil {
		if sharedErrors.IsNotFoundError(err) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "error updating user")
	}

	return updatedUser, nil
}

// IsUsernameAvailable checks if the given username is available.
func (s *userService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}

	// Check if the username is already taken
	count, err := s.userRepo.Count(ctx, filter)
	if err != nil {
		return false, sharedErrors.WrapError(err, "error checking username availability")
	}

	// If the count is 0, the username is available
	return count == 0, nil
}
