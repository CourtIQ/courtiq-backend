package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMongoDBConfig(t *testing.T) {
	config := DefaultMongoDBConfig()
	
	assert.Equal(t, 10*time.Second, config.ConnectTimeout)
	assert.Equal(t, uint64(100), config.MaxPoolSize)
	assert.Empty(t, config.URI)
}

func TestNewMongoDB(t *testing.T) {
	// Skip this test if no MongoDB URI is provided
	mongoURI := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	mongodb, err := NewMongoDB(ctx, mongoURI)
	if err != nil {
		t.Skip("Skipping test due to MongoDB connection error:", err)
		return
	}
	
	// Make sure to close the connection
	defer mongodb.Close(context.Background())
	
	// Verify the connection
	assert.NotNil(t, mongodb)
	assert.NotNil(t, mongodb.client)
	assert.NotNil(t, mongodb.db)
	assert.Equal(t, DatabaseName, mongodb.db.Name())
}

func TestGetCollection(t *testing.T) {
	// Skip this test if no MongoDB URI is provided
	mongoURI := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	mongodb, err := NewMongoDB(ctx, mongoURI)
	if err != nil {
		t.Skip("Skipping test due to MongoDB connection error:", err)
		return
	}
	
	// Make sure to close the connection
	defer mongodb.Close(context.Background())
	
	// Test getting a collection
	collection := mongodb.GetCollection(UsersCollection)
	assert.NotNil(t, collection)
	assert.Equal(t, UsersCollection, collection.Name())
}