package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/resolvers"
	"github.com/CourtIQ/courtiq-backend/user-service/services"
)

func main() {
	// Retrieve port from environment variables, defaulting to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize MongoDB connection
	services.ConnectToMongoDB()
	defer services.DisconnectFromMongoDB()

	// Initialize your UserService implementation
	userService := services.NewUserService() // Use the proper constructor

	// Create a new Resolver and inject the service
	resolver := &resolvers.Resolver{
		UserService: userService,
	}

	// Create the GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Set up HTTP routes for the GraphQL playground and query handler
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Start the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
