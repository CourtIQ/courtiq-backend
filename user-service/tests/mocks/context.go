package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Context key type to avoid collisions
type contextKey string

const mongoIDKey contextKey = "mongoId"

// ContextWithMongoID returns a context with the mongoId value set
func ContextWithMongoID(id primitive.ObjectID) context.Context {
	return context.WithValue(context.Background(), mongoIDKey, id.Hex())
}