package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/relationship-service/graph"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Friendship is the resolver for the friendship field.
func (r *queryResolver) Friendship(ctx context.Context, id primitive.ObjectID) (*model.Friendship, error) {
	return r.RelationshipService.Friendship(ctx, id)
}

// MyFriends is the resolver for the myFriends field.
func (r *queryResolver) MyFriends(ctx context.Context, limit *int, offset *int) ([]*model.Friendship, error) {
	return r.RelationshipService.MyFriends(ctx, limit, offset)
}

// Friends is the resolver for the friends field.
func (r *queryResolver) Friends(ctx context.Context, ofUserID primitive.ObjectID, limit *int, offset *int) ([]*model.Friendship, error) {
	return r.RelationshipService.Friends(ctx, ofUserID, limit, offset)
}

// MyFriendRequests is the resolver for the myFriendRequests field.
func (r *queryResolver) MyFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
	return r.RelationshipService.MyFriendRequests(ctx)
}

// SentFriendRequests is the resolver for the sentFriendRequests field.
func (r *queryResolver) SentFriendRequests(ctx context.Context) ([]*model.Friendship, error) {
	return r.RelationshipService.SentFriendRequests(ctx)
}

// FriendshipStatus is the resolver for the friendshipStatus field.
func (r *queryResolver) FriendshipStatus(ctx context.Context, otherUserID primitive.ObjectID) (*model.RelationshipStatus, error) {
	return r.RelationshipService.FriendshipStatus(ctx, otherUserID)
}

// Coachship is the resolver for the coachship field.
func (r *queryResolver) Coachship(ctx context.Context, id primitive.ObjectID) (*model.Coachship, error) {
	return r.RelationshipService.Coachship(ctx, id)
}

// MyCoaches is the resolver for the myCoaches field.
func (r *queryResolver) MyCoaches(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	return r.RelationshipService.MyCoaches(ctx, limit, offset)
}

// MyStudents is the resolver for the myStudents field.
func (r *queryResolver) MyStudents(ctx context.Context, limit *int, offset *int) ([]*model.Coachship, error) {
	return r.RelationshipService.MyStudents(ctx, limit, offset)
}

// MyStudentRequests is the resolver for the myStudentRequests field.
func (r *queryResolver) MyStudentRequests(ctx context.Context) ([]*model.Coachship, error) {
	return r.RelationshipService.MyStudentRequests(ctx)
}

// MyCoachRequests is the resolver for the myCoachRequests field.
func (r *queryResolver) MyCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	return r.RelationshipService.MyCoachRequests(ctx)
}

// SentCoachRequests is the resolver for the sentCoachRequests field.
func (r *queryResolver) SentCoachRequests(ctx context.Context) ([]*model.Coachship, error) {
	return r.RelationshipService.SentCoachRequests(ctx)
}

// SentStudentRequests is the resolver for the sentStudentRequests field.
func (r *queryResolver) SentStudentRequests(ctx context.Context) ([]*model.Coachship, error) {
	return r.RelationshipService.SentStudentRequests(ctx)
}

// IsStudent is the resolver for the isStudent field.
func (r *queryResolver) IsStudent(ctx context.Context, studentID primitive.ObjectID) (*model.RelationshipStatus, error) {
	return r.RelationshipService.IsStudent(ctx, studentID)
}

// IsCoach is the resolver for the isCoach field.
func (r *queryResolver) IsCoach(ctx context.Context, coachID primitive.ObjectID) (*model.RelationshipStatus, error) {
	return r.RelationshipService.IsCoach(ctx, coachID)
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
