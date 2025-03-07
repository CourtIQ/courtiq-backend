package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/user-service/internal/services"
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

	// Create repository using the shared package's repository
	userRepo := repository.NewBaseRepository[model.User](mongodb.GetCollection(db.UsersCollection))

	// Create the service with the repository
	userService := services.NewUserService(userRepo)

	// Set up the GraphQL server with the resolver that has the service injected
	// Using non-deprecated handler.New instead of NewDefaultServer
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			UserServiceIntf: userService,
		},
	}))

	// Configure the server with defaults
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.Introspection{})

	// Create router mux
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
