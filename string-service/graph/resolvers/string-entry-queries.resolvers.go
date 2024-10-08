package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.53

import (
	"context"
	"fmt"

	"github.com/CourtIQ/backend-courtiq/string-service/graph"
	"github.com/CourtIQ/backend-courtiq/string-service/graph/model"
)

// GetStringEntry is the resolver for the getStringEntry field.
func (r *queryResolver) GetStringEntry(ctx context.Context, id string) (*model.StringEntry, error) {
	panic(fmt.Errorf("not implemented: GetStringEntry - getStringEntry"))
}

// GetAllStringEntries is the resolver for the getAllStringEntries field.
func (r *queryResolver) GetAllStringEntries(ctx context.Context, userID string) ([]*model.StringEntry, error) {
	panic(fmt.Errorf("not implemented: GetAllStringEntries - getAllStringEntries"))
}

// GetUniqueStringEntries is the resolver for the getUniqueStringEntries field.
func (r *queryResolver) GetUniqueStringEntries(ctx context.Context) ([]*model.StringEntry, error) {
	panic(fmt.Errorf("not implemented: GetUniqueStringEntries - getUniqueStringEntries"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
