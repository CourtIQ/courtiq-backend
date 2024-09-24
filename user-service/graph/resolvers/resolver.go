package resolvers

import (
	"github.com/CourtIQ/backend-courtiq/user-service/graph/generated"
)

// Resolver struct should implement the ResolverRoot interface
type Resolver struct{}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Implement your query and mutation methods here
