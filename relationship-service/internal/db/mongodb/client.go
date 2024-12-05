// internal/db/mongodb/client.go
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const RelationshipsCollection = "relationships"

type Client struct {
	client        *mongo.Client
	database      string
	relationships *mongo.Collection
}

type ClientOptions struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func NewClient(ctx context.Context, opts *ClientOptions) (*Client, error) {
	clientOpts := options.Client().
		ApplyURI(opts.URI).
		SetConnectTimeout(opts.Timeout)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify connection with a timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Initialize collection
	db := client.Database(opts.Database)

	return &Client{
		client:        client,
		database:      opts.Database,
		relationships: db.Collection(RelationshipsCollection),
	}, nil
}

// GetRelationships returns the relationships collection for database operations
func (c *Client) GetRelationships() *mongo.Collection {
	return c.relationships
}

// Close disconnects from MongoDB
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// CreateIndexes ensures required indexes exist for optimal query performance
func (c *Client) CreateIndexes(ctx context.Context) error {
	indexes := getRelationshipsIndexes()

	if _, err := c.relationships.Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}
	return nil
}
