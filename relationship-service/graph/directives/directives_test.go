// graph/directives/directives_test.go
package directives_test

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/directives"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	mockrepo "github.com/CourtIQ/courtiq-backend/relationship-service/tests/mocks/repository"
)

// MockGetCurrentUserID simulates retrieving the current user ID.
func MockGetCurrentUserID(userID string, err error) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		return userID, err
	}
}

func withFieldContextArgs(ctx context.Context, args map[string]interface{}) context.Context {
	fc := &graphql.FieldContext{Args: args}
	return graphql.WithFieldContext(ctx, fc)
}

func mockNextResolver(returnVal interface{}, err error) graphql.Resolver {
	return func(ctx context.Context) (interface{}, error) {
		return returnVal, err
	}
}

func TestSatisfiesDirective(t *testing.T) {
	baseCtx := context.Background()
	currentUserID := "user123"
	friendshipID := "507f1f77bcf86cd799439011"
	receiverUserID := "user456"

	t.Run("unauthorized if GetCurrentUserID fails", func(t *testing.T) {
		t.Log("Test scenario: Unauthorized if no user ID can be retrieved")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo

		t.Log("Simulating GetCurrentUserID failure")
		directives.GetCurrentUserID = MockGetCurrentUserID("", errors.New("no user"))

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{"friendshipId": friendshipID})

		t.Log("Calling SatisfiesDirective expecting unauthorized error")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("nextCalled", nil),
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "unauthorized")
	})

	t.Run("noExistingFriendship requires receiverId", func(t *testing.T) {
		t.Log("Test scenario: noExistingFriendship directive requires receiverId")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{})

		t.Log("Calling SatisfiesDirective without receiverId expecting an error")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("nextCalled", nil),
			nil,
			nil,
			nil,
			nil,
			nil,
			&[]bool{true}[0],
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "noExistingFriendship requires receiverId")
	})

	t.Run("noExistingFriendship fails if friendship exists", func(t *testing.T) {
		t.Log("Test scenario: noExistingFriendship should fail if a friendship already exists")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)

		filter := map[string]interface{}{
			"type": "FRIENDSHIP",
			"participantIds": map[string]interface{}{
				"$all": []string{currentUserID, receiverUserID},
			},
			"status": map[string]interface{}{
				"$in": []string{"PENDING", "ACTIVE"},
			},
		}

		t.Logf("Setting repo expectation: Count with filter: %#v", filter)
		repo.On("Count", mock.Anything, filter).Return(int64(1), nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{
			"receiverId": receiverUserID,
		})

		t.Log("Calling SatisfiesDirective expecting forbidden error due to existing friendship")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("nextCalled", nil),
			nil,
			nil,
			nil,
			nil,
			nil,
			&[]bool{true}[0],
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "forbidden: existing friendship or pending request already present")
		repo.AssertExpectations(t)
	})

	t.Run("noExistingFriendship passes if no friendship exists", func(t *testing.T) {
		t.Log("Test scenario: noExistingFriendship should pass if no friendship exists")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)

		filter := map[string]interface{}{
			"type": "FRIENDSHIP",
			"participantIds": map[string]interface{}{
				"$all": []string{currentUserID, receiverUserID},
			},
			"status": map[string]interface{}{
				"$in": []string{"PENDING", "ACTIVE"},
			},
		}

		t.Logf("Setting repo expectation: Count with filter: %#v", filter)
		repo.On("Count", mock.Anything, filter).Return(int64(0), nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{
			"receiverId": receiverUserID,
		})

		t.Log("Calling SatisfiesDirective expecting success and nextCalled result")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("nextCalled", nil),
			nil,
			nil,
			nil,
			nil,
			nil,
			&[]bool{true}[0],
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.NoError(t, err)
		assert.Equal(t, "nextCalled", result)
		repo.AssertExpectations(t)
	})

	t.Run("missing friendshipId for requireParticipant check", func(t *testing.T) {
		t.Log("Test scenario: requireParticipant but friendshipId missing")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{})

		participant := true
		t.Log("Calling SatisfiesDirective expecting error due to missing friendshipId")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("shouldNotBeCalled", nil),
			nil,
			nil,
			&participant,
			nil,
			nil,
			nil,
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "missing friendshipId")
	})

	t.Run("forbidden if conditions not met", func(t *testing.T) {
		t.Log("Test scenario: conditions not met, should return forbidden")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)
		participant := true

		t.Log("Setting repo expectation using mock.MatchedBy to ignore exact _id type")
		repo.On("Count", mock.Anything, mock.MatchedBy(func(filter interface{}) bool {
			t.Logf("Matcher sees filter: %#v", filter)
			f, ok := filter.(map[string]interface{})
			if !ok {
				return false
			}
			// Check participantIds
			if f["participantIds"] != currentUserID {
				t.Log("participantIds does not match currentUserID")
				return false
			}
			// Ensure _id is present
			if _, hasID := f["_id"]; !hasID {
				t.Log("No _id field found")
				return false
			}
			return true
		})).Return(int64(0), nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{
			"friendshipId": friendshipID,
		})

		t.Log("Calling SatisfiesDirective expecting forbidden error")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("shouldNotBeCalled", nil),
			nil,
			nil,
			&participant,
			nil,
			nil,
			nil,
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "forbidden: conditions not met")
		repo.AssertExpectations(t)
	})

	t.Run("conditions met should proceed", func(t *testing.T) {
		t.Log("Test scenario: conditions are met, should proceed to next")
		repo := &mockrepo.MockRelationshipRepository{}
		directives.RelationshipRepo = repo
		directives.GetCurrentUserID = MockGetCurrentUserID(currentUserID, nil)

		participant := true
		status := model.RelationshipStatusActive
		rType := model.RelationshipTypeFriendship

		t.Log("Setting repo expectation using mock.MatchedBy to handle _id")
		repo.On("Count", mock.Anything, mock.MatchedBy(func(filter interface{}) bool {
			t.Logf("Matcher sees filter: %#v", filter)
			f, ok := filter.(map[string]interface{})
			if !ok {
				return false
			}
			if f["status"] != "ACTIVE" {
				t.Log("status not ACTIVE")
				return false
			}
			if f["type"] != "FRIENDSHIP" {
				t.Log("type not FRIENDSHIP")
				return false
			}
			if f["participantIds"] != currentUserID {
				t.Log("participantIds not user123")
				return false
			}
			if _, hasID := f["_id"]; !hasID {
				t.Log("No _id field found")
				return false
			}
			return true
		})).Return(int64(1), nil)

		ctx := withFieldContextArgs(baseCtx, map[string]interface{}{
			"friendshipId": friendshipID,
		})

		t.Log("Calling SatisfiesDirective expecting success and nextCalled")
		result, err := directives.SatisfiesDirective(
			ctx,
			nil,
			mockNextResolver("nextCalled", nil),
			&status,
			&rType,
			&participant,
			nil,
			nil,
			nil,
		)

		t.Logf("Result: %v, Error: %v", result, err)
		assert.NoError(t, err)
		assert.Equal(t, "nextCalled", result)
		repo.AssertExpectations(t)
	})
}
