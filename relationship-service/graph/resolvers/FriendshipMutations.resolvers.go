package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

// SendFriendRequest is the resolver for the sendFriendRequest field.
func (r *mutationResolver) SendFriendRequest(ctx context.Context, userID string) (*model.Friendship, error) {
	panic(fmt.Errorf("not implemented: SendFriendRequest - sendFriendRequest"))
}

// AcceptFriendRequest is the resolver for the acceptFriendRequest field.
func (r *mutationResolver) AcceptFriendRequest(ctx context.Context, friendshipID string) (*model.Friendship, error) {
	panic(fmt.Errorf("not implemented: AcceptFriendRequest - acceptFriendRequest"))
}

// RejectFriendRequest is the resolver for the rejectFriendRequest field.
func (r *mutationResolver) RejectFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: RejectFriendRequest - rejectFriendRequest"))
}

// CancelFriendRequest is the resolver for the cancelFriendRequest field.
func (r *mutationResolver) CancelFriendRequest(ctx context.Context, friendshipID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: CancelFriendRequest - cancelFriendRequest"))
}

// EndFriendship is the resolver for the endFriendship field.
func (r *mutationResolver) EndFriendship(ctx context.Context, friendshipID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: EndFriendship - endFriendship"))
}