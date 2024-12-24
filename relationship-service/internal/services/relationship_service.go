package services

import (
	"context"
	"errors"
	"time"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNotImplemented = errors.New("not implemented")

// RelationshipService is a placeholder struct that implements RelationshipServiceIntf
type RelationshipService struct {
	friendshipRepo repository.FriendshipRepository
	coachshipRepo  repository.CoachshipRepository
}

// NewRelationshipService returns a new instance of RelationshipService
func NewRelationshipService(friendshipRepo repository.FriendshipRepository, coachshipRepo repository.CoachshipRepository) RelationshipServiceIntf {
	return &RelationshipService{
		friendshipRepo: friendshipRepo,
		coachshipRepo:  coachshipRepo,
	}
}

// ---------------------------------------------------------------------------
// Friendship Queries
// ---------------------------------------------------------------------------

func (s *RelationshipService) Friendship(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	friendship, err := s.friendshipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if friendship == nil {
		return nil, utils.ErrFriendshipNotFound
	}

	if err := s.checkFriendshipPermission(friendship, currentUserID); err != nil {
		return nil, err
	}

	return friendship, nil
}

func (s *RelationshipService) MyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"participants": currentUserID,
		"status":       model.RelationshipStatusActive,
	}

	friendships, err := s.friendshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return friendships, nil
}

func (s *RelationshipService) Friends(ctx context.Context, ofUserId primitive.ObjectID, limit *int, offset *int,
) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	checkFilter := bson.M{
		"participants": bson.M{
			"$all": []primitive.ObjectID{currentUserID, ofUserId},
		},
		"status": model.RelationshipStatusActive,
	}

	activeFriendship, err := s.friendshipRepo.Find(ctx, checkFilter, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(activeFriendship) == 0 {
		return nil, utils.ErrFriendshipForbidden
	}

	filter := bson.M{
		"participants": ofUserId,
		"status":       model.RelationshipStatusActive,
	}

	friends, err := s.friendshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *RelationshipService) MyFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"receiverId": currentUserID,
		"status":     model.RelationshipStatusPending,
	}

	friendships, err := s.friendshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return friendships, nil
}

func (s *RelationshipService) SentFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"senderId": currentUserID,
		"status":   model.RelationshipStatusPending,
	}

	friendships, err := s.friendshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return friendships, nil
}

func (s *RelationshipService) FriendshipStatus(ctx context.Context, otherUserId primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"participants": bson.M{
			"$all": []primitive.ObjectID{currentUserID, otherUserId},
		},
	}

	friendships, err := s.friendshipRepo.Find(ctx, filter, nil, nil)
	if err != nil {
		return "", err
	}

	if len(friendships) == 0 {
		return model.RelationshipStatusNone, nil
	}

	return friendships[0].Status, nil
}

// ---------------------------------------------------------------------------
// Coachship Queries
// ---------------------------------------------------------------------------

func (s *RelationshipService) Coach(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	return nil, ErrNotImplemented
}

func (s *RelationshipService) Student(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	return nil, ErrNotImplemented
}

func (s *RelationshipService) MyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"studentId": currentUserID,
		"status":    model.RelationshipStatusActive,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) MyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"coachId": currentUserID,
		"status":  model.RelationshipStatusActive,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) MyStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"coachId": currentUserID,
		"status":  model.RelationshipStatusPending,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) MyCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"studentId": currentUserID,
		"status":    model.RelationshipStatusPending,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) SentCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// "coachId" == current user & status == "PENDING"
	filter := bson.M{
		"coachId": currentUserID,
		"status":  model.RelationshipStatusPending,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) SentStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// "studentId" == current user & status == "PENDING"
	filter := bson.M{
		"studentId": currentUserID,
		"status":    model.RelationshipStatusPending,
	}

	coachships, err := s.coachshipRepo.Find(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return coachships, nil
}

