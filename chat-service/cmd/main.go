package main

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph"
	"github.com/CourtIQ/courtiq-backend/chat-service/graph/resolvers"
	chatRepo "github.com/CourtIQ/courtiq-backend/chat-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/chat-service/internal/services"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/server"
)

func main() {
	// Create server config from shared config
	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Initialize MongoDB client
	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create Repository Factory
	repoFactory := repository.NewRepositoryFactory(mongodb)

	// Instantiate repositories
	messageRepo := chatRepo.NewMessageRepository(repoFactory)
	chatRepo := chatRepo.NewChatRepository(repoFactory) // Assuming this exists or will be created

	chatService := services.NewChatService(chatRepo, messageRepo)

	resolver := &resolvers.Resolver{
		ChatService: chatService,
	}

	// Initialize GraphQL server
	serverConfig := server.DefaultServerConfig()
	serverConfig.ServiceName = config.ServiceName
	serverConfig.Port = config.Port
	serverConfig.Environment = config.Environment
	serverConfig.PlaygroundEnabled = config.PlaygroundEnabled
	serverConfig.MongoDBURL = config.MongoDBURL
	serverConfig.DatabaseName = db.DatabaseName

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
