package integration

import (
	"context"
	"testing"
	"time"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/services"
	"github.com/CourtIQ/courtiq-backend/user-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIntegration_UserService(t *testing.T) {
	// Skip if short mode is enabled (for quick test runs)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	mongodb := setupTestDB(t)
	defer disconnectDB(t, mongodb)

	// Clean up test collection
	cleanupCollection(t, mongodb, db.UsersCollection)

	// Create repository and service
	userRepo := repository.NewBaseRepository[model.User](mongodb.GetCollection(db.UsersCollection))
	userService := services.NewUserService(userRepo)

	// Create test user
	ctx := context.Background()
	testUser := createTestUser(t, ctx, userRepo)

	// Run tests
	t.Run("GetUser", func(t *testing.T) {
		user, err := userService.GetUser(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Setup update data
		firstName := "Updated"
		lastName := "User"
		bio := "Updated bio"

		input := &model.UpdateUserInput{
			FirstName: &firstName,
			LastName:  &lastName,
			Bio:       &bio,
		}

		// Update user
		mockCtx := mocks.ContextWithMongoID(testUser.ID)
		updatedUser, err := userService.UpdateUser(mockCtx, input)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, updatedUser.ID)
		assert.Equal(t, *input.FirstName, *updatedUser.FirstName)
		assert.Equal(t, *input.LastName, *updatedUser.LastName)
		assert.Equal(t, *input.Bio, *updatedUser.Bio)
		assert.Equal(t, "Updated User", *updatedUser.DisplayName)
	})

	t.Run("IsUsernameAvailable", func(t *testing.T) {
		// Test username that should be available
		available, err := userService.IsUsernameAvailable(ctx, "available_username")
		require.NoError(t, err)
		assert.True(t, available)

		// Set a username for our test user
		username := "taken_username"
		_, err = userRepo.FindOneAndUpdate(
			ctx,
			bson.M{"_id": testUser.ID},
			bson.M{"$set": bson.M{"username": username}},
		)
		require.NoError(t, err)

		// Test username that should be taken
		available, err = userService.IsUsernameAvailable(ctx, username)
		require.NoError(t, err)
		assert.False(t, available)
	})
}

// Helper to create a test user
func createTestUser(t *testing.T, ctx context.Context, repo repository.Repository[model.User]) *model.User {
	now := time.Now()
	displayName := "Test User"

	user := &model.User{
		ID:          primitive.NewObjectID(),
		FirebaseID:  "test-firebase-id",
		Email:       "test@example.com",
		DisplayName: &displayName,
		CreatedAt:   &now,
		LastUpdated: &now,
	}

	insertedUser, err := repo.Insert(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, insertedUser)

	return insertedUser
}
