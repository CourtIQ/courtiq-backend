package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
)

// Owner is the resolver for the owner field.
func (r *tennisRacketResolver) Owner(ctx context.Context, obj *model.TennisRacket, federationRequires map[string]interface{}) (*model.User, error) {
	return &model.User{
		ID: obj.OwnerID, // Assuming TennisRacket has an OwnerID field
	}, nil
}

// TennisRacket returns graph.TennisRacketResolver implementation.
func (r *Resolver) TennisRacket() graph.TennisRacketResolver { return &tennisRacketResolver{r} }

type tennisRacketResolver struct{ *Resolver }
