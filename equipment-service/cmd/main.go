package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/resolvers"

	// "github.com/CourtIQ/courtiq-backend/equipment-service/graph/resolvers"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port              int
	Environment       string
	LogLevel          string
	GinMode           string
	PlaygroundEnabled bool
}

func loadConfig() *Config {
	// Load port with fallback
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080 // Default internal port from config
	}

	// Parse playground setting
	playgroundEnabled := false
	if playground := os.Getenv("GRAPHQL_PLAYGROUND"); playground == "true" {
		playgroundEnabled = true
	}

	return &Config{
		Port:              port,
		Environment:       os.Getenv("GO_ENV"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		GinMode:           os.Getenv("GIN_MODE"),
		PlaygroundEnabled: playgroundEnabled,
	}
}

func setupLogging(config *Config) {
	// Configure logging based on environment
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// You could add more sophisticated logging setup here
	switch config.LogLevel {
	case "debug":
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	case "info":
		log.SetFlags(log.LstdFlags)
	case "warn":
		log.SetFlags(log.LstdFlags)
	default:
		log.SetFlags(log.LstdFlags)
	}
}

func main() {
	// Load configuration from environment
	config := loadConfig()

	// Setup logging
	setupLogging(config)

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
