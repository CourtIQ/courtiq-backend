// cmd/main.go
package main

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/services"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	sharedRepo "github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/server"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create repository factory from shared package
	repoFactory := sharedRepo.NewRepositoryFactory(mongodb)

	// Initialize repositories using factory
	friendshipRepo := repository.NewFriendshipRepository(repoFactory)
	coachshipRepo := repository.NewCoachshipRepository(repoFactory)

	// Initialize services
	relationshipService := services.NewRelationshipService(friendshipRepo, coachshipRepo)

	// Initialize resolver with service dependency
	resolver := &resolvers.Resolver{
		RelationshipService: relationshipService,
	}

	// Initialize GraphQL server
	serverConfig := server.DefaultServerConfig()
	serverConfig.ServiceName = config.ServiceName
	serverConfig.Port = config.Port
	serverConfig.Environment = config.Environment
	serverConfig.PlaygroundEnabled = config.PlaygroundEnabled
	serverConfig.MongoDBURL = config.MongoDBURL
	serverConfig.DatabaseName = db.DatabaseName // Use constant from shared package

	// Create GraphQL server
	gqlServer, err := server.NewServer(serverConfig, graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	log.Printf("%s running in %s mode on port %d", config.ServiceName, config.Environment, config.Port)
	if err := gqlServer.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
