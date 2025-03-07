package testing

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupTestEnv sets up common test environment variables
func SetupTestEnv(t *testing.T) {
	t.Setenv("SERVICE_NAME", "test-service")
	t.Setenv("PORT", "8999")
	t.Setenv("GO_ENV", "test")
	t.Setenv("LOG_LEVEL", "debug")
}

// CreateTestMongoClient creates a MongoDB client for testing
// It will connect to a MongoDB instance if MONGODB_TEST_URL is set,
// otherwise it will return nil
func CreateTestMongoClient(t *testing.T) *mongo.Client {
	mongoURL := os.Getenv("MONGODB_TEST_URL")
	if mongoURL == "" {
		t.Skip("Skipping MongoDB test: MONGODB_TEST_URL not set")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	require.NoError(t, err, "Failed to connect to test MongoDB instance")

	// Ping to ensure connection is established
	err = client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping test MongoDB instance")

	return client
}

// CleanupTestMongoClient disconnects the MongoDB client
func CleanupTestMongoClient(t *testing.T, client *mongo.Client) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	require.NoError(t, err, "Failed to disconnect from test MongoDB instance")
}