// interfaces/types.go
package interfaces

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/user-service/models"
)

// Resolver is the resolver root
type Resolver interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

// QueryResolver is the query resolver
type QueryResolver interface {
	Me(ctx context.Context) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
}

// MutationResolver is the mutation resolver
type MutationResolver interface {
	UpdateUser(ctx context.Context, input model.UpdateUserInput) (*models.User, error)
}
