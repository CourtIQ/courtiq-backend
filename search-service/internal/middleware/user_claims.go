package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserClaims represents the JSON structure of the user claims
type UserClaims struct {
	MongoID       string `json:"mongoId"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	AuthTime      int64  `json:"auth_time"`
	UserID        string `json:"user_id"`
	Sub           string `json:"sub"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Firebase      struct {
		Identities struct {
			Email []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
	UID string `json:"uid"`
}

// GetMongoID extracts the mongoId from the X-User-Claims header in the given request.
// Returns the mongoId if found, otherwise an error.
func GetMongoID(r *http.Request) (string, error) {
	claimsStr := r.Header.Get("X-User-Claims")
	if claimsStr == "" {
		return "", errors.New("no user claims found in X-User-Claims header")
	}

	var claims UserClaims
	if err := json.Unmarshal([]byte(claimsStr), &claims); err != nil {
		return "", err
	}

	if claims.MongoID == "" {
		return "", errors.New("mongoId not present in user claims")
	}

	return claims.MongoID, nil
}

// Context key type to avoid collisions
type contextKey string

const mongoIDKey contextKey = "mongoId"

// WithUserClaims is an HTTP middleware that reads the "X-User-Claims" header,
// extracts the mongoId, and places it into the request context.
func WithUserClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mongoID, err := GetMongoID(r)
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
