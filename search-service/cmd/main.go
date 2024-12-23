package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/search-service/graph"
	"github.com/CourtIQ/courtiq-backend/search-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/configs"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/middleware"
)

func main() {

	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Initialize MongoDB client
	// mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB: %v", err)
	// }

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{},
	}))

	// Configure transports with specific order and options
	srv.AddTransport(transport.POST{})    // Standard GraphQL POST requests
	srv.AddTransport(transport.Options{}) // CORS preflight requests
	srv.AddTransport(transport.GET{})     // Simple queries via GET

	// Add basic extensions
	srv.Use(extension.Introspection{}) // Enable schema introspection

	// Rest of your code remains the same...
	mux := http.NewServeMux()
	mux.Handle("/graphql", middleware.WithUserClaims(srv))

	// Setup playground if enabled
	if config.PlaygroundEnabled {
		mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		log.Printf("GraphQL Playground enabled at /")
	}

	// Start server
	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("%s running in %s mode on %s", config.ServiceName, config.Environment, address)
	if config.PlaygroundEnabled {
		log.Printf("GraphQL Playground available at http://localhost:%d", config.Port)
	}
	log.Printf("GraphQL endpoint available at http://localhost:%d/graphql", config.Port)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