func (s *RelationshipService) IsStudent(ctx context.Context, studentId primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"coachId":   currentUserID,
		"studentId": studentId,
		"status":    model.RelationshipStatusActive,
	}

	coachship, err := s.coachshipRepo.FindOne(ctx, filter)
	if err != nil {
		return "", err
	}

	if coachship == nil {
		return model.RelationshipStatusNone, nil
	}

	return coachship.Status, nil
}

func (s *RelationshipService) IsCoach(ctx context.Context, coachId primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"coachId":   coachId,
		"studentId": currentUserID,
		"status":    model.RelationshipStatusActive,
	}

	coachship, err := s.coachshipRepo.FindOne(ctx, filter)
	if err != nil {
		return "", err
	}

	if coachship == nil {
		return model.RelationshipStatusNone, nil
	}

	return coachship.Status, nil
}

// ---------------------------------------------------------------------------
// Friendship Mutations
// ---------------------------------------------------------------------------

func (s *RelationshipService) SendFriendRequest(ctx context.Context, receiverId primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if currentUserID == receiverId {
		return nil, utils.ErrCannotSendFriendRequestToSelf
	}

	filter := bson.M{
		"participants": bson.M{
			"$all": []primitive.ObjectID{currentUserID, receiverId},
		},
		"status": bson.M{
			"$in": []model.RelationshipStatus{
				model.RelationshipStatusPending,
				model.RelationshipStatusActive,
			},
		},
	}

	existing, err := s.friendshipRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, utils.ErrFriendshipAlreadyExists
	}

	newFriendship := &model.Friendship{
		ID:           primitive.NewObjectID(),
		Participants: []primitive.ObjectID{currentUserID, receiverId},
		Type:         model.RelationshipTypeFriendship,
		Status:       model.RelationshipStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		SenderID:     currentUserID,
		ReceiverID:   receiverId,
	}

	inserted, err := s.friendshipRepo.Insert(ctx, newFriendship)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (s *RelationshipService) AcceptFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        friendshipId,
		"receiverId": currentUserID,
		"status":     model.RelationshipStatusPending,
	}

	update := bson.M{
		"$set": bson.M{
			"status":    model.RelationshipStatusActive,
			"updatedAt": time.Now(),
		},
	}

	updatedFriendship, err := s.friendshipRepo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedFriendship, nil
}

func (s *RelationshipService) RejectFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        friendshipId,
		"receiverId": currentUserID,
		"status":     model.RelationshipStatusPending,
	}

	deletedFriendship, err := s.friendshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedFriendship, nil
}

func (s *RelationshipService) CancelFriendRequest(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Only match a PENDING request where the current user is the sender
	filter := bson.M{
		"_id":      friendshipId,
		"senderId": currentUserID,
		"status":   model.RelationshipStatusPending,
	}

	deletedFriendship, err := s.friendshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedFriendship, nil
}

func (s *RelationshipService) EndFriendship(ctx context.Context, friendshipId primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// If you store participants in an array:
	filter := bson.M{
		"_id":          friendshipId,
		"participants": currentUserID, // ensures the user is in participants
		"status":       model.RelationshipStatusActive,
	}

	deletedFriendship, err := s.friendshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedFriendship, nil
}

// ---------------------------------------------------------------------------
// Coachship Mutations
// ---------------------------------------------------------------------------

func (s *RelationshipService) RequestToBeStudent(ctx context.Context, ofUserId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if currentUserID == ofUserId {
		return nil, utils.ErrCannotSendStudentRequestToSelf
	}

	filter := bson.M{
		"coachId":   ofUserId,
		"studentId": currentUserID,
		"status": bson.M{
			"$in": []model.RelationshipStatus{
				model.RelationshipStatusPending,
				model.RelationshipStatusActive,
			},
		},
	}

	existing, err := s.coachshipRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, utils.ErrAlreadyStudentOrPending
	}

	now := time.Now()
	newCoachship := &model.Coachship{
		ID:           primitive.NewObjectID(),
		Participants: []primitive.ObjectID{ofUserId, currentUserID},
		CoachID:      ofUserId,
		StudentID:    currentUserID,
		Status:       model.RelationshipStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
		Type:         model.RelationshipTypeCoachship,
	}

	inserted, err := s.coachshipRepo.Insert(ctx, newCoachship)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (s *RelationshipService) AcceptStudentRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     coachshipId,
		"coachId": currentUserID,
		"status":  model.RelationshipStatusPending,
	}

	update := bson.M{
		"$set": bson.M{
			"status":    model.RelationshipStatusActive,
			"updatedAt": time.Now(),
		},
	}

	updatedCoachship, err := s.coachshipRepo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedCoachship, nil
}

