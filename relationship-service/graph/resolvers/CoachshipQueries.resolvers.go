package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
)

// Coachship is the resolver for the coachship field.
func (r *queryResolver) Coachship(ctx context.Context, id string) (*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: Coachship - coachship"))
}

// Coaches is the resolver for the coaches field.
func (r *queryResolver) Coaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: Coaches - coaches"))
}

// Students is the resolver for the students field.
func (r *queryResolver) Students(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: Students - students"))
}

// SentCoacheeRequests is the resolver for the sentCoacheeRequests field.
func (r *queryResolver) SentCoacheeRequests(ctx context.Context) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: SentCoacheeRequests - sentCoacheeRequests"))
}

// ReceivedCoachRequests is the resolver for the receivedCoachRequests field.
func (r *queryResolver) ReceivedCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: ReceivedCoachRequests - receivedCoachRequests"))
}

// SentCoachRequests is the resolver for the sentCoachRequests field.
func (r *queryResolver) SentCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: SentCoachRequests - sentCoachRequests"))
}

// ReceivedCoacheeRequests is the resolver for the receivedCoacheeRequests field.
func (r *queryResolver) ReceivedCoacheeRequests(ctx context.Context) ([]*model.Coachship, error) {
	panic(fmt.Errorf("not implemented: ReceivedCoacheeRequests - receivedCoacheeRequests"))
}

// CoachshipStatus is the resolver for the coachshipStatus field.
func (r *queryResolver) CoachshipStatus(ctx context.Context, otherUserID string) (*model.RelationshipStatus, error) {
	panic(fmt.Errorf("not implemented: CoachshipStatus - coachshipStatus"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }