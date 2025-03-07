package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"SERVICE_NAME":    os.Getenv("SERVICE_NAME"),
		"PORT":            os.Getenv("PORT"),
		"GO_ENV":          os.Getenv("GO_ENV"),
		"LOG_LEVEL":       os.Getenv("LOG_LEVEL"),
		"MONGODB_URL":     os.Getenv("MONGODB_URL"),
		"FIREBASE_CONFIG": os.Getenv("FIREBASE_CONFIG"),
	}

	// Restore original environment variables after test
	defer func() {
		for k, v := range originalEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Test 1: Default values when environment variables are not set
	os.Unsetenv("SERVICE_NAME")
	os.Unsetenv("PORT")
	os.Unsetenv("GO_ENV")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("GRAPHQL_PLAYGROUND")
	os.Unsetenv("MONGODB_URL")
	os.Unsetenv("FIREBASE_CONFIG")
	
	config := LoadConfig()
	
	assert.Equal(t, "relationship-service", config.ServiceName)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "info", config.LogLevel)
	assert.False(t, config.PlaygroundEnabled)
	assert.Equal(t, "", config.MongoDBURL)

	// Test 2: Custom values from environment variables
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("PORT", "9000")
	os.Setenv("GO_ENV", "test")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("GRAPHQL_PLAYGROUND", "true")
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	os.Setenv("FIREBASE_CONFIG", "{\"testKey\": \"testValue\"}")
	
	config = LoadConfig()
	
	assert.Equal(t, "test-service", config.ServiceName)
	assert.Equal(t, 9000, config.Port)
	assert.Equal(t, "test", config.Environment)
	assert.Equal(t, "debug", config.LogLevel)
	assert.True(t, config.PlaygroundEnabled)
	assert.Equal(t, "mongodb://localhost:27017", config.MongoDBURL)
	assert.Equal(t, "{\"testKey\": \"testValue\"}", config.FirebaseConfig)
}

func TestSetupLogging(t *testing.T) {
	// This test just ensures the function doesn't panic
	config := &Config{
		LogLevel: "info",
	}
	
	// Should not panic
	SetupLogging(config)
	
	config.LogLevel = "debug"
	SetupLogging(config)
	
	config.LogLevel = "warn"
	SetupLogging(config)
}