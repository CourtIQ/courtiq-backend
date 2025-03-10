package repository

import (
	"context"
	"testing"

	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test entity
type TestEntity struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string            `bson:"name"`
	Age  int               `bson:"age"`
}

// TestRepository is a mock implementation for testing
type TestRepository struct {
	*BaseRepository[TestEntity]
	findByIDFunc             func(ctx context.Context, id primitive.ObjectID) (*TestEntity, error)
	findByIDWithFiltersFunc  func(ctx context.Context, id primitive.ObjectID, additionalFilters bson.M) (*TestEntity, error)
	findOneFunc              func(ctx context.Context, filter interface{}) (*TestEntity, error)
	insertFunc               func(ctx context.Context, entity *TestEntity) (*TestEntity, error)
	updateFunc               func(ctx context.Context, id primitive.ObjectID, entity *TestEntity) (*TestEntity, error)
	deleteFunc               func(ctx context.Context, id primitive.ObjectID) error
	findFunc                 func(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*TestEntity, error)
	countFunc                func(ctx context.Context, filter interface{}) (int64, error)
	findOneAndUpdateFunc     func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*TestEntity, error)
	findOneAndDeleteFunc     func(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (*TestEntity, error)
}

// Override methods to use mock functions
func (r *TestRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*TestEntity, error) {
	if r.findByIDFunc != nil {
		return r.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (r *TestRepository) FindByIDWithFilters(ctx context.Context, id primitive.ObjectID, additionalFilters bson.M) (*TestEntity, error) {
	if r.findByIDWithFiltersFunc != nil {
		return r.findByIDWithFiltersFunc(ctx, id, additionalFilters)
	}
	return nil, nil
}

func (r *TestRepository) FindOne(ctx context.Context, filter interface{}) (*TestEntity, error) {
	if r.findOneFunc != nil {
		return r.findOneFunc(ctx, filter)
	}
	return nil, nil
}

func (r *TestRepository) Insert(ctx context.Context, entity *TestEntity) (*TestEntity, error) {
	if r.insertFunc != nil {
		return r.insertFunc(ctx, entity)
	}
	return entity, nil
}

func (r *TestRepository) Update(ctx context.Context, id primitive.ObjectID, entity *TestEntity) (*TestEntity, error) {
	if r.updateFunc != nil {
		return r.updateFunc(ctx, id, entity)
	}
	return entity, nil
}

func (r *TestRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if r.deleteFunc != nil {
		return r.deleteFunc(ctx, id)
	}
	return nil
}

func (r *TestRepository) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*TestEntity, error) {
	if r.findFunc != nil {
		return r.findFunc(ctx, filter, opts...)
	}
	return nil, nil
}

func (r *TestRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	if r.countFunc != nil {
		return r.countFunc(ctx, filter)
	}
	return 0, nil
}

func (r *TestRepository) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*TestEntity, error) {
	if r.findOneAndUpdateFunc != nil {
		return r.findOneAndUpdateFunc(ctx, filter, update, opts...)
	}
	return nil, nil
}

func (r *TestRepository) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (*TestEntity, error) {
	if r.findOneAndDeleteFunc != nil {
		return r.findOneAndDeleteFunc(ctx, filter, opts...)
	}
	return nil, nil
}

// NewTestRepository creates a test repository with mocked methods
func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func TestRepositoryFindByID(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test entity
	id := primitive.NewObjectID()
	entity := &TestEntity{
		ID:   id,
		Name: "Test User",
		Age:  30,
	}

	// Set up mock function
	repo.findByIDFunc = func(ctx context.Context, testID primitive.ObjectID) (*TestEntity, error) {
		if testID == id {
			return entity, nil
		}
		return nil, sharedErrors.ErrNotFound
	}

	// Test success case
	result, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Name, result.Name)
	assert.Equal(t, entity.Age, result.Age)

	// Test not found case
	result, err = repo.FindByID(context.Background(), primitive.NewObjectID())
	assert.Error(t, err)
	assert.True(t, sharedErrors.IsNotFoundError(err))
	assert.Nil(t, result)
}

