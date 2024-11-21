package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB encapsulates the Mongo client and database references
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB client and returns a MongoDB struct containing the client and the database reference
func NewMongoDB(uri, dbName string) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection to ensure it's successful
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	database := client.Database(dbName)
	log.Printf("Connected to MongoDB at %s, using database '%s'", uri, dbName)

	return &MongoDB{
		Client:   client,
		Database: database,
	}
}
