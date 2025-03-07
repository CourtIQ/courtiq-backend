package repository

import (
	"context"
	"errors"

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines a generic interface for all repositories
type Repository[T any] interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*T, error)
	FindOne(ctx context.Context, filter interface{}) (*T, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*T, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
	Insert(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, id primitive.ObjectID, entity *T) (*T, error)
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

// Insert creates a new entity
func (r *BaseRepository[T]) Insert(ctx context.Context, entity *T) (*T, error) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to insert entity")
	}

	return entity, nil
}

// Update updates an existing entity
func (r *BaseRepository[T]) Update(ctx context.Context, id primitive.ObjectID, entity *T) (*T, error) {
	filter := bson.M{"_id": id}
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
