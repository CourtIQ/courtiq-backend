// cmd/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/configs"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/db"
	// If you have a `utils` package with middleware:
	// "github.com/CourtIQ/courtiq-backend/equipment-service/internal/utils"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Initialize MongoDB client
	_, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create repositories
	// racketRepo := repository.NewCoachshipRepository(mongodb)
	// stringRepo := repository.NewFriendshipRepository(mongodb)

	// Create the service with the repositories
	// equipmentService := services.NewRelationshipService(racketRepo, stringRepo)

	// Set up the GraphQL server with the resolver that has the service injected
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			// EquipmentServiceIntf: equipmentService,
		},
	}))

	// Create router mux
	mux := http.NewServeMux()

	// If you'd like to add user middleware, for example:
	// mux.Handle("/graphql", utils.NewUserMiddleware(utils.AuthConfig{EnableAuth: true})(srv))
	// Otherwise, just use srv directly.
	// For now, let's just mount srv directly:
	mux.Handle("/graphql", srv)

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
