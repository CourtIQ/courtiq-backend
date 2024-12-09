// filters_test.go
package satisfies

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

// Mock function to simulate user retrieval from context
func mockUserIDFromContext(ctx context.Context) (string, error) {
	uid, ok := ctx.Value("mockUserID").(string)
	if !ok || uid == "" {
		return "", errors.New("no mock user id in context")
	}
	return uid, nil
}

func TestBuildRoleFilter(t *testing.T) {
	// Override the getUserIDFromContext function for the duration of the test
	originalGetUser := getUserIDFromContext
	getUserIDFromContext = mockUserIDFromContext
	defer func() { getUserIDFromContext = originalGetUser }()

	ctx := context.WithValue(context.Background(), "mockUserID", "currentUser123")
	fc := &graphql.FieldContext{
		Args: map[string]interface{}{
			"ofUserId": "otherUser456",
		},
	}
	ctx = graphql.WithFieldContext(ctx, fc)

	// Test nil role
	filter, err := BuildRoleFilter(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, filter)

	// With participants required
	role := &model.RoleConditions{
		RequireParticipants: boolPtr(true),
	}
	filter, err = BuildRoleFilter(ctx, role)
	require.NoError(t, err)
	require.Equal(t, []string{"currentUser123", "otherUser456"}, filter["participantIds"])

	// RequireSender
	role = &model.RoleConditions{
		RequireSender: boolPtr(true),
	}
	filter, err = BuildRoleFilter(ctx, role)
	require.NoError(t, err)
	require.Equal(t, "currentUser123", filter["senderId"])

	// Multiple conditions
	role = &model.RoleConditions{
		RequireReceiver:     boolPtr(true),
		RequireParticipants: boolPtr(true),
	}
	filter, err = BuildRoleFilter(ctx, role)
	require.NoError(t, err)
	require.Equal(t, "currentUser123", filter["receiverId"])
	require.Equal(t, []string{"currentUser123", "otherUser456"}, filter["participantIds"])
}

func TestBuildExistenceFilter(t *testing.T) {
	// Since BuildExistenceFilter does not depend on user retrieval,
	// we do not need to mock getUserIDFromContext here.

	ctx := context.Background()

	// No existence conditions
	filter, err := BuildExistenceFilter(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, filter)

	// FRIENDSHIP with friendshipId
	fc := &graphql.FieldContext{
		Args: map[string]interface{}{
			"friendshipId": primitive.NewObjectID().Hex(),
		},
	}
	ctx = graphql.WithFieldContext(ctx, fc)

	existence := &model.ExistenceConditions{
		RelationshipType:   relTypePtr(model.RelationshipTypeFriendship),
		RelationshipStatus: relStatusPtr(model.RelationshipStatusPending),
	}
	filter, err = BuildExistenceFilter(ctx, existence)
	require.NoError(t, err)
	require.Equal(t, "FRIENDSHIP", filter["type"])
	require.Equal(t, "PENDING", filter["status"])
	require.NotNil(t, filter["_id"])

	// FRIENDSHIP without friendshipId
	fc = &graphql.FieldContext{
		Args: map[string]interface{}{}, // no friendshipId
	}
	ctx = graphql.WithFieldContext(context.Background(), fc)
	filter, err = BuildExistenceFilter(ctx, existence)
	require.NoError(t, err)
	require.Equal(t, "FRIENDSHIP", filter["type"])
	require.Equal(t, "PENDING", filter["status"])
	_, hasId := filter["_id"]
	require.False(t, hasId)

	// COACHSHIP with invalid coachshipId
	fc = &graphql.FieldContext{
		Args: map[string]interface{}{
			"coachshipId": "invalidHex",
		},
	}
	ctx = graphql.WithFieldContext(context.Background(), fc)
	existence.RelationshipType = relTypePtr(model.RelationshipTypeCoachship)
	filter, err = BuildExistenceFilter(ctx, existence)
	require.NoError(t, err)
	require.Equal(t, "COACHSHIP", filter["type"])
	require.Equal(t, "PENDING", filter["status"])
	_, hasId = filter["_id"]
	require.False(t, hasId)
}

func TestBuildNonExistenceFilter(t *testing.T) {
	// Mock user retrieval
	originalGetUser := getUserIDFromContext
	getUserIDFromContext = mockUserIDFromContext
	defer func() { getUserIDFromContext = originalGetUser }()

	ctx := context.WithValue(context.Background(), "mockUserID", "currentUser123")

	// No conditions
	filter, err := BuildNonExistenceFilter(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, filter)

	// NoExistingFriendship
	fc := &graphql.FieldContext{
		Args: map[string]interface{}{
			"receiverId": "receiverUser456",
		},
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	nonExistence := &model.NonExistenceConditions{
		NoExistingFriendship: boolPtr(true),
	}
	filter, err = BuildNonExistenceFilter(ctx, nonExistence)
	require.NoError(t, err)
	require.Equal(t, "FRIENDSHIP", filter["type"])
	require.Equal(t, []string{"currentUser123", "receiverUser456"}, filter["participantIds"])

	// NotExistingCoach
	fc = &graphql.FieldContext{
		Args: map[string]interface{}{
			"ofUserId": "studentUser789",
		},
	}
	ctx = graphql.WithFieldContext(context.Background(), fc)
	ctx = context.WithValue(ctx, "mockUserID", "coachUser000")
	nonExistence = &model.NonExistenceConditions{
		NotExistingCoach: boolPtr(true),
	}
	filter, err = BuildNonExistenceFilter(ctx, nonExistence)
	require.NoError(t, err)
	require.Equal(t, "COACHSHIP", filter["type"])
	require.Equal(t, "coachUser000", filter["coachId"])
	require.Equal(t, "studentUser789", filter["studentId"])

	// NotExistingStudent
	fc = &graphql.FieldContext{
		Args: map[string]interface{}{
			"ofUserId": "coachUser222",
		},
	}
	ctx = graphql.WithFieldContext(context.Background(), fc)
	ctx = context.WithValue(ctx, "mockUserID", "studentUser333")
	nonExistence = &model.NonExistenceConditions{
		NotExistingStudent: boolPtr(true),
	}
	filter, err = BuildNonExistenceFilter(ctx, nonExistence)
	require.NoError(t, err)
	require.Equal(t, "COACHSHIP", filter["type"])
	require.Equal(t, "coachUser222", filter["coachId"])
	require.Equal(t, "studentUser333", filter["studentId"])

	// Test missing required args
	fc = &graphql.FieldContext{
		Args: map[string]interface{}{},
	}
	ctx = graphql.WithFieldContext(context.Background(), fc)
	nonExistence = &model.NonExistenceConditions{
		NotExistingCoach: boolPtr(true),
	}
	filter, err = BuildNonExistenceFilter(ctx, nonExistence)
	require.Error(t, err)
	require.Nil(t, filter)
}

// Helper functions

func boolPtr(b bool) *bool {
	return &b
}

func relTypePtr(r model.RelationshipType) *model.RelationshipType {
	return &r
}

func relStatusPtr(s model.RelationshipStatus) *model.RelationshipStatus {
	return &s
}
