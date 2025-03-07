package repository

import (
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryFactory provides methods for creating repositories
type RepositoryFactory struct {
	db *db.MongoDB
}

// NewRepositoryFactory creates a new RepositoryFactory
func NewRepositoryFactory(db *db.MongoDB) *RepositoryFactory {
	return &RepositoryFactory{db: db}
}

// NewRepository creates a new repository for the given entity type and collection name
// Correct approach: This is a standalone generic function
func NewRepository[T any](factory *RepositoryFactory, collectionName string) *BaseRepository[T] {
	collection := factory.db.GetCollection(collectionName)
	return NewBaseRepository[T](collection)
}

// Alternative approach: Use a method that takes a type as a parameter, not a generic method
func (f *RepositoryFactory) CreateRepository(collectionName string) interface{} {
	collection := f.db.GetCollection(collectionName)
	// You'd need to implement logic to determine which type of repository to create
	// based on other parameters or context
	return collection
}

// GetCollection returns a mongo.Collection for the given name
func (f *RepositoryFactory) GetCollection(name string) *mongo.Collection {
	return f.db.GetCollection(name)
}
