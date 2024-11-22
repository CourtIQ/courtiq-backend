package main

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/user-service/configs"
	"github.com/CourtIQ/courtiq-backend/user-service/db"
	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/user-service/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from environment
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Set Gin mode
	gin.SetMode(config.GinMode)

	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(mongodb)

	// Initialize Resolver
	resolver := &resolvers.Resolver{
		UserService: userService,
	}

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
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

	// Start server
	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("Server running on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
