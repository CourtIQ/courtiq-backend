package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipService implements RelationshipServiceIntf
type RelationshipService struct {
	friendshipRepo repository.FriendshipRepository
	coachshipRepo  repository.CoachshipRepository
}

// NewRelationshipService creates a new RelationshipService
func NewRelationshipService(friendshipRepo repository.FriendshipRepository, coachshipRepo repository.CoachshipRepository) *RelationshipService {
	return &RelationshipService{
		friendshipRepo: friendshipRepo,
		coachshipRepo:  coachshipRepo,
	}
}

// ---------------------------------------------------------------------------
// Entity Resolvers
// ---------------------------------------------------------------------------

// FindCoachshipByID finds a coachship by ID
func (s *RelationshipService) FindCoachshipByID(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	return s.coachshipRepo.FindByID(ctx, id)
}

// FindFriendshipByID finds a friendship by ID
func (s *RelationshipService) FindFriendshipByID(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	return s.friendshipRepo.FindByID(ctx, id)
}

// ---------------------------------------------------------------------------
// Friendship Queries
// ---------------------------------------------------------------------------

// GetFriendship gets the friendship between the current user and another user
func (s *RelationshipService) GetFriendship(ctx context.Context, friendID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	friendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, friendID)
	if err != nil {
		return nil, utils.NewFriendshipNotFoundError()
	}

	return friendship, nil
}

// CheckFriendshipStatus checks the friendship status between the current user and another user
func (s *RelationshipService) CheckFriendshipStatus(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return model.RelationshipStatusNone, err
	}

	friendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err != nil {
		return model.RelationshipStatusNone, nil
	}

	return friendship.Status, nil
}

// GetMyFriends gets all friends of the current user
func (s *RelationshipService) GetMyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.friendshipRepo.GetFriendships(ctx, currentUserID, model.RelationshipStatusAccepted, limit, offset)
}

// GetSentFriendRequests gets all friend requests sent by the current user
func (s *RelationshipService) GetSentFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.friendshipRepo.GetSentRequests(ctx, currentUserID, limit, offset)
}

// GetReceivedFriendRequests gets all friend requests received by the current user
func (s *RelationshipService) GetReceivedFriendRequests(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.friendshipRepo.GetReceivedRequests(ctx, currentUserID, limit, offset)
}

// ---------------------------------------------------------------------------
// Friendship Mutations
// ---------------------------------------------------------------------------

// SendFriendRequest sends a friend request to another user
func (s *RelationshipService) SendFriendRequest(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Check if trying to friend self
	if currentUserID == userID {
		return nil, utils.NewSelfRequestError("friend")
	}

	// Check if a friendship already exists
	existingFriendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err == nil && existingFriendship != nil {
		if existingFriendship.Status == model.RelationshipStatusPending {
			return nil, utils.NewRelationshipAlreadyExistsError("FRIENDSHIP")
		}
		if existingFriendship.Status == model.RelationshipStatusAccepted {
			return nil, utils.NewUIError("FRIENDSHIP_ALREADY_ACCEPTED", "You are already friends with this user", nil)
		}
		if existingFriendship.Status == model.RelationshipStatusBlocked {
			return nil, utils.NewRelationshipForbiddenError("Unable to send friend request")
		}
	}

	// Create a new friendship
	now := time.Now()
	friendship := &model.Friendship{
		Type:        model.RelationshipTypeFriendship,
		Status:      model.RelationshipStatusPending,
		InitiatorID: currentUserID,
		ReceiverID:  userID,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	return s.friendshipRepo.Create(ctx, friendship)
}

// AcceptFriendRequest accepts a friend request
func (s *RelationshipService) AcceptFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Get the friendship
	friendship, err := s.friendshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, utils.NewFriendshipNotFoundError()
	}

	// Verify that the current user is the receiver of the request
	if friendship.ReceiverID != currentUserID {
		return nil, utils.NewRelationshipForbiddenError("You can only accept friend requests sent to you")
	}

	// Verify that the friendship is pending
	if friendship.Status != model.RelationshipStatusPending {
		return nil, utils.NewUIError("INVALID_REQUEST_STATUS", "You can only accept pending friend requests", nil)
	}

	// Update the friendship status
	now := time.Now()
	friendship.Status = model.RelationshipStatusAccepted
	friendship.UpdatedAt = &now

	return s.friendshipRepo.Update(ctx, friendship)
}

