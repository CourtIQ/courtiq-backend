// internal/services/search_service.go

package services

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/search-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
)

// SearchService is an interface to keep your code testable & decoupled.
type SearchService interface {
	Search(ctx context.Context, query string, resourceTypes []model.ResourceType, limit *int, offset *int) ([]model.SearchResult, error)
}

type searchService struct {
	searchRepo repository.SearchRepository
}

// NewSearchService constructs the SearchService with the required repository.
func NewSearchService(searchRepo repository.SearchRepository) SearchService {
	return &searchService{
		searchRepo: searchRepo,
	}
}

// Search implements searching for the given resourceTypes (currently just USER).
func (s *searchService) Search(
	ctx context.Context,
	query string,
	resourceTypes []model.ResourceType,
	limit *int,
	offset *int,
) ([]model.SearchResult, error) {

	// 1) Convert pointers to actual int values (with defaults).
	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}
	offsetVal := 0
	if offset != nil {
		offsetVal = *offset
	}

	// 2) Get current user ID so we can exclude them from user results
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized or missing user: %w", err)
	}

	// 3) Accumulate final results from each resource type
	var finalResults []model.SearchResult

	for _, rt := range resourceTypes {
		switch rt {
		case model.ResourceTypeUser:
			// search USERS
			userResults, err := s.searchRepo.SearchUsers(ctx, query, currentUserID, limitVal, offsetVal)
			if err != nil {
				return nil, err
			}

			// convert []*UserSearchResult -> []SearchResult
			for _, ur := range userResults {
				finalResults = append(finalResults, ur)
			}

		// In the future, handle other resource types (CLUB, MATCHUP, etc.)
		default:
			// For now, do nothing if an unknown or unimplemented resourceType is requested
			// or you could return an error if you prefer strictness.
		}
	}

	return finalResults, nil
}
