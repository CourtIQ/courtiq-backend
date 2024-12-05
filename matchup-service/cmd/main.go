package main

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/matchup-service/configs"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/resolvers"

	// "github.com/CourtIQ/courtiq-backend/matchup-service/internal/db"

	// "github.com/CourtIQ/courtiq-backend/matchup-service/internal/repository"
	// "github.com/CourtIQ/courtiq-backend/matchup-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Set Gin mode
	gin.SetMode(config.GinMode)

	// Initialize MongoDB connection using config.MongoDBURL
	// Use "matchupdb" as the database name
	// mongoDB := db.NewMongoDB(config.MongoDBURL, "courtiq-db")

	// Initialize repositories and services
	// matchupRepository := repository.NewMatchupRepository(mongoDB)
	// matchupService := service.NewMatchupService(matchupRepository)

	// Initialize GraphQL server with the resolved schemas and dependencies
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			// MatchupService: matchupService,
		},
	}))

	// Create a Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	if config.LogLevel == "debug" {
		router.Use(gin.Logger())
	}

	// Setup routes
	if config.PlaygroundEnabled {
		log.Printf("GraphQL Playground enabled at /")
		router.GET("/", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))
	}

	router.POST("/graphql", gin.WrapH(srv))

	// Start the server
	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("%s running in %s mode on %s", config.ServiceName, config.Environment, address)
	if config.PlaygroundEnabled {
		log.Printf("GraphQL Playground available at http://localhost:%d", config.Port)
	}
	log.Printf("GraphQL endpoint available at http://localhost:%d/graphql", config.Port)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
