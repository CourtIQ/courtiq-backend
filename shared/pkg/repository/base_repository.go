package repository

import (
	"context"
	"errors"
	"reflect" // Add reflection package

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines a generic interface for all repositories
type Repository[T any] interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*T, error)
	FindByIDWithFilters(ctx context.Context, id primitive.ObjectID, additionalFilters bson.M) (*T, error)
	FindOne(ctx context.Context, filter interface{}) (*T, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*T, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
	// Change Insert signature to return the inserted entity (interface{}) and error
	Insert(ctx context.Context, entity interface{}) (interface{}, error)
	// Change Update signature to accept interface{}
	Update(ctx context.Context, id primitive.ObjectID, entity interface{}) (*T, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*T, error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (*T, error)
}

// BaseRepository provides a basic implementation of Repository
type BaseRepository[T any] struct {
	collection *mongo.Collection
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository[T any](collection *mongo.Collection) *BaseRepository[T] {
	return &BaseRepository[T]{
		collection: collection,
	}
}

// FindByID finds an entity by its ObjectID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
	filter := bson.M{"_id": id}

	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "failed to find entity by ID")
	}

	return &entity, nil
}

// FindOne finds a single entity matching the given filter
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter interface{}) (*T, error) {
	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No error, just nil entity
		}
		return nil, sharedErrors.WrapError(err, "failed to find entity")
	}

	return &entity, nil
}

// Find finds all entities matching the given filter
func (r *BaseRepository[T]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*T, error) {
	findOptions := options.Find()
	if len(opts) > 0 {
		findOptions = opts[0]
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to find entities")
	}
	defer cursor.Close(ctx)

	var entities []*T
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, sharedErrors.WrapError(err, "failed to decode entities")
	}

	return entities, nil
}

// Count counts documents matching the given filter
func (r *BaseRepository[T]) Count(ctx context.Context, filter interface{}) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, sharedErrors.WrapError(err, "failed to count entities")
	}

	return count, nil
}

// Insert creates a new entity and attempts to set its ID field using reflection.
// It returns the original entity interface (potentially modified with the ID) and an error.
func (r *BaseRepository[T]) Insert(ctx context.Context, entity interface{}) (interface{}, error) {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		// Check for duplicate key errors specifically if needed
		// if mongo.IsDuplicateKeyError(err) { ... }
		return nil, sharedErrors.WrapError(err, "failed to insert entity")
	}

	// Assert the inserted ID to primitive.ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		// This should not happen with standard MongoDB ObjectID usage
		// Return the entity without ID set in this edge case, but log the error.
		// log.Printf("Warning: InsertedID was not a primitive.ObjectID: %T", result.InsertedID)
		return entity, sharedErrors.WrapError(errors.New("invalid type for inserted ID"), "failed to get inserted ID type")

	}

	// Use reflection to set the ID field on the entity
	// This assumes the entity is a pointer to a struct with an exported field named "ID" of type primitive.ObjectID
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityElem := entityValue.Elem()
		if entityElem.Kind() == reflect.Struct {
			idField := entityElem.FieldByName("ID")
			// Check if field exists, is the correct type, and is settable
			if idField.IsValid() && idField.Type() == reflect.TypeOf(primitive.ObjectID{}) && idField.CanSet() {
				idField.Set(reflect.ValueOf(insertedID))
			} else {
				// Optional: Log a warning if ID field couldn't be set
				// log.Printf("Warning: Could not set ID field on inserted entity type %T", entity)
			}
		}
	} else {
		// Optional: Log a warning if entity is not a pointer
		// log.Printf("Warning: Inserted entity is not a pointer type %T, ID field cannot be set", entity)
	}

	// Return the original entity interface, potentially updated with the ID
	return entity, nil
}

// Update updates an existing entity
// Modify entity parameter to interface{} for flexibility
func (r *BaseRepository[T]) Update(ctx context.Context, id primitive.ObjectID, entity interface{}) (*T, error) {
	filter := bson.M{"_id": id}
	// Use bson.M{"$set": entity} which works well if entity is a struct or map
	update := bson.M{"$set": entity}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedEntity T
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "failed to update entity")
	}

	return &updatedEntity, nil
}

// FindByIDWithFilters finds an entity by its ObjectID with additional filters
func (r *BaseRepository[T]) FindByIDWithFilters(ctx context.Context, id primitive.ObjectID, additionalFilters bson.M) (*T, error) {
	// Start with the ID filter
	filter := bson.M{"_id": id}

	// Merge with additional filters
	for k, v := range additionalFilters {
		filter[k] = v
	}

	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "failed to find entity by ID with filters")
	}

	return &entity, nil
}

// Delete deletes an entity by its ID
func (r *BaseRepository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return sharedErrors.WrapError(err, "failed to delete entity")
	}

	if result.DeletedCount == 0 {
		return sharedErrors.ErrNotFound
	}

	return nil
}

// FindOneAndUpdate finds an entity matching the given filter, applies the update, and returns the updated document
func (r *BaseRepository[T]) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	findOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if len(opts) > 0 {
		findOpts = opts[0]
	}

	var updatedEntity T
	err := r.collection.FindOneAndUpdate(ctx, filter, update, findOpts).Decode(&updatedEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "failed to find and update entity")
	}

	return &updatedEntity, nil
}

// FindOneAndDelete finds an entity matching the given filter, deletes it, and returns the deleted document
func (r *BaseRepository[T]) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (*T, error) {
	findOpts := options.FindOneAndDelete()
	if len(opts) > 0 {
		findOpts = opts[0]
	}

	var deletedEntity T
	err := r.collection.FindOneAndDelete(ctx, filter, findOpts).Decode(&deletedEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.WrapError(err, "failed to find and delete entity")
	}

	return &deletedEntity, nil
}