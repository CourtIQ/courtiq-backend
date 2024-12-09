// filters.go
package satisfies

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var getUserIDFromContext = middleware.GetUserIDFromContext

func BuildRoleFilter(
	ctx context.Context,
	role *model.RoleConditions,
) (map[string]interface{}, error) {
	if role == nil {
		return nil, nil
	}

	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil, errors.New("no field context found")
	}

	currentUserID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := make(map[string]interface{})

	if role.RequireParticipants != nil && *role.RequireParticipants {
		participantIds := []string{currentUserID}
		if ofUserID, ok := fc.Args["ofUserId"].(string); ok && ofUserID != "" {
			participantIds = append(participantIds, ofUserID)
		}
		filter["participantIds"] = bson.M{"$all": participantIds}
	}

	if role.RequireSender != nil && *role.RequireSender {
		filter["senderId"] = currentUserID
	}

	if role.RequireReceiver != nil && *role.RequireReceiver {
		filter["receiverId"] = currentUserID
	}

	if role.RequireStudent != nil && *role.RequireStudent {
		filter["studentId"] = currentUserID
	}

	if role.RequireCoach != nil && *role.RequireCoach {
		filter["coachId"] = currentUserID
	}

	return filter, nil
}

// BuildExistenceFilter constructs a filter map for checking existence conditions.
func BuildExistenceFilter(
	ctx context.Context,
	existence *model.ExistenceConditions,
) (map[string]interface{}, error) {
	if existence == nil {
		return nil, nil
	}

	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil, errors.New("no field context found")
	}

	filter := make(map[string]interface{})

	if existence.RelationshipType != nil {
		if *existence.RelationshipType == model.RelationshipTypeFriendship {
			if friendshipID, ok := fc.Args["friendshipId"].(string); ok && friendshipID != "" {
				oid, err := primitive.ObjectIDFromHex(friendshipID)
				if err == nil {
					filter["_id"] = oid
				}
			}
		} else if *existence.RelationshipType == model.RelationshipTypeCoachship {
			if coachshipId, ok := fc.Args["coachshipId"].(string); ok && coachshipId != "" {
				oid, err := primitive.ObjectIDFromHex(coachshipId)
				if err == nil {
					filter["_id"] = oid
				}
			}

		}
		filter["type"] = string(*existence.RelationshipType)
	}

	if existence.RelationshipStatus != nil {
		filter["status"] = string(*existence.RelationshipStatus)
	}

	return filter, nil
}

// BuildNonExistenceFilter constructs a filter map for checking non-existence conditions.
func BuildNonExistenceFilter(
	ctx context.Context,
	nonExistence *model.NonExistenceConditions,
) (map[string]interface{}, error) {
	if nonExistence == nil {
		return nil, nil
	}

	currentUserID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil, errors.New("no field context found")
	}

	filter := make(map[string]interface{})

	if nonExistence.NoExistingFriendship != nil && *nonExistence.NoExistingFriendship {
		receiverID, ok := fc.Args["receiverId"].(string)
		if !ok || receiverID == "" {
			return nil, errors.New("receiverId argument is required and must be a non-empty string")
		}
		filter["type"] = "FRIENDSHIP"
		filter["participantIds"] = []string{currentUserID, receiverID}
	}

	if nonExistence.NotExistingCoach != nil && *nonExistence.NotExistingCoach {
		ofUserID, ok := fc.Args["ofUserId"].(string)
		if !ok || ofUserID == "" {
			return nil, errors.New("ofUserId argument is required and must be a non-empty string")
		}
		filter["type"] = "COACHSHIP"
		filter["coachId"] = currentUserID
		filter["studentId"] = ofUserID
	}

	if nonExistence.NotExistingStudent != nil && *nonExistence.NotExistingStudent {
		ofUserID, ok := fc.Args["ofUserId"].(string)
		if !ok || ofUserID == "" {
			return nil, errors.New("ofUserId argument is required and must be a non-empty string")
		}
		filter["type"] = "COACHSHIP"
		filter["coachId"] = ofUserID
		filter["studentId"] = currentUserID
	}

	return filter, nil
}
