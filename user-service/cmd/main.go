package main

import (
	"fmt"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/user-service/configs"
	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/resolvers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Set Gin mode
	gin.SetMode(config.GinMode)

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{},
	}))

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	if config.LogLevel == "debug" {
		router.Use(gin.Logger())
	}

	// Setup routes
	if config.PlaygroundEnabled {
		log.Printf("GraphQL Playground enabled at /")
		router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
	}

	router.POST("/graphql", gin.WrapH(srv))

	// Get the service name for logging
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "service"
	}

	// Start server
	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("%s running in %s mode on %s", serviceName, config.Environment, address)
	if config.PlaygroundEnabled {
		log.Printf("GraphQL Playground available at http://localhost:%d", config.Port)
	}
	log.Printf("GraphQL endpoint available at http://localhost:%d/graphql", config.Port)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
