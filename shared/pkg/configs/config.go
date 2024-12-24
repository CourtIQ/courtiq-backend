// shared/pkg/configs/config.go
package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName            string
	Port                   int
	Environment            string
	LogLevel               string
	PlaygroundEnabled      bool
	MongoDBURL             string
	FirebaseConfig         string
	FirebaseServiceAccount string
}

func LoadConfig() *Config {
	// Load environment variables from .env if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	// Convert PORT from string to int with fallback
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080 // Default port
	}

	return &Config{
		ServiceName:            getEnvWithDefault("SERVICE_NAME", "relationship-service"),
		Port:                   port,
		Environment:            getEnvWithDefault("GO_ENV", "development"),
		LogLevel:               getEnvWithDefault("LOG_LEVEL", "info"),
		PlaygroundEnabled:      os.Getenv("GRAPHQL_PLAYGROUND") == "true",
		MongoDBURL:             os.Getenv("MONGODB_URL"),
		FirebaseConfig:         os.Getenv("FIREBASE_CONFIG"),
		FirebaseServiceAccount: os.Getenv("FIREBASE_SERVICE_ACCOUNT"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func SetupLogging(config *Config) {
	// Configure logging with timestamp and file information
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
