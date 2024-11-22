// graph/resolvers/helpers.go
package resolvers

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The key used by the API gateway to pass user ID
const userIDContextKey = "user_id"

// getUserIDFromContext extracts the user ID from context
func getUserIDFromContext(ctx context.Context) (primitive.ObjectID, error) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("user ID not found in context")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid user ID format")
	}

	return objectID, nil
}
