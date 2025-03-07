# Shared Repository Package

This package provides a generic repository pattern implementation for MongoDB.

## Features

- Generic CRUD operations that work with any model type
- Consistent error handling
- Pagination support
- Flexible query capabilities

## Usage

### Creating a Repository

```go
import (
    "github.com/CourtIQ/courtiq-backend/shared/pkg/db"
    "github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
)

// Create a repository factory
mongodb, _ := db.NewMongoDB(ctx, mongoURI)
factory := repository.NewRepositoryFactory(mongodb)

// Create a repository for a specific model
userRepo := factory.NewRepository[UserModel](db.UsersCollection)
```

### Repository Methods

```go
// Find by ID
user, err := userRepo.FindByID(ctx, id)

// Insert
newUser, err := userRepo.Insert(ctx, user)

// Update
updatedUser, err := userRepo.Update(ctx, id, user)

// Delete
err := userRepo.Delete(ctx, id)

// Find with filter
users, err := userRepo.Find(ctx, filter, findOptions)

// Find one with filter
user, err := userRepo.FindOne(ctx, filter)

// Find one and update
updatedUser, err := userRepo.FindOneAndUpdate(ctx, filter, update)

// Find one and delete
deletedUser, err := userRepo.FindOneAndDelete(ctx, filter)

// Count
count, err := userRepo.Count(ctx, filter)
```

### Working with service-specific repositories

For complex queries or custom operations, create a service-specific repository that wraps the BaseRepository:

```go
// Define your repository interface
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindActiveUsers(ctx context.Context) ([]*User, error)
    // Include all base repository methods
    repository.Repository[User]
}

// Implement the repository
type userRepository struct {
    *repository.BaseRepository[User]
}

func NewUserRepository(mdb *db.MongoDB) UserRepository {
    baseRepo := repository.NewBaseRepository[User](mdb.GetCollection(db.UsersCollection))
    return &userRepository{baseRepo}
}

// Implement custom methods
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    return r.FindOne(ctx, bson.M{"email": email})
}

func (r *userRepository) FindActiveUsers(ctx context.Context) ([]*User, error) {
    return r.Find(ctx, bson.M{"status": "active"})
}
```

## Error Handling

All repository methods return domain-specific errors from the shared errors package:

- `errors.ErrNotFound`: When a requested entity doesn't exist
- `errors.ErrDatabaseOperation`: For general database errors
- `errors.ErrInvalidInput`: For invalid inputs or parameters

Example error handling:

```go
user, err := userRepo.FindByID(ctx, id)
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle not found
    } else {
        // Handle other errors
    }
}
```

## Indexes and Performance

For microservice-specific indexes and performance optimizations:

```go
// In your microservice initialization
func initIndexes(ctx context.Context, mdb *db.MongoDB) error {
    // Create indexes for your specific collections
    indexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{{Key: "createdAt", Value: 1}},
        },
    }
    
    return mdb.EnsureIndexes(ctx, db.UsersCollection, indexes)
}
```