package repository

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type coachshipRepository struct {
	collection *mongo.Collection
}

// NewCoachshipRepository returns an implementation of CoachshipRepository backed by MongoDB.
func NewCoachshipRepository(mdb *db.MongoDB) CoachshipRepository {
	return &coachshipRepository{
		collection: mdb.GetCollection(db.CoachshipsCollection),
	}
}

// FindByID looks up a coachship by its ObjectID. If not found, returns ErrCoachshipNotFound.
func (r *coachshipRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	filter := bson.M{"_id": id}

	var coachship model.Coachship
	err := r.collection.FindOne(ctx, filter).Decode(&coachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Return a user-friendly error
			return nil, utils.ErrCoachshipNotFound
		}
		return nil, err // Some other DB error
	}
	return &coachship, nil
}

// Insert creates a new coachship document and returns the inserted record.
func (r *coachshipRepository) Insert(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	if coachship.ID.IsZero() {
		coachship.ID = primitive.NewObjectID()
	}

	// Insert the document
	res, err := r.collection.InsertOne(ctx, coachship)
	if err != nil {
		return nil, err
	}

	// Retrieve the newly inserted document to return
	var insertedCoachship model.Coachship
	if err := r.collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&insertedCoachship); err != nil {
		return nil, err
	}

	return &insertedCoachship, nil
}

// Update modifies an existing coachship document, returning the updated record.
// If ID is missing, returns ErrMissingCoachshipID; if not found, returns ErrCoachshipNotFound.
func (r *coachshipRepository) Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	if coachship.ID.IsZero() {
		return nil, utils.ErrMissingCoachshipID
	}

	filter := bson.M{"_id": coachship.ID}
	update := bson.M{"$set": coachship}

	// Return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCoachship model.Coachship
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCoachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrCoachshipNotFound
		}
		return nil, err
	}

	return &updatedCoachship, nil
}

// Delete removes a coachship by its ObjectID, returning the deleted document.
// If no document was deleted, returns ErrCoachshipNotFound.
func (r *coachshipRepository) Delete(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	filter := bson.M{"_id": id}

	// Find and delete in one step
	var deletedCoachship model.Coachship
	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&deletedCoachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrCoachshipNotFound
		}
		return nil, err
	}

	return &deletedCoachship, nil
}

// Find returns an array of coachships matching the provided filter.
// If none found, it returns an empty slice (and no error).
func (r *coachshipRepository) Find(ctx context.Context, filter interface{}, limit *int, offset *int) ([]*model.Coachship, error) {
	findOpts := utils.BuildFindOptions(limit, offset)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var coachships []*model.Coachship
	if err := cursor.All(ctx, &coachships); err != nil {
		return nil, err
	}
	return coachships, nil
}
func (r *coachshipRepository) FindOne(ctx context.Context, filter interface{}) (*model.Coachship, error) {
	var coachship model.Coachship
	err := r.collection.FindOne(ctx, filter).Decode(&coachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &coachship, nil
}

// FindOneAndUpdate finds a single Friendship that matches `filter`, applies `update`,
func (r *coachshipRepository) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*model.Coachship, error) {

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCoachship model.Coachship
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCoachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrCoachshipNotFound
		}
		return nil, err
	}

	return &updatedCoachship, nil
}

// FindOneAndDelete finds a single Friendship matching `filter`, deletes it,
func (r *coachshipRepository) FindOneAndDelete(ctx context.Context, filter interface{}) (*model.Coachship, error) {
	var deletedCoachship model.Coachship
	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&deletedCoachship)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrCoachshipNotFound
		}
		return nil, err
	}
	return &deletedCoachship, nil
}
