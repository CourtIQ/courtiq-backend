package database

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureUserIndexes creates the necessary indexes for the users collection
func EnsureUserIndexes(ctx context.Context, mdb *db.MongoDB) error {
	// Define indexes for users collection
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetSparse(true). // Allow documents without username field
				SetName("unique_username"),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_email"),
		},
		{
			Keys: bson.D{
				{Key: "displayName", Value: "text"},
				{Key: "username", Value: "text"},
			},
			Options: options.Index().
				SetName("text_search_user"),
		},
		{
			Keys: bson.D{{Key: "location.city", Value: 1}},
			Options: options.Index().SetName("location_city"),
		},
		{
			Keys: bson.D{{Key: "lastUpdated", Value: -1}},
			Options: options.Index().SetName("last_updated"),
		},
	}

	// Create indexes
	err := mdb.EnsureIndexes(ctx, db.UsersCollection, userIndexes)
	if err != nil {
		log.Printf("Failed to create user indexes: %v", err)
		return err
	}

	log.Println("User indexes created successfully")
	return nil
}