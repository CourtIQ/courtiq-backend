package repository

import (
	"context"
	"time"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/db"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type coachshipRepository struct {
	coll *mongo.Collection
}

// NewCoachshipRepository creates a new NewCoachshipRepository implementation
func NewCoachshipRepository(mdb *db.MongoDB) CoachshipRepository {
	return &coachshipRepository{
		coll: mdb.GetCollection(db.CoachshipsCollection),
	}
}

func (r *coachshipRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	var coachship model.Coachship
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&coachship)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &coachship, nil
}

func (r *coachshipRepository) GetMyCoaches(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Coachship, error) {
	// User is a student, getting ACTIVE coachships where this user is the student
	filter := bson.M{
		"studentId": userID,
		"status":    "ACTIVE",
	}
	findOpts := utils.BuildFindOptions(limit, offset)
	cursor, err := r.coll.Find(ctx, filter, findOpts)
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

func (r *coachshipRepository) GetMyStudents(ctx context.Context, userID string, limit *int, offset *int) ([]*model.Coachship, error) {
	// User is a coach, getting ACTIVE coachships where this user is the coach
	filter := bson.M{
		"coachId": userID,
		"status":  "ACTIVE",
	}
	findOpts := utils.BuildFindOptions(limit, offset)
	cursor, err := r.coll.Find(ctx, filter, findOpts)
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

func (r *coachshipRepository) GetMyStudentRequests(ctx context.Context, userID string) ([]*model.Coachship, error) {
	// User is a coach, these are PENDING requests from students
	filter := bson.M{
		"coachId": userID,
		"status":  "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *coachshipRepository) GetMyCoachRequests(ctx context.Context, userID string) ([]*model.Coachship, error) {
	// User is a student, these are PENDING requests from coaches
	filter := bson.M{
		"studentId": userID,
		"status":    "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *coachshipRepository) GetSentCoachRequests(ctx context.Context, userID string) ([]*model.Coachship, error) {
	// Requests sent by this user as a coach (coachId = userID and status = PENDING)
	filter := bson.M{
		"coachId": userID,
		"status":  "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *coachshipRepository) GetSentStudentRequests(ctx context.Context, userID string) ([]*model.Coachship, error) {
	// Requests sent by this user as a student (studentId = userID and status = PENDING)
	filter := bson.M{
		"studentId": userID,
		"status":    "PENDING",
	}
	cursor, err := r.coll.Find(ctx, filter)
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

func (r *coachshipRepository) Insert(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	coachship.ID = primitive.NewObjectID()
	now := time.Now()
	coachship.CreatedAt = now
	coachship.UpdatedAt = now
	_, err := r.coll.InsertOne(ctx, coachship)
	if err != nil {
		return nil, err
	}
	return coachship, nil
}

func (r *coachshipRepository) Update(ctx context.Context, coachship *model.Coachship) (*model.Coachship, error) {
	coachship.UpdatedAt = time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": coachship.ID}, bson.M{"$set": coachship})
	if err != nil {
		return nil, err
	}
	return coachship, nil
}

func (r *coachshipRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
