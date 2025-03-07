# Shared Package for CourtIQ Microservices

This package contains shared functionality used across all CourtIQ microservices.

## Package Structure

- **configs**: Configuration loading and management
- **db**: MongoDB client and collection constants
- **errors**: Standardized error types and handling
- **health**: Health check implementations
- **middleware**: Common HTTP and GraphQL middleware
- **repository**: Generic repository pattern implementation

## Usage Guidelines

### 1. Error Handling

Use the shared error types from `errors` package for consistent error handling across services:

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"

// Return standard errors
return errors.ErrNotFound
return errors.ErrInvalidInput

// Wrap errors with context
return errors.WrapError(err, "failed to process request")

// Create validation errors
return errors.NewValidationError("email", "invalid format")

// Check error types
if errors.IsNotFoundError(err) {
    // Handle not found
}
```

### 2. MongoDB Access

Use the MongoDB client from the `db` package:

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/db"

// Create a new MongoDB client
mongodb, err := db.NewMongoDB(ctx, mongoURI)
if err != nil {
    // Handle error
}

// Access collections using constants
usersCollection := mongodb.GetCollection(db.UsersCollection)
```

### 3. Repository Pattern

Use the generic repository from the `repository` package:

```go
import (
    "github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
    "github.com/CourtIQ/courtiq-backend/shared/pkg/db"
)

// Create a repository factory
mongodb, _ := db.NewMongoDB(ctx, mongoURI)
factory := repository.NewRepositoryFactory(mongodb)

// Create a repository for a specific model
userRepo := factory.NewRepository[UserModel](db.UsersCollection)

// Use repository methods
user, err := userRepo.FindByID(ctx, id)
users, err := userRepo.Find(ctx, filter)
```

### 4. Health Checks

Use the health check handler from the `health` package:

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/health"

// Create a health handler
healthHandler := health.NewHandler(5 * time.Second)

// Add checks
healthHandler.AddMongoDBCheck("mongodb", mongoClient)

// Register HTTP handlers
http.HandleFunc("/health/live", healthHandler.HandleLiveness)
http.HandleFunc("/health/ready", healthHandler.HandleReadiness)
```

### 5. Middleware

Use middleware from the `middleware` package:

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"

// Add user claims middleware
http.Handle("/api", middleware.WithUserClaims(apiHandler))

// Get user from context
userID, err := middleware.GetMongoIDFromContext(ctx)
```

## Extending Shared Functionality

When adding new functionality to the shared package:

1. Keep shared code generic enough to be used across multiple services
2. Add thorough documentation and examples
3. Write tests for all shared code
4. Create domain-specific extensions in your microservice