package services

import (
	"context"
	"time"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct{}

// FindUserByID returns a dummy user by ID
func (s *UserServiceImpl) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "test-user-id" {
		return &model.User{
			ID:          "test-user-ab",
			Username:    strPtr("testuser"),
			DisplayName: strPtr("Test User"),
			Email:       "testuser@example.com",
			Gender:      strPtr("Male"),
			Nationality: strPtr("Testland"),
			Dob:         strPtr("1990-01-01"),
			ProfileImage: &model.ProfileImage{
				Small:  "small-url",
				Medium: "medium-url",
				Large:  "large-url",
			},
			CreatedAt:   "2023-01-01T00:00:00Z",
			LastUpdated: strPtr(time.Now().Format(time.RFC3339)),
		}, nil
	}
	return nil, nil // User not found
}

// UpdateUser simulates updating a user
func (s *UserServiceImpl) UpdateUser(ctx context.Context, input model.UserUpdateInput) (*model.User, error) {
	return &model.User{
		ID:          input.ID,
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Email:       *input.Email,
		Gender:      input.Gender,
		Nationality: input.Nationality,
		Dob:         input.Dob,
		ProfileImage: &model.ProfileImage{
			Small:  "updated-small-url",
			Medium: "updated-medium-url",
			Large:  "updated-large-url",
		},
		CreatedAt:   "2023-01-01T00:00:00Z",
		LastUpdated: strPtr(time.Now().Format(time.RFC3339)),
	}, nil
}

// DeleteUser simulates user deletion
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) (bool, error) {
	if id == "test-user-id" {
		return true, nil
	}
	return false, nil
}

// IsUsernameAvailable checks if a username is available
func (s *UserServiceImpl) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	if username == "takenusername" {
		return false, nil
	}
	return true, nil
}

// FetchCurrentUser returns a dummy logged-in user
func (s *UserServiceImpl) FetchCurrentUser(ctx context.Context) (*model.User, error) {
	return &model.User{
		ID:          "my-user-id",
		Username:    strPtr("myusername"),
		DisplayName: strPtr("My User"),
		Email:       "myemail@example.com",
		Gender:      strPtr("Female"),
		Nationality: strPtr("Myland"),
		Dob:         strPtr("1995-05-05"),
		ProfileImage: &model.ProfileImage{
			Small:  "small-url",
			Medium: "medium-url",
			Large:  "large-url",
		},
		CreatedAt:   "2023-01-01T00:00:00Z",
		LastUpdated: strPtr("2023-09-01T00:00:00Z"),
	}, nil
}

// Helper function to return *string
func strPtr(s string) *string {
	return &s
}
