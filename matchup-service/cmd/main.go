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
	sharedRepo "github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/server"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	configs.SetupLogging(config)

	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	repoFactory := sharedRepo.NewRepositoryFactory(mongodb)

	// Create repositories
	matchUpRepo := repository.NewMatchupsRepository(repoFactory)
	pointsRepo := repository.NewShotsRepository(repoFactory)

	// Create matchup service
	matchUpService := services.NewMatchUpService(matchUpRepo, pointsRepo)

	// Initialize resolver
	resolver := &resolvers.Resolver{
		MatchUpServiceInterface: matchUpService,
	}

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