// RejectFriendRequest rejects a friend request
func (s *RelationshipService) RejectFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Get the friendship
	friendship, err := s.friendshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, utils.NewFriendshipNotFoundError()
	}

	// Verify that the current user is the receiver of the request
	if friendship.ReceiverID != currentUserID {
		return nil, utils.NewRelationshipForbiddenError("You can only reject friend requests sent to you")
	}

	// Verify that the friendship is pending
	if friendship.Status != model.RelationshipStatusPending {
		return nil, utils.NewUIError("INVALID_REQUEST_STATUS", "You can only reject pending friend requests", nil)
	}

	// Update the friendship status
	now := time.Now()
	friendship.Status = model.RelationshipStatusDeclined
	friendship.UpdatedAt = &now

	return s.friendshipRepo.Update(ctx, friendship)
}

// CancelFriendRequest cancels a friend request
func (s *RelationshipService) CancelFriendRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Get the friendship
	friendship, err := s.friendshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, utils.NewFriendshipNotFoundError()
	}

	// Verify that the current user is the initiator of the request
	if friendship.InitiatorID != currentUserID {
		return nil, utils.NewRelationshipForbiddenError("You can only cancel friend requests you sent")
	}

	// Verify that the friendship is pending
	if friendship.Status != model.RelationshipStatusPending {
		return nil, utils.NewUIError("INVALID_REQUEST_STATUS", "You can only cancel pending friend requests", nil)
	}

	// Update the friendship status
	now := time.Now()
	friendship.Status = model.RelationshipStatusDeclined
	friendship.UpdatedAt = &now

	return s.friendshipRepo.Update(ctx, friendship)
}

// RemoveFriend removes a friend
func (s *RelationshipService) RemoveFriend(ctx context.Context, friendID primitive.ObjectID) (bool, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return false, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Get the friendship
	friendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, friendID)
	if err != nil {
		return false, utils.NewFriendshipNotFoundError()
	}

	// Verify that the friendship is accepted
	if friendship.Status != model.RelationshipStatusAccepted {
		return false, utils.NewUIError("NOT_FRIENDS", "You are not friends with this user", nil)
	}

	// Delete the friendship
	return s.friendshipRepo.Delete(ctx, friendship.ID)
}

// BlockUser blocks a user
func (s *RelationshipService) BlockUser(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Check if trying to block self
	if currentUserID == userID {
		return nil, utils.NewSelfRequestError("block")
	}

	// Check if a friendship already exists
	existingFriendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err == nil && existingFriendship != nil {
		now := time.Now()
		existingFriendship.Status = model.RelationshipStatusBlocked
		existingFriendship.UpdatedAt = &now
		return s.friendshipRepo.Update(ctx, existingFriendship)
	}

	// Create a new friendship with blocked status
	now := time.Now()
	friendship := &model.Friendship{
		Type:        model.RelationshipTypeFriendship,
		Status:      model.RelationshipStatusBlocked,
		InitiatorID: currentUserID,
		ReceiverID:  userID,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	return s.friendshipRepo.Create(ctx, friendship)
}

// UnblockUser unblocks a user
func (s *RelationshipService) UnblockUser(ctx context.Context, userID primitive.ObjectID) (*model.Friendship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	// Find the blocked relationship
	friendship, err := s.friendshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err != nil {
		return nil, utils.NewFriendshipNotFoundError()
	}

	// Verify that the relationship is blocked
	if friendship.Status != model.RelationshipStatusBlocked {
		return nil, utils.NewUIError("NOT_BLOCKED", "User is not blocked", nil)
	}

	// Verify that the current user is the one who blocked the other user
	if friendship.InitiatorID != currentUserID {
		return nil, utils.NewRelationshipForbiddenError("You cannot unblock this user")
	}

	// Update the relationship status
	now := time.Now()
	friendship.Status = model.RelationshipStatusNone
	friendship.UpdatedAt = &now

	return s.friendshipRepo.Update(ctx, friendship)
}

// ---------------------------------------------------------------------------
// Coachship Queries
// ---------------------------------------------------------------------------

// GetCoachship gets the coaching relationship between the current user and another user
func (s *RelationshipService) GetCoachship(ctx context.Context, userID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, utils.NewUIError("UNAUTHORIZED", "Unable to determine user identity", err)
	}

	coachship, err := s.coachshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err != nil {
		return nil, utils.NewCoachshipNotFoundError()
	}

	return coachship, nil
}

