package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/middleware"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService constructs a userService with the given UserRepository.
func NewUserService(userRepo repository.UserRepository) UserServiceIntf {
	return &userService{
		userRepo: userRepo,
	}
}

// Me retrieves the profile of the currently authenticated user.
func (s *userService) Me(ctx context.Context) (*model.User, error) {
	// Extract mongoId from context (already validated in middleware or utils)
	mongoID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: could not retrieve user id from context: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, mongoID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUser retrieves a user's profile by their unique ID.
func (s *userService) GetUser(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user by ID: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser updates a user's profile based on the input data.
func (s *userService) UpdateUser(ctx context.Context, input *model.UpdateUserInput) (*model.User, error) {
	// Extract user id from context
	mongoID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: could not retrieve user id from context: %w", err)
	}

	// Attempt to update the user in the repository
	updatedUser, err := s.userRepo.UpdateUser(ctx, mongoID, input)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}
	if updatedUser == nil {
		// If nil is returned, it might mean no user was found to update
		return nil, errors.New("user not found")
	}

	return updatedUser, nil
}

// IsUsernameAvailable checks if the given username is available.
func (s *userService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	// Check if the username is available
	available, err := s.userRepo.Count(ctx, username)
	if err != nil {
		return false, fmt.Errorf("error checking username availability: %w", err)
	}

	// If the count is 0, the username is available
	isAvailable := available == 0

	return isAvailable, nil
}
