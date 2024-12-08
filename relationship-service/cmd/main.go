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
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/directives"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/configs"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/db"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/middleware"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/services"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Setup logging
	configs.SetupLogging(config)

	// Initialize MongoDB client
	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Get the relationships collection
	coll := mongodb.GetCollection(db.RelationshipsCollection)

	// Create the repository using the collection
	relationshipRepo := repository.NewRelationshipRepository(coll)

	// Create the service with the repository
	relationshipService := services.NewRelationshipService(relationshipRepo)

	// Set the directive dependencies
	directives.RelationshipRepo = relationshipRepo
	directives.GetCurrentUserID = middleware.GetUserIDFromContext

	// Build gqlgen config and assign the directive
	c := graph.Config{
		Resolvers: &resolvers.Resolver{
			RelationshipService: relationshipService,
		},
	}

	c.Directives.Satisfies = directives.SatisfiesDirective
	// Create the executable schema
	schema := graph.NewExecutableSchema(c)

	// Set up the GraphQL server
	srv := handler.NewDefaultServer(schema)

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
