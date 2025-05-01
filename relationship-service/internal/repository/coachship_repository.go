package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CoachshipRepositoryImpl implements CoachshipRepository
type CoachshipRepositoryImpl struct {
	baseRepo *repository.BaseRepository[model.Coachship]
}

// NewCoachshipRepository creates a new instance of CoachshipRepository
func NewCoachshipRepository(factory *repository.RepositoryFactory) CoachshipRepository {
	// Use the shared package's NewRepository function to create a base repository
	baseRepo := repository.NewRepository[model.Coachship](factory, RelationshipsCollection)
	return &CoachshipRepositoryImpl{
		baseRepo: baseRepo,
	}
}

// FindByID finds a coachship by its ID
func (r *CoachshipRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	filter := bson.M{
		"type": model.RelationshipTypeCoachship,
	}
	coachship, err := r.baseRepo.FindByIDWithFilters(ctx, id, filter)
	if err != nil {
		return nil, WrapRepositoryError(err, "find by ID", "coachship")
	}
	return coachship, nil
}

// FindBetweenUsers finds a coachship between two users
func (r *CoachshipRepositoryImpl) FindBetweenUsers(ctx context.Context, coachID, studentID primitive.ObjectID) (*model.Coachship, error) {
	filter := bson.M{
		"coach._id":   coachID,
		"student._id": studentID,
		"type":        model.RelationshipTypeCoachship,
	}
	coachship, err := r.baseRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, WrapRepositoryError(err, "find between users", "coachship")
	}
	return coachship, nil
}

// GetCoachships gets all coachships for a user with a specific status
func (r *CoachshipRepositoryImpl) GetCoachships(ctx context.Context, userID primitive.ObjectID, status model.RelationshipStatus, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"coach._id": userID},
			{"student._id": userID},
		},
		"status": status,
		"type":   model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetCoaches gets all coaches for a student
func (r *CoachshipRepositoryImpl) GetCoaches(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"student._id": studentID,
		"status._id":  model.RelationshipStatusAccepted,
		"type":        model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetStudents gets all students for a coach
func (r *CoachshipRepositoryImpl) GetStudents(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"coach._id": coachID,
		"status":    model.RelationshipStatusAccepted,
		"type":      model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetSentRequests gets all coaching requests sent by a user as either a coach or student
func (r *CoachshipRepositoryImpl) GetSentRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"initiator._id": userID,
		"status":        model.RelationshipStatusPending,
		"type":          model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetReceivedRequests gets all coaching requests received by a user as either a coach or student
func (r *CoachshipRepositoryImpl) GetReceivedRequests(ctx context.Context, userID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"receiver._id": userID,
		"status":       model.RelationshipStatusPending,
		"type":         model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetSentCoachRequests gets coaching requests sent by a user as a coach
func (r *CoachshipRepositoryImpl) GetSentCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"coach._id":     coachID,
		"initiator._id": coachID, // User initiated as coach
		"status":        model.RelationshipStatusPending,
		"type":          model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetReceivedCoachRequests gets coaching requests received by a user as a coach
func (r *CoachshipRepositoryImpl) GetReceivedCoachRequests(ctx context.Context, coachID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"coach._id":    coachID,
		"receiver._id": coachID, // User received as coach
		"status":       model.RelationshipStatusPending,
		"type":         model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetSentStudentRequests gets coaching requests sent by a user as a student
func (r *CoachshipRepositoryImpl) GetSentStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"student._id":   studentID,
		"initiator._id": studentID, // User initiated as student
		"status":        model.RelationshipStatusPending,
		"type":          model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// GetReceivedStudentRequests gets coaching requests received by a user as a student
func (r *CoachshipRepositoryImpl) GetReceivedStudentRequests(ctx context.Context, studentID primitive.ObjectID, limit, offset *int) ([]*model.Coachship, error) {
	filter := bson.M{
		"student._id":  studentID,
		"receiver._id": studentID, // User received as student
		"status":       model.RelationshipStatusPending,
		"type":         model.RelationshipTypeCoachship,
	}

	opts := options.Find()
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}
	if offset != nil {
		opts.SetSkip(int64(*offset))
	}

	return r.baseRepo.Find(ctx, filter, opts)
}

// Create creates a new coachship
func (r *CoachshipRepositoryImpl) Create(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	// Set _id if not already set
	if coachship.ID == primitive.NilObjectID {
		coachship.ID = primitive.NewObjectID()
	}

	// Ensure type is set to COACHSHIP
	coachship.Type = model.RelationshipTypeCoachship

	created, err := r.baseRepo.Insert(ctx, coachship)
	if err != nil {
		return nil, WrapRepositoryError(err, "create", "coachship")
	}
	return created.(*model.Coachship), nil
}

// Update updates an existing coachship
func (r *CoachshipRepositoryImpl) Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	updated, err := r.baseRepo.Update(ctx, coachship.ID, coachship)
	if err != nil {
		return nil, WrapRepositoryError(err, "update", "coachship")
	}
	return updated, nil
}

// Delete deletes a coachship
func (r *CoachshipRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	err := r.baseRepo.Delete(ctx, id)
	if err != nil {
		return false, WrapRepositoryError(err, "delete", "coachship")
	}
	return true, nil
}
