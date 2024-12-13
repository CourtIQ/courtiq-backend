// // internal/services/relationship_service.go
package services

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
// 	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
// 	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/middleware"
// 	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
// )

// type relationshipService struct {
// 	repo repository.RelationshipRepository
// }

// func NewRelationshipService(repo repository.RelationshipRepository) RelationshipService {
// 	return &relationshipService{repo: repo}
// }

// // Helper to create a "not implemented" error with function name
// func notImplemented(funcName string) error {
// 	return errors.New(funcName + " not implemented")
// }

// // =========================
// // Friendships Mutations
// // =========================

// // SendFriendRequest sends a friend request from the current user to the receiver.
// func (s *relationshipService) SendFriendRequest(ctx context.Context, receiverID string) (*bool, error) {
// 	senderID, err := middleware.GetUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	friendship := domain.NewFriendship(senderID, receiverID)

// 	if err := s.repo.Create(ctx, friendship); err != nil {
// 		return nil, err
// 	}

// 	success := true
// 	return &success, nil
// }

// func (s *relationshipService) AcceptFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
// 	fields := map[string]interface{}{
// 		"status":    domain.RelationshipStatusActive,
// 		"updatedAt": time.Now().UTC(),
// 	}

// 	if err := s.repo.Update(ctx, friendshipID, fields); err != nil {
// 		return nil, err
// 	}

// 	success := true
// 	return &success, nil
// }

// func (s *relationshipService) RejectFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
// 	if err := s.repo.Delete(ctx, friendshipID); err != nil {
// 		return nil, err
// 	}

// 	success := true
// 	return &success, nil
// }

// func (s *relationshipService) CancelFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
// 	if err := s.repo.Delete(ctx, friendshipID); err != nil {
// 		return nil, err
// 	}

// 	success := true
// 	return &success, nil
// }

// func (s *relationshipService) EndFriendship(ctx context.Context, friendshipID string) (*bool, error) {
// 	if err := s.repo.Delete(ctx, friendshipID); err != nil {
// 		return nil, err

// 	}

// 	success := true
// 	return &success, nil
// }

// // =========================
// // Friendships Queries
// // =========================
// func (s *relationshipService) GetFriendship(ctx context.Context, id string) (*model.Friendship, error) {
// 	friendship, err := s.repo.GetFriendshipByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	modelFriendship := domainFriendshipToModel(friendship)
// 	if modelFriendship == nil {
// 		return nil, err
// 	}

// 	return modelFriendship, nil
// }

// func (s *relationshipService) ListMyFriends(ctx context.Context, limit int, offset int) ([]*model.Friendship, error) {
// 	return nil, notImplemented("ListMyFriends")
// }

// func (s *relationshipService) ListFriends(ctx context.Context, ofUserID string, limit int, offset int) ([]*model.Friendship, error) {
// 	return nil, notImplemented("ListFriends")
// }

// func (s *relationshipService) ListFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
// 	return nil, notImplemented("ListFriendRequests")
// }

// func (s *relationshipService) ListSentFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
// 	return nil, notImplemented("ListSentFriendRequests")
// }

// func (s *relationshipService) CheckFriendshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error) {
// 	return nil, notImplemented("CheckFriendshipStatus")
// }

// // =========================
// // Coachship Mutations
// // =========================
// func (s *relationshipService) RequestToBeStudent(ctx context.Context, ofUserID string) (*bool, error) {
// 	return nil, notImplemented("RequestToBeStudent")
// }

// func (s *relationshipService) AcceptStudentRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("AcceptStudentRequest")
// }

// func (s *relationshipService) RejectStudentRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("RejectStudentRequest")
// }

// func (s *relationshipService) CancelStudentRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("CancelStudentRequest")
// }

// func (s *relationshipService) RemoveStudent(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("RemoveStudent")
// }

// func (s *relationshipService) RequestToBeCoach(ctx context.Context, ofUserID string) (*bool, error) {
// 	return nil, notImplemented("RequestToBeCoach")
// }

// func (s *relationshipService) AcceptCoachRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("AcceptCoachRequest")
// }

// func (s *relationshipService) RejectCoachRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("RejectCoachRequest")
// }

// func (s *relationshipService) CancelCoachRequest(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("CancelCoachRequest")
// }

// func (s *relationshipService) RemoveCoach(ctx context.Context, coachshipID string) (*bool, error) {
// 	return nil, notImplemented("RemoveCoach")
// }

// // =========================
// // Coachship Queries
// // =========================
// func (s *relationshipService) GetCoachship(ctx context.Context, id string) (*model.Coachship, error) {
// 	return nil, notImplemented("GetCoachship")
// }

// func (s *relationshipService) ListCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListCoaches")
// }

// func (s *relationshipService) ListStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListStudents")
// }

// func (s *relationshipService) ListSentStudentRequests(ctx context.Context) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListSentStudentRequests")
// }

// func (s *relationshipService) ListReceivedCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListReceivedCoachRequests")
// }

// func (s *relationshipService) ListSentCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListSentCoachRequests")
// }

// func (s *relationshipService) ListReceivedStudentRequests(ctx context.Context) ([]*model.Coachship, error) {
// 	return nil, notImplemented("ListReceivedStudentRequests")
// }

// func (s *relationshipService) IsStudent(ctx context.Context, studentID string) (*model.RelationshipStatus, error) {
// 	return nil, notImplemented("IsStudent")
// }

// func (s *relationshipService) IsCoach(ctx context.Context, coachID string) (*model.RelationshipStatus, error) {
// 	return nil, notImplemented("IsCoach")
// }
