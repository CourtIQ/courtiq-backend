package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents the MongoDB client wrapper
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// MongoDBConfig contains configuration for MongoDB connection
type MongoDBConfig struct {
	URI            string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
}

// DefaultMongoDBConfig returns default MongoDB configuration
func DefaultMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
	}
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(ctx context.Context, uri string) (*MongoDB, error) {
	config := DefaultMongoDBConfig()
	config.URI = uri
	return NewMongoDBWithConfig(ctx, config)
}

// NewMongoDBWithConfig creates a new MongoDB instance with the given config
func NewMongoDBWithConfig(ctx context.Context, config MongoDBConfig) (*MongoDB, error) {
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetConnectTimeout(config.ConnectTimeout).
		SetMaxPoolSize(config.MaxPoolSize)

	ctx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(DatabaseName)
	log.Printf("Connected to MongoDB database: %s", DatabaseName)

	return &MongoDB{
		client: client,
		db:     db,
	}, nil
}

// GetClient returns the MongoDB client
func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}

// GetDatabase returns the MongoDB database
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.db
}

// GetCollection returns a mongo.Collection for the given name
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// Close disconnects from MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Ping checks if the MongoDB server is reachable
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, readpref.Primary())
}

// EnsureIndexes creates indexes for the specified collection
func (m *MongoDB) EnsureIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error {
	_, err := m.GetCollection(collection).Indexes().CreateMany(ctx, indexes)
	return err
}
