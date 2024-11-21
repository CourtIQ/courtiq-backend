package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config holds the service configuration values, including secrets.
type Config struct {
	ServiceName            string
	Port                   int
	Environment            string
	LogLevel               string
	GinMode                string
	PlaygroundEnabled      bool
	MongoDBURL             string
	FirebaseConfig         string
	FirebaseServiceAccount string
}

// LoadConfig loads configuration from environment variables or a .env file.
func LoadConfig() *Config {
	// Load environment variables from .env if present (for local development).
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	// Convert PORT from string to int with a fallback default.
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080 // Default internal port
	}

	// Check if GraphQL Playground is enabled
	playgroundEnabled := os.Getenv("GRAPHQL_PLAYGROUND") == "true"

	// Set default values for optional variables if not set
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "equipment-service"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}

	// Construct and return the Config struct
	return &Config{
		ServiceName:            serviceName,
		Port:                   port,
		Environment:            os.Getenv("GO_ENV"),
		LogLevel:               os.Getenv("LOG_LEVEL"),
		GinMode:                ginMode,
		PlaygroundEnabled:      playgroundEnabled,
		MongoDBURL:             os.Getenv("MONGODB_URL"),
		FirebaseConfig:         os.Getenv("FIREBASE_CONFIG"),
		FirebaseServiceAccount: os.Getenv("FIREBASE_SERVICE_ACCOUNT"),
	}
}

// SetupLogging configures logging based on environment and log level.
func SetupLogging(config *Config) {
	// Configure logging with timestamp and file information.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Adjust logging detail based on LogLevel
	switch config.LogLevel {
	case "debug":
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	case "info", "warn":
		log.SetFlags(log.LstdFlags)
	default:
		log.SetFlags(log.LstdFlags)
	}
}
