// internal/services/relationship_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/domain"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/middleware"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
)

type relationshipService struct {
	repo repository.RelationshipRepository
}

func NewRelationshipService(repo repository.RelationshipRepository) RelationshipService {
	return &relationshipService{repo: repo}
}

// Helper to create a "not implemented" error with function name
func notImplemented(funcName string) error {
	return errors.New(funcName + " not implemented")
}

// Friendships
func (s *relationshipService) SendFriendRequest(ctx context.Context, receiverID string) (bool, error) {
	requesterID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	friendship := domain.NewFriendship(requesterID, receiverID)

	if err := s.repo.Create(ctx, friendship); err != nil {
		return false, fmt.Errorf("failed to create friendship: %w", err)
	}

	return true, nil
}

func (s *relationshipService) AcceptFriendRequest(ctx context.Context, friendshipID string) (bool, error) {
	fields := map[string]interface{}{
		"status":    domain.RelationshipStatusActive,
		"updatedAt": time.Now().UTC(),
	}

	if err := s.repo.Update(ctx, friendshipID, fields); err != nil {
		return false, fmt.Errorf("failed to accept friend request: %w", err)
	}

	return true, nil
}

func (s *relationshipService) RejectFriendRequest(ctx context.Context, friendshipID string) (bool, error) {
	if err := s.repo.Delete(ctx, friendshipID); err != nil {
		return false, fmt.Errorf("failed to reject friend request: %w", err)
	}
	return true, nil
}

func (s *relationshipService) CancelFriendRequest(ctx context.Context, friendshipID string) (bool, error) {
	if err := s.repo.Delete(ctx, friendshipID); err != nil {
		return false, fmt.Errorf("failed to cancel friend request: %w", err)
	}
	return true, nil
}

func (s *relationshipService) EndFriendship(ctx context.Context, friendshipID string) (bool, error) {
	if err := s.repo.Delete(ctx, friendshipID); err != nil {
		return false, fmt.Errorf("failed to end friendship: %w", err)
	}
	return true, nil
}

// Coachships
func (s *relationshipService) SendCoachRequest(ctx context.Context, userID string) (*model.Coachship, error) {
	return nil, notImplemented("SendCoachRequest")
}

func (s *relationshipService) SendCoacheeRequest(ctx context.Context, userID string) (*model.Coachship, error) {
	return nil, notImplemented("SendCoacheeRequest")
}

func (s *relationshipService) AcceptCoachRequest(ctx context.Context, coachshipID string) (*model.Coachship, error) {
	return nil, notImplemented("AcceptCoachRequest")
}

func (s *relationshipService) AcceptCoacheeRequest(ctx context.Context, coachshipID string) (*model.Coachship, error) {
	return nil, notImplemented("AcceptCoacheeRequest")
}

func (s *relationshipService) DeclineCoachRequest(ctx context.Context, coachshipID string) (bool, error) {
	return false, notImplemented("DeclineCoachRequest")
}

func (s *relationshipService) DeclineCoacheeRequest(ctx context.Context, coachshipID string) (bool, error) {
	return false, notImplemented("DeclineCoacheeRequest")
}

func (s *relationshipService) CancelCoachRequest(ctx context.Context, coachshipID string) (bool, error) {
	return false, notImplemented("CancelCoachRequest")
}

func (s *relationshipService) CancelCoacheeRequest(ctx context.Context, coachshipID string) (bool, error) {
	return false, notImplemented("CancelCoacheeRequest")
}

func (s *relationshipService) EndCoachship(ctx context.Context, coachshipID string) (bool, error) {
	return false, notImplemented("EndCoachship")
}

// Queries - Coachships
func (s *relationshipService) GetCoachship(ctx context.Context, id string) (*model.Coachship, error) {
	return nil, notImplemented("GetCoachship")
}

func (s *relationshipService) ListCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	return nil, notImplemented("ListCoaches")
}

func (s *relationshipService) ListStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	return nil, notImplemented("ListStudents")
}

func (s *relationshipService) ListSentCoacheeRequests(ctx context.Context) ([]*model.Coachship, error) {
	return nil, notImplemented("ListSentCoacheeRequests")
}

func (s *relationshipService) ListReceivedCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	return nil, notImplemented("ListReceivedCoachRequests")
}

func (s *relationshipService) ListSentCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	return nil, notImplemented("ListSentCoachRequests")
}

func (s *relationshipService) ListReceivedCoacheeRequests(ctx context.Context) ([]*model.Coachship, error) {
	return nil, notImplemented("ListReceivedCoacheeRequests")
}

func (s *relationshipService) CheckCoachshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error) {
	return nil, notImplemented("CheckCoachshipStatus")
}

// Queries - Friendships
func (s *relationshipService) GetFriendship(ctx context.Context, id string) (*model.Friendship, error) {
	return nil, notImplemented("GetFriendship")
}

func (s *relationshipService) ListFriends(ctx context.Context, limit int, offset int) ([]*model.Friendship, error) {
	return nil, notImplemented("ListFriends")
}

func (s *relationshipService) ListPendingFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
	return nil, notImplemented("ListPendingFriendRequests")
}

func (s *relationshipService) ListSentFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
	return nil, notImplemented("ListSentFriendRequests")
}

func (s *relationshipService) CheckFriendshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error) {
	return nil, notImplemented("CheckFriendshipStatus")
}