func (s *RelationshipService) RejectStudentRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     coachshipId,
		"coachId": currentUserID,
		"status":  model.RelationshipStatusPending,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

func (s *RelationshipService) CancelRequestToBeStudent(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":       coachshipId,
		"studentId": currentUserID,
		"status":    model.RelationshipStatusPending,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

func (s *RelationshipService) RemoveStudent(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Only match an ACTIVE coachship where the current user is the coach
	filter := bson.M{
		"_id":     coachshipId,
		"coachId": currentUserID,
		"status":  model.RelationshipStatusActive,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

func (s *RelationshipService) RequestToBeCoach(ctx context.Context, ofUserId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if currentUserID == ofUserId {
		return nil, utils.ErrCannotSendCoachRequestToSelf
	}

	filter := bson.M{
		"coachId":   currentUserID,
		"studentId": ofUserId,
		"status": bson.M{
			"$in": []model.RelationshipStatus{
				model.RelationshipStatusPending,
				model.RelationshipStatusActive,
			},
		},
	}

	existing, err := s.coachshipRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, utils.ErrAlreadyCoachOrPending
	}

	now := time.Now()
	newCoachship := &model.Coachship{
		ID:           primitive.NewObjectID(),
		Participants: []primitive.ObjectID{currentUserID, ofUserId},
		CoachID:      currentUserID,
		StudentID:    ofUserId,
		Status:       model.RelationshipStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
		Type:         model.RelationshipTypeCoachship,
	}

	inserted, err := s.coachshipRepo.Insert(ctx, newCoachship)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (s *RelationshipService) AcceptCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":       coachshipId,
		"studentId": currentUserID,
		"status":    model.RelationshipStatusPending,
	}

	update := bson.M{
		"$set": bson.M{
			"status":    model.RelationshipStatusActive,
			"updatedAt": time.Now(),
		},
	}

	updatedCoachship, err := s.coachshipRepo.FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedCoachship, nil
}

func (s *RelationshipService) RejectCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Only match a PENDING coachship where the current user is the student
	filter := bson.M{
		"_id":       coachshipId,
		"studentId": currentUserID,
		"status":    model.RelationshipStatusPending,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

func (s *RelationshipService) CancelCoachRequest(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Only match a PENDING coachship where the current user is the coach
	filter := bson.M{
		"_id":     coachshipId,
		"coachId": currentUserID,
		"status":  model.RelationshipStatusPending,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

func (s *RelationshipService) RemoveCoach(ctx context.Context, coachshipId primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Only match an ACTIVE coachship where the current user is the student
	filter := bson.M{
		"_id":       coachshipId,
		"studentId": currentUserID,
		"status":    model.RelationshipStatusActive,
	}

	deletedCoachship, err := s.coachshipRepo.FindOneAndDelete(ctx, filter)
	if err != nil {
		return nil, err
	}

	return deletedCoachship, nil
}

// ---------------------------------------------------------------------------
// Entity Resolvers
// ---------------------------------------------------------------------------

func (s *RelationshipService) FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	return nil, ErrNotImplemented
}

func (s *RelationshipService) FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	return nil, ErrNotImplemented
}

func (s *RelationshipService) FindRelationshipByID(ctx context.Context, id primitive.ObjectID) (model.Relationship, error) {
	return nil, ErrNotImplemented
}

func (s *RelationshipService) checkFriendshipPermission(friendship *model.Friendship, userID primitive.ObjectID) error {
	for _, participantID := range friendship.Participants {
		if participantID == userID {
			return nil
		}
	}
	return utils.ErrFriendshipForbidden
}
