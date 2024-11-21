package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(ctx context.Context, uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (c *Client) Database() *mongo.Database {
	return c.db
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