func TestRepositoryFindOne(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test entity
	entity := &TestEntity{
		ID:   primitive.NewObjectID(),
		Name: "Test User",
		Age:  30,
	}

	// Set up mock function
	repo.findOneFunc = func(ctx context.Context, filter interface{}) (*TestEntity, error) {
		// Basic filter check
		if f, ok := filter.(bson.M); ok {
			if name, ok := f["name"]; ok && name == "Test User" {
				return entity, nil
			}
		}
		return nil, nil // Not found, but no error
	}

	// Test success case
	result, err := repo.FindOne(context.Background(), bson.M{"name": "Test User"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Name, result.Name)

	// Test not found case - should return nil, nil
	result, err = repo.FindOne(context.Background(), bson.M{"name": "Non-existent"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestRepositoryInsert(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test entity
	entity := &TestEntity{
		ID:   primitive.NewObjectID(),
		Name: "Test User",
		Age:  30,
	}

	// Error entity
	errorEntity := &TestEntity{
		ID:   primitive.NewObjectID(),
		Name: "Error User",
		Age:  25,
	}

	// Set up mock function
	repo.insertFunc = func(ctx context.Context, e *TestEntity) (*TestEntity, error) {
		if e.Name == "Error User" {
			return nil, mongo.ErrClientDisconnected
		}
		return e, nil
	}

	// Test success case
	result, err := repo.Insert(context.Background(), entity)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Name, result.Name)
	assert.Equal(t, entity.Age, result.Age)

	// Test error case
	result, err = repo.Insert(context.Background(), errorEntity)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRepositoryUpdate(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test entity
	id := primitive.NewObjectID()
	entity := &TestEntity{
		ID:   id,
		Name: "Updated User",
		Age:  35,
	}

	// Set up mock function
	repo.updateFunc = func(ctx context.Context, testID primitive.ObjectID, e *TestEntity) (*TestEntity, error) {
		if testID == id {
			return e, nil
		}
		return nil, sharedErrors.ErrNotFound
	}

	// Test success case
	result, err := repo.Update(context.Background(), id, entity)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Name, result.Name)

	// Test not found case
	result, err = repo.Update(context.Background(), primitive.NewObjectID(), entity)
	assert.Error(t, err)
	assert.True(t, sharedErrors.IsNotFoundError(err))
	assert.Nil(t, result)
}

func TestRepositoryFindByIDWithFilters(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test entity
	id := primitive.NewObjectID()
	entity := &TestEntity{
		ID:   id,
		Name: "Test User",
		Age:  30,
	}

	// Set up mock function
	repo.findByIDWithFiltersFunc = func(ctx context.Context, testID primitive.ObjectID, filters bson.M) (*TestEntity, error) {
		if testID == id {
			// Check if age filter matches
			if age, ok := filters["age"]; ok && age == 30 {
				return entity, nil
			}
			// If age filter doesn't match
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.ErrNotFound
	}

	// Test success case - ID and filter match
	result, err := repo.FindByIDWithFilters(context.Background(), id, bson.M{"age": 30})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Name, result.Name)
	assert.Equal(t, entity.Age, result.Age)

	// Test not found case - ID matches but filter doesn't
	result, err = repo.FindByIDWithFilters(context.Background(), id, bson.M{"age": 25})
	assert.Error(t, err)
	assert.True(t, sharedErrors.IsNotFoundError(err))
	assert.Nil(t, result)

	// Test not found case - ID doesn't match
	result, err = repo.FindByIDWithFilters(context.Background(), primitive.NewObjectID(), bson.M{"age": 30})
	assert.Error(t, err)
	assert.True(t, sharedErrors.IsNotFoundError(err))
	assert.Nil(t, result)
}

func TestRepositoryDelete(t *testing.T) {
	// Create test repository
	repo := NewTestRepository()
	
	// Create a test ID
	id := primitive.NewObjectID()

	// Set up mock function
	repo.deleteFunc = func(ctx context.Context, testID primitive.ObjectID) error {
		if testID == id {
			return nil
		}
		return sharedErrors.ErrNotFound
	}

	// Test success case
	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)

	// Test not found case
	err = repo.Delete(context.Background(), primitive.NewObjectID())
	assert.Error(t, err)
	assert.True(t, sharedErrors.IsNotFoundError(err))
}