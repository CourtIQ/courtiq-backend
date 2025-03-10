package main

import (
	"context"
	"log"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/matchup-service/internal/services"
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
	matchUpRepo := repository.NewMatchUpRepository(mongodb)
	pointsRepo := repository.NewPointsRepositoru(mongodb)

	// Create matchup service
	matchUpService := services.NewMatchUpService(matchUpRepo, pointsRepo)

	// Initialize resolver
	resolver := &resolvers.Resolver{
		MatchUpServiceInterface: matchUpService,
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
