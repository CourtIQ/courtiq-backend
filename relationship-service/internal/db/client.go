// db/client.go
package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName          = "courtiq-db"
	FriendshipsCollection = "friendships"
	CoachshipsCollection  = "coachships"
)

// MongoDB represents the MongoDB client wrapper
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	db := client.Database(DatabaseName)
	log.Printf("Connected to MongoDB database: %s", DatabaseName)

	return &MongoDB{
		client: client,
		db:     db,
	}, nil
}

// GetCollection returns a mongo.Collection for the given name
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// Close disconnects from MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
