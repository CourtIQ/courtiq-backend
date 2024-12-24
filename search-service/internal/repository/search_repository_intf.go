// internal/repository/search_repository.go

package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CourtIQ/courtiq-backend/search-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
)

// SearchRepository defines the interface for searching users.
type SearchRepository interface {
	SearchUsers(ctx context.Context, query string, excludeUserID primitive.ObjectID, limit, offset int) ([]*model.UserSearchResult, error)
}

type searchRepository struct {
	coll *mongo.Collection
}

// NewSearchRepository sets up the repo with a reference to the "users" collection (for now).
func NewSearchRepository(mdb *db.MongoDB) SearchRepository {
	return &searchRepository{
		coll: mdb.GetCollection(db.UsersCollection),
	}
}

func (r *searchRepository) SearchUsers(
	ctx context.Context,
	query string,
	excludeUserID primitive.ObjectID,
	limit, offset int,
) ([]*model.UserSearchResult, error) {

	// 1) Build the pipeline using your pipeline builder
	pipeline := utils.BuildUserSearchPipeline(query, excludeUserID, limit, offset)

	// 2) Aggregate
	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(ctx)

	// 3) Decode
	var docs []model.UserSearchResult
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	// Convert to []*model.UserSearchResult if you prefer pointer slices
	results := make([]*model.UserSearchResult, len(docs))
	for i := range docs {
		results[i] = &docs[i]
	}

	return results, nil
}
