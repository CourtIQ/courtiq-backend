package services

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
)

// userService is the struct that implements the UserService interface.
type userService struct{}

// Helper to create a "not implemented" error with the function name.
func notImplemented(funcName string) error {
	return errors.New(funcName + " not implemented")
}

// NewUserService creates a new instance of userService.
func NewUserService() UserService {
	return &userService{}
}

// Me retrieves the profile of the currently authenticated user.
func (s *userService) Me(ctx context.Context) (*model.User, error) {
	return nil, notImplemented("Me")
}

// GetUser retrieves a user's profile by their unique ID.
func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return nil, notImplemented("GetUser")
}

// UpdateUser updates a user's profile based on the input data.
func (s *userService) UpdateUser(ctx context.Context, input *model.UpdateUserInput) (*model.User, error) {
	return nil, notImplemented("UpdateUser")
}
