package main

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/services"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/server"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Create server config from shared config
	serverConfig := server.DefaultServerConfig()
	serverConfig.LoadFromSharedConfig(*config)

	// Initialize MongoDB client
	mongodb, err := db.NewMongoDB(context.Background(), serverConfig.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create repositories
	racketRepo := repository.NewTennisRacketMongoRepo(mongodb)
	stringRepo := repository.NewTennisStringMongoRepo(mongodb)

	// Create equipment service
	equipmentService := services.NewEquipmentService(racketRepo, stringRepo)

	// Initialize resolver
	resolver := &resolvers.Resolver{
		EquipmentServiceIntf: equipmentService,
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
