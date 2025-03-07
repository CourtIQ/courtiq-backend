# Shared Database Package

This package provides MongoDB client functionality and collection constants for all microservices.

## Features

- MongoDB client with connection pooling
- Collection name constants
- Health check support
- Index management

## Usage

### Connecting to MongoDB

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/db"

// Basic connection
mongodb, err := db.NewMongoDB(ctx, "mongodb://localhost:27017")
if err != nil {
    // Handle error
}

// Advanced connection with config
config := db.DefaultMongoDBConfig()
config.URI = os.Getenv("MONGODB_URL")
config.MaxPoolSize = 50
config.ConnectTimeout = 5 * time.Second

mongodb, err := db.NewMongoDBWithConfig(ctx, config)
if err != nil {
    // Handle error
}
```

### Accessing Collections

Use the predefined collection constants:

```go
// Get collection using constants
usersCollection := mongodb.GetCollection(db.UsersCollection)
friendshipsCollection := mongodb.GetCollection(db.FriendshipsCollection)
```

### Creating Indexes

```go
// Create indexes
indexes := []mongo.IndexModel{
    {
        Keys:    bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),
    },
    {
        Keys: bson.D{{Key: "createdAt", Value: 1}},
    },
}

err := mongodb.EnsureIndexes(ctx, db.UsersCollection, indexes)
if err != nil {
    // Handle index creation error
}
```

### Health Checks

```go
import "github.com/CourtIQ/courtiq-backend/shared/pkg/health"

// Create health handler
healthHandler := health.NewHandler(5 * time.Second)

// Add MongoDB health check
healthHandler.AddCheck("mongodb", func(ctx context.Context) (health.Status, error) {
    err := mongodb.Ping(ctx)
    if err != nil {
        return health.StatusDown, err
    }
    return health.StatusUp, nil
})
```

### Graceful Shutdown

```go
// In your main function
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := mongodb.Close(ctx); err != nil {
        log.Printf("Error closing MongoDB connection: %v", err)
    }
}()
```

## Collection Constants

Available collections:

```go
const (
    DatabaseName                   = "courtiq-db"
    UsersCollection                = "users"
    FriendshipsCollection          = "friendships"
    CoachshipsCollection           = "coachships"
    TennisRacketsCollection        = "tennis_rackets"
    TennisStringsCollection        = "tennis_strings"
    TennisCourtsCollection         = "tennis_courts"
    TennisMatchupsCollection       = "tennis_matchups"
    TennisMatchupsPointsCollection = "tennis_matchups_points"
)
```

## Adding Custom Collections

If your microservice needs additional collections, define them in your microservice's constants file:

```go
// In your microservice's package
package constants

import "github.com/CourtIQ/courtiq-backend/shared/pkg/db"

const (
    // Reference shared collections
    UsersCollection = db.UsersCollection
    
    // Add service-specific collections
    NotificationsCollection = "notifications"
    SettingsCollection      = "settings"
)
```