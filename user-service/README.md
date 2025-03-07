# User Service

This service handles user management and profiles in the CourtIQ platform.

## Architecture

The User Service follows a layered architecture:

- **GraphQL API Layer**: Handles incoming GraphQL requests
- **Service Layer**: Contains business logic and use cases
- **Repository Layer**: Uses the shared repository pattern for data access
- **Database Layer**: MongoDB for data persistence

## Key Components

### Services

The service layer implements the business logic:

- `UserServiceIntf`: Interface defining user operations
- `userService`: Implementation of the interface

### Repositories

This service uses the generic repository pattern from the shared package:

- `repository.Repository[model.User]`: Generic MongoDB repository for User model

### GraphQL

- `graph/schema.graphqls`: GraphQL schema definition
- `graph/resolvers`: GraphQL resolvers that use the service layer

## Getting Started

### Prerequisites

- Go 1.23 or later
- MongoDB instance
- Environment variables configured

### Running the Service

```bash
# Set environment variables or use .env file
export MONGODB_URL=mongodb://localhost:27017
export PORT=8080
export SERVICE_NAME=user-service
export GO_ENV=development
export GRAPHQL_PLAYGROUND=true

# Run the service
go run cmd/main.go
```

### Usage

Once the service is running, you can:

- Access the GraphQL playground at `http://localhost:8080/`
- Send GraphQL queries to `http://localhost:8080/graphql`

## API Operations

- `me`: Get the authenticated user's profile
- `getUser`: Get a user by ID
- `updateUser`: Update a user's profile
- `isUsernameAvailable`: Check if a username is available

## Design Decisions

### Shared Code

This service uses the shared package for:

- Repository pattern
- MongoDB connections
- Error handling
- Middleware (authentication)
- Configuration management

### Direct Repository Usage

Instead of implementing a service-specific repository, this service directly uses the shared repository pattern, which:

- Reduces code duplication
- Ensures consistent data access patterns
- Leverages type-safe generics
- Simplifies maintenance