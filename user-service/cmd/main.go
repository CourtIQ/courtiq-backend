package main

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/server"
	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/services"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Create server config from shared config
	serverConfig := server.DefaultServerConfig()
	serverConfig.LoadFromSharedConfig(*config)

	// Initialize MongoDB repository for users
	mongodb, err := db.NewMongoDB(context.Background(), serverConfig.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create user repository
	userRepo := repository.NewBaseRepository[model.User](mongodb.GetCollection(db.UsersCollection))

	// Create user service
	userService := services.NewUserService(userRepo)

	// Initialize GraphQL server with the resolver
	resolver := &resolvers.Resolver{
		UserService: userService,
	}

	// Create executable schema
	execSchema := graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	})

	// Create and configure the server
	srv, err := server.NewServer(serverConfig, execSchema)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Register directives
	srv.RegisterDirectives(server.GetDefaultDirectives())

	// Start the server
	if err := srv.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
