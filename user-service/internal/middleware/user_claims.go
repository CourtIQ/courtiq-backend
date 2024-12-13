package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/CourtIQ/courtiq-backend/user-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Context key type to avoid collisions
type contextKey string

const mongoIDKey contextKey = "mongoId"

// WithUserClaims is an HTTP middleware that reads the "X-User-Claims" header,
// extracts the mongoId, and places it into the request context.
func WithUserClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mongoID, err := utils.GetMongoID(r)
		if err != nil {
			// If we cannot get mongoId, we can log it or handle as needed.
			// For now, we will just proceed without a mongoId in context.
			// If needed, you can return an HTTP error or leave ctx unchanged.
		} else {
			// Store mongoId in context
			ctx := context.WithValue(r.Context(), mongoIDKey, mongoID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// GetMongoIDFromContext retrieves the mongoId from the context.
// Returns empty string if not found.
func GetMongoIDFromContext(ctx context.Context) (primitive.ObjectID, error) {
	v := ctx.Value(mongoIDKey)
	if v == nil {
		return primitive.NilObjectID, errors.New("no mongoId in context")
	}

	mongoIDStr, ok := v.(string)
	if !ok {
		return primitive.NilObjectID, errors.New("mongoId in context is not a string")
	}

	oid, err := primitive.ObjectIDFromHex(mongoIDStr)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return oid, nil
}
