package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
)

// MyTennisString is the resolver for the myTennisString field.
func (r *queryResolver) MyTennisString(ctx context.Context) ([]*model.TennisString, error) {
	panic(fmt.Errorf("not implemented: MyTennisString - myTennisString"))
}

// GetTennisString is the resolver for the getTennisString field.
func (r *queryResolver) GetTennisString(ctx context.Context, id string) (*model.TennisString, error) {
	panic(fmt.Errorf("not implemented: GetTennisString - getTennisString"))
}
