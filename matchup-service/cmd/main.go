package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Create the service (using dummy implementation for now)
	// matchupService := service.NewMatchUpService()

	// Set up the GraphQL server with the resolver that has the service injected
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		// Resolvers: resolvers.NewRcesolver(matchupService),
	}))

	// Create router mux
	mux := http.NewServeMux()

	// Setup routes
	if config.PlaygroundEnabled {
		mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		log.Printf("GraphQL Playground enabled at /")
	}

	mux.Handle("/graphql", srv)

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
