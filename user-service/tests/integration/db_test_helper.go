package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// TestDBName is the name of the test database
const TestDBName = "courtiq-test-db"

// setupTestDB connects to MongoDB and returns a client
func setupTestDB(t *testing.T) *db.MongoDB {
	// Load .env.test if it exists
	_ = godotenv.Load("../../.env.test")

	// Get MongoDB URI from environment variable with fallback
	mongoURI := os.Getenv("TEST_MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI") // Fall back to regular URI
		if mongoURI == "" {
			t.Fatal("MongoDB URI not set in environment variables")
		}
	}

	// Override the database name to use test database
	config := db.DefaultMongoDBConfig()
	config.URI = mongoURI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodb, err := db.NewMongoDBWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return mongodb
}

// cleanupCollection deletes all documents from a collection
func cleanupCollection(t *testing.T, mongodb *db.MongoDB, collectionName string) {
	collection := mongodb.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Logf("Warning: Failed to clean up collection %s: %v", collectionName, err)
	}
}

// disconnectDB closes the database connection
func disconnectDB(t *testing.T, mongodb *db.MongoDB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mongodb.Close(ctx)
	if err != nil {
		t.Logf("Warning: Failed to disconnect from MongoDB: %v", err)
	}
}
