package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph"
	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMatchUp is the resolver for the getMatchUp field.
func (r *queryResolver) GetMatchUp(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error) {
	panic(fmt.Errorf("not implemented: GetMatchUp - getMatchUp"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }