// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/configs"
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

	// Get the relationships collection
	// coll := mongodb.GetCollection(db.EquipmentCollection)

	// Create the repository using the collection
	// equipmentRepo := repository.NewEquipmentService(coll)

	// Create the service with the repository
	// equipmentService := services.NewRelationshipService(equipmentRepo)

	// Set up the GraphQL server with the resolver that has the service injected
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			// EquipmentService: equipmentService,
		},
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
