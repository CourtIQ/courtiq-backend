package services

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectToMongoDB initializes the MongoDB connection
func ConnectToMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in environment variables")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Ping to ensure a successful connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	log.Println("Successfully connected to MongoDB")
}

// DisconnectFromMongoDB closes the MongoDB connection
func DisconnectFromMongoDB() {
	if err := MongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("Disconnected from MongoDB")
}

// GetMongoCollection returns a MongoDB collection by database and collection name
func GetMongoCollection(dbName string, collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoClient is not initialized. Please connect to MongoDB first.")
	}
	return MongoClient.Database(dbName).Collection(collectionName)
}