// IsCoachOf checks if the current user is a coach of the given user
func (s *RelationshipService) IsCoachOf(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return model.RelationshipStatusNone, err
	}

	// Find the coachship
	coachship, err := s.coachshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err != nil {
		return model.RelationshipStatusNone, nil
	}

	// Check if current user is the coach
	if coachship.Coach == currentUserID && coachship.Student == userID {
		return coachship.Status, nil
	}

	return model.RelationshipStatusNone, nil
}

// IsStudentOf checks if the current user is a student of the given user
func (s *RelationshipService) IsStudentOf(ctx context.Context, userID primitive.ObjectID) (model.RelationshipStatus, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return model.RelationshipStatusNone, err
	}

	// Find the coachship
	coachship, err := s.coachshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err != nil {
		return model.RelationshipStatusNone, nil
	}

	// Check if current user is the student
	if coachship.Student == currentUserID && coachship.Coach == userID {
		return coachship.Status, nil
	}

	return model.RelationshipStatusNone, nil
}

// GetCoachships gets all coaching relationships for the current user
func (s *RelationshipService) GetCoachships(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetCoachships(ctx, currentUserID, model.RelationshipStatusAccepted, limit, offset)
}

// GetMyCoaches gets all coaches for the current user
func (s *RelationshipService) GetMyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetCoaches(ctx, currentUserID, limit, offset)
}

// GetMyStudents gets all students for the current user
func (s *RelationshipService) GetMyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetStudents(ctx, currentUserID, limit, offset)
}

// GetSentCoachRequests gets all coaching requests sent by the current user as a coach
func (s *RelationshipService) GetSentCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetSentCoachRequests(ctx, currentUserID, limit, offset)
}

// GetReceivedCoachRequests gets all coaching requests received by the current user as a coach
func (s *RelationshipService) GetReceivedCoachRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetReceivedCoachRequests(ctx, currentUserID, limit, offset)
}

// GetSentStudentRequests gets all coaching requests sent by the current user as a student
func (s *RelationshipService) GetSentStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetSentStudentRequests(ctx, currentUserID, limit, offset)
}

// GetReceivedStudentRequests gets all coaching requests received by the current user as a student
func (s *RelationshipService) GetReceivedStudentRequests(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.coachshipRepo.GetReceivedStudentRequests(ctx, currentUserID, limit, offset)
}

// ---------------------------------------------------------------------------
// Coachship Mutations
// ---------------------------------------------------------------------------

// RequestToBeCoachOf sends a request to be a coach of another user
func (s *RelationshipService) RequestToBeCoachOf(ctx context.Context, userID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if a coachship already exists
	existingCoachship, err := s.coachshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err == nil && existingCoachship != nil {
		return nil, fmt.Errorf("a coaching relationship already exists between you and this user")
	}

	// Create a new coachship
	now := time.Now()
	coachship := &model.Coachship{
		Type:        model.RelationshipTypeCoachship,
		Status:      model.RelationshipStatusPending,
		InitiatorID: currentUserID,
		ReceiverID:  userID,
		CreatedAt:   now,
		UpdatedAt:   &now,
		Coach:       currentUserID, // Current user wants to be the coach
		Student:     userID,        // Target user will be the student
	}

	return s.coachshipRepo.Create(ctx, coachship)
}

// RequestToBeCoachedBy sends a request to be coached by another user
func (s *RelationshipService) RequestToBeCoachedBy(ctx context.Context, userID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check if a coachship already exists
	existingCoachship, err := s.coachshipRepo.FindBetweenUsers(ctx, currentUserID, userID)
	if err == nil && existingCoachship != nil {
		return nil, fmt.Errorf("a coaching relationship already exists between you and this user")
	}

	// Create a new coachship
	now := time.Now()
	coachship := &model.Coachship{
		Type:        model.RelationshipTypeCoachship,
		Status:      model.RelationshipStatusPending,
		InitiatorID: currentUserID,
		ReceiverID:  userID,
		CreatedAt:   now,
		UpdatedAt:   &now,
		Coach:       userID,        // Target user will be the coach
		Student:     currentUserID, // Current user wants to be the student
	}

	return s.coachshipRepo.Create(ctx, coachship)
}

