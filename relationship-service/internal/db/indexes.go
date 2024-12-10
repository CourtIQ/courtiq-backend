// internal/db/indexes.go
package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureRelationshipsIndexes(ctx context.Context, coll *mongo.Collection) error {
	// Define the indexes
	models := getRelationshipsIndexes()

	// Create them
	_, err := coll.Indexes().CreateMany(ctx, models, options.CreateIndexes())
	return err
}

func getRelationshipsIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "participantIds", Value: 1},
				{Key: "type", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
	}
}
