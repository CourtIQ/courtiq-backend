package server

import (
	"os"
	"strconv"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
)

// ServerConfig holds configuration for GraphQL server
type ServerConfig struct {
	// Basic server settings
	ServiceName      string
	Port             int
	Environment      string
	PlaygroundEnabled bool

	// Database settings
	MongoDBURL       string
	DatabaseName     string

	// Optional settings
	EnableMetrics     bool
	EnableAccessControl bool
	RelationshipDBURL  string // For access control
}

// DefaultServerConfig returns a default configuration with values from environment
func DefaultServerConfig() ServerConfig {
	// Get port from env or use default
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Get service name or use default
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "graphql-service"
	}

	// Get environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Default MongoDB URL
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}

	// Default database name
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "courtiq-db"
	}

	// Check if playground is enabled
	playgroundEnabled := true
	if playgroundStr := os.Getenv("PLAYGROUND_ENABLED"); playgroundStr == "false" {
		playgroundEnabled = false
	}

	// Check if metrics are enabled
	enableMetrics := true
	if metricsStr := os.Getenv("ENABLE_METRICS"); metricsStr == "false" {
		enableMetrics = false
	}

	// Check if access control is enabled
	enableAccessControl := true
	if acStr := os.Getenv("ENABLE_ACCESS_CONTROL"); acStr == "false" {
		enableAccessControl = false
	}

	return ServerConfig{
		ServiceName:        serviceName,
		Port:               port,
		Environment:        env,
		PlaygroundEnabled:  playgroundEnabled,
		MongoDBURL:         mongoURL,
		DatabaseName:       dbName,
		EnableMetrics:      enableMetrics,
		EnableAccessControl: enableAccessControl,
		RelationshipDBURL:   mongoURL, // Default same as main DB
	}
}

// LoadFromSharedConfig loads configuration from shared configs
func (c *ServerConfig) LoadFromSharedConfig(config configs.Config) {
	c.ServiceName = config.ServiceName
	c.Port = config.Port
	c.Environment = config.Environment
	c.PlaygroundEnabled = config.PlaygroundEnabled
	c.MongoDBURL = config.MongoDBURL
	// Other fields keep their default values unless explicitly set
}