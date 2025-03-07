package unit

import (
	"context"
	"testing"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/services"
	"github.com/CourtIQ/courtiq-backend/user-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUser(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockRepository)
	userService := services.NewUserService(mockRepo)

	userID := primitive.NewObjectID()
	expectedUser := &model.User{
		ID:         userID,
		Email:      "test@example.com",
		FirebaseID: "firebase123",
	}

	// Configure mock
	mockRepo.On("FindByID", mock.Anything, userID).Return(expectedUser, nil)

	// Execute
	result, err := userService.GetUser(context.Background(), userID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockRepository)
	userService := services.NewUserService(mockRepo)

	userID := primitive.NewObjectID()

	// Configure mock to return not found error
	mockRepo.On("FindByID", mock.Anything, userID).Return(nil, errors.ErrNotFound)

	// Execute
	result, err := userService.GetUser(context.Background(), userID)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.IsNotFoundError(err))
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockRepository)
	userService := services.NewUserService(mockRepo)

	userID := primitive.NewObjectID()
	ctx := mocks.ContextWithMongoID(userID)

	firstName := "John"
	lastName := "Doe"
	bio := "Tennis enthusiast"

	input := &model.UpdateUserInput{
		FirstName: &firstName,
		LastName:  &lastName,
		Bio:       &bio,
	}

	expectedDisplayName := "John Doe"
	expectedUser := &model.User{
		ID:          userID,
		Email:       "test@example.com",
		FirebaseID:  "firebase123",
		FirstName:   &firstName,
		LastName:    &lastName,
		DisplayName: &expectedDisplayName,
		Bio:         &bio,
	}

	// Configure mock
	mockRepo.On("FindOneAndUpdate",
		mock.Anything,
		bson.M{"_id": userID},
		mock.Anything,
		mock.Anything).Return(expectedUser, nil)

	// Execute
	result, err := userService.UpdateUser(ctx, input)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestIsUsernameAvailable_Available(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockRepository)
	userService := services.NewUserService(mockRepo)

	username := "available_username"

	// Configure mock
	mockRepo.On("Count", mock.Anything, bson.M{"username": username}).Return(int64(0), nil)

	// Execute
	available, err := userService.IsUsernameAvailable(context.Background(), username)

	// Verify
	assert.NoError(t, err)
	assert.True(t, available)
	mockRepo.AssertExpectations(t)
}

func TestIsUsernameAvailable_Taken(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockRepository)
	userService := services.NewUserService(mockRepo)

	username := "taken_username"

	// Configure mock
	mockRepo.On("Count", mock.Anything, bson.M{"username": username}).Return(int64(1), nil)

	// Execute
	available, err := userService.IsUsernameAvailable(context.Background(), username)

	// Verify
	assert.NoError(t, err)
	assert.False(t, available)
	mockRepo.AssertExpectations(t)
}
