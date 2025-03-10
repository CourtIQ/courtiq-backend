package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/access"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server represents a GraphQL server instance
type Server struct {
	Config         ServerConfig
	GraphQLHandler *handler.Server
	Router         *chi.Mux
	MongoDB        *db.MongoDB
	AccessChecker  access.Checker
}

// NewServer creates a new server with the given configuration and GraphQL schema
func NewServer(config ServerConfig, es graphql.ExecutableSchema) (*Server, error) {
	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(context.Background(), config.MongoDBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Create the GraphQL server
	srv := handler.NewDefaultServer(es)

	// Add default transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Add default extensions
	srv.Use(extension.Introspection{})

	// Create router
	router := chi.NewRouter()

	// Create server instance
	server := &Server{
		Config:         config,
		GraphQLHandler: srv,
		Router:         router,
		MongoDB:        mongodb,
	}

	// Setup middleware and routes
	server.setupMiddleware()
	server.setupRoutes()

	// Setup access control if enabled
	if config.EnableAccessControl {
		if err := server.setupAccessControl(); err != nil {
			return nil, fmt.Errorf("failed to setup access control: %w", err)
		}
	}

	return server, nil
}

// setupMiddleware adds middleware to the router
func (s *Server) setupMiddleware() {
	// Add user claims middleware
	s.Router.Use(middleware.WithUserClaims)

	// Add metrics middleware if enabled
	if s.Config.EnableMetrics {
		s.GraphQLHandler.Use(middleware.GetMetricsConfig(s.Config.ServiceName))
	}
}

// setupRoutes configures server routes
func (s *Server) setupRoutes() {
	// Add GraphQL endpoint
	s.Router.Handle("/graphql", s.GraphQLHandler)

	// Add playground if enabled
	if s.Config.PlaygroundEnabled {
		s.Router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	}

	// Add health check endpoint
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add metrics endpoint if enabled
	if s.Config.EnableMetrics {
		s.Router.Handle("/metrics", promhttp.Handler())
	}
}

// setupAccessControl initializes the access control checker
func (s *Server) setupAccessControl() error {
	// Setup relationship checker
	config := access.DefaultRelationshipCheckerConfig()
	config.Client = s.MongoDB.GetClient()
	config.Database = s.Config.DatabaseName

	// Create access checker and middleware
	s.AccessChecker = access.NewRelationshipChecker(config)

	// Add middleware to the GraphQL server - using AroundOperations instead of Use
	s.GraphQLHandler.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		ctx = context.WithValue(ctx, access.CheckerContextKey, s.AccessChecker)
		return next(ctx)
	})

	return nil
}

// Serve starts the HTTP server with graceful shutdown
func (s *Server) Serve() error {
	// Create HTTP server
	addr := fmt.Sprintf(":%d", s.Config.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	// Channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Channel to track server errors
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("%s running in %s mode on %s", s.Config.ServiceName, s.Config.Environment, addr)
		if s.Config.PlaygroundEnabled {
			log.Printf("GraphQL Playground available at http://localhost:%d/", s.Config.Port)
		}
		log.Printf("GraphQL endpoint available at http://localhost:%d/graphql", s.Config.Port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Block until we receive a shutdown signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-shutdown:
		log.Println("Shutting down server...")

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt to gracefully shutdown
		if err := httpServer.Shutdown(ctx); err != nil {
			// If graceful shutdown fails, force it
			if err := httpServer.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}

		// Close MongoDB connection
		if err := s.MongoDB.Close(ctx); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}

	return nil
}

// RegisterDirectives registers GraphQL directives with the server
func (s *Server) RegisterDirectives(directives DirectivesMap) {
	for name, fn := range directives {
		s.GraphQLHandler.AroundFields(CreateFieldDirectiveMiddleware(name, fn))
	}
}

// Close cleans up server resources
func (s *Server) Close(ctx context.Context) error {
	if s.MongoDB != nil {
		return s.MongoDB.Close(ctx)
	}
	return nil
}