package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.61

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph"
	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return r.Me(ctx)
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return r.GetUser(ctx, id)
}

// IsUsernameAvailable is the resolver for the isUsernameAvailable field.
func (r *queryResolver) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	return r.IsUsernameAvailable(ctx, username)
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