// AcceptToBeCoachOf accepts a request to be a coach
func (s *RelationshipService) AcceptToBeCoachOf(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the coach and receiver of the request
	if coachship.Coach != currentUserID || coachship.ReceiverID != currentUserID {
		return nil, fmt.Errorf("you can only accept requests to be a coach that were sent to you")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only accept pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusAccepted
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// RejectToBeCoachOf rejects a request to be a coach
func (s *RelationshipService) RejectToBeCoachOf(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the coach and receiver of the request
	if coachship.Coach != currentUserID || coachship.ReceiverID != currentUserID {
		return nil, fmt.Errorf("you can only reject requests to be a coach that were sent to you")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only reject pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusDeclined
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// AcceptToBeCoachedBy accepts a request to be coached
func (s *RelationshipService) AcceptToBeCoachedBy(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the student and receiver of the request
	if coachship.Student != currentUserID || coachship.ReceiverID != currentUserID {
		return nil, fmt.Errorf("you can only accept requests to be coached that were sent to you")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only accept pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusAccepted
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// RejectToBeCoachedBy rejects a request to be coached
func (s *RelationshipService) RejectToBeCoachedBy(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the student and receiver of the request
	if coachship.Student != currentUserID || coachship.ReceiverID != currentUserID {
		return nil, fmt.Errorf("you can only reject requests to be coached that were sent to you")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only reject pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusDeclined
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// CancelCoachRequest cancels a request to be a coach
func (s *RelationshipService) CancelCoachRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the coach and initiator of the request
	if coachship.Coach != currentUserID || coachship.InitiatorID != currentUserID {
		return nil, fmt.Errorf("you can only cancel coach requests you sent as a coach")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only cancel pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusDeclined
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// CancelStudentRequest cancels a request to be a student
func (s *RelationshipService) CancelStudentRequest(ctx context.Context, requestID primitive.ObjectID) (*model.Coachship, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("coaching request not found: %v", err)
	}

	// Verify that the current user is the student and initiator of the request
	if coachship.Student != currentUserID || coachship.InitiatorID != currentUserID {
		return nil, fmt.Errorf("you can only cancel student requests you sent as a student")
	}

	// Verify that the coachship is pending
	if coachship.Status != model.RelationshipStatusPending {
		return nil, fmt.Errorf("you can only cancel pending coaching requests")
	}

	// Update the coachship status
	now := time.Now()
	coachship.Status = model.RelationshipStatusDeclined
	coachship.UpdatedAt = &now

	return s.coachshipRepo.Update(ctx, coachship)
}

// EndCoachingAsCoach ends a coaching relationship as the coach
func (s *RelationshipService) EndCoachingAsCoach(ctx context.Context, coachshipID primitive.ObjectID) (bool, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, coachshipID)
	if err != nil {
		return false, fmt.Errorf("coaching relationship not found: %v", err)
	}

	// Verify that the current user is the coach
	if coachship.Coach != currentUserID {
		return false, fmt.Errorf("you can only end coaching relationships where you are the coach")
	}

	// Verify that the coachship is accepted
	if coachship.Status != model.RelationshipStatusAccepted {
		return false, fmt.Errorf("you can only end active coaching relationships")
	}

	// Delete the coachship
	return s.coachshipRepo.Delete(ctx, coachshipID)
}

// EndCoachingAsStudent ends a coaching relationship as the student
func (s *RelationshipService) EndCoachingAsStudent(ctx context.Context, coachshipID primitive.ObjectID) (bool, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	// Get the coachship
	coachship, err := s.coachshipRepo.FindByID(ctx, coachshipID)
	if err != nil {
		return false, fmt.Errorf("coaching relationship not found: %v", err)
	}

	// Verify that the current user is the student
	if coachship.Student != currentUserID {
		return false, fmt.Errorf("you can only end coaching relationships where you are the student")
	}

	// Verify that the coachship is accepted
	if coachship.Status != model.RelationshipStatusAccepted {
		return false, fmt.Errorf("you can only end active coaching relationships")
	}

	// Delete the coachship
	return s.coachshipRepo.Delete(ctx, coachshipID)
}
