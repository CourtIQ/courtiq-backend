package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.61

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph"
	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTennisRacket is the resolver for the createTennisRacket field.
func (r *mutationResolver) CreateTennisRacket(ctx context.Context, input model.CreateTennisRacketInput) (*model.TennisRacket, error) {
	return r.EquipmentServiceIntf.CreateTennisRacket(ctx, input)
}

// UpdateMyTennisRacket is the resolver for the updateMyTennisRacket field.
func (r *mutationResolver) UpdateMyTennisRacket(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisRacketInput) (*model.TennisRacket, error) {
	return r.EquipmentServiceIntf.UpdateMyTennisRacket(ctx, id, input)
}

// DeleteMyTennisRacket is the resolver for the deleteMyTennisRacket field.
func (r *mutationResolver) DeleteMyTennisRacket(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.EquipmentServiceIntf.DeleteMyTennisRacket(ctx, id)
}

// CreateTennisString is the resolver for the createTennisString field.
func (r *mutationResolver) CreateTennisString(ctx context.Context, input model.CreateTennisStringInput) (*model.TennisString, error) {
	return r.EquipmentServiceIntf.CreateTennisString(ctx, input)
}

// UpdateMyTennisString is the resolver for the updateMyTennisString field.
func (r *mutationResolver) UpdateMyTennisString(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisStringInput) (*model.TennisString, error) {
	return r.EquipmentServiceIntf.UpdateMyTennisString(ctx, id, input)
}

// DeleteMyTennisString is the resolver for the deleteMyTennisString field.
func (r *mutationResolver) DeleteMyTennisString(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.EquipmentServiceIntf.DeleteMyTennisString(ctx, id)
}

// AssignRacketToString is the resolver for the assignRacketToString field.
func (r *mutationResolver) AssignRacketToString(ctx context.Context, racketID primitive.ObjectID, stringID primitive.ObjectID) (*model.TennisString, error) {
	panic(fmt.Errorf("not implemented: AssignRacketToString - assignRacketToString"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
