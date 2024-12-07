// internal/db/mongodb/indexes.go
package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getRelationshipsIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "participants", Value: 1},
				{Key: "type", Value: 1},
				{Key: "status", Value: 1},
			},
			// Additional index options can be set here if needed
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: -1}},
			// Additional index options can be set here if needed
		},
	}
}
