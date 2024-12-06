package resolvers

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

// SendFriendRequest is the resolver for the sendFriendRequest field.
func (r *mutationResolver) SendFriendRequest(ctx context.Context, userID string) (*model.Friendship, error) {
	return r.RelationshipService.SendFriendRequest(ctx, userID)
}

// AcceptFriendRequest is the resolver for the acceptFriendRequest field.
func (r *mutationResolver) AcceptFriendRequest(ctx context.Context, friendshipID string) (*model.Friendship, error) {
	return r.RelationshipService.AcceptFriendRequest(ctx, friendshipID)
}

// RejectFriendRequest is the resolver for the rejectFriendRequest field.
func (r *mutationResolver) RejectFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
	success, err := r.RelationshipService.RejectFriendRequest(ctx, friendshipID)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

// CancelFriendRequest is the resolver for the cancelFriendRequest field.
func (r *mutationResolver) CancelFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
	success, err := r.RelationshipService.CancelFriendRequest(ctx, friendshipID)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

// EndFriendship is the resolver for the endFriendship field.
func (r *mutationResolver) EndFriendship(ctx context.Context, friendshipID string) (*bool, error) {
	success, err := r.RelationshipService.EndFriendship(ctx, friendshipID)
	if err != nil {
		return nil, err
	}
	return &success, nil
}
