// internal/services/search_service.go

package services

import (
	"context"
	"fmt"
	"log"

	"github.com/CourtIQ/courtiq-backend/search-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/scalar"
)

// SearchService is an interface to keep your code testable & decoupled.
type SearchService interface {
	Search(ctx context.Context, query string, resourceTypes []model.ResourceType,
		limit *int, offset *int, near *scalar.GeoPoint) ([]model.SearchResult, error)

	SearchUsers(ctx context.Context, query string, limit *int, offset *int) ([]*model.UserSearchResult, error)

	SearchTennisCourts(ctx context.Context, query string, limit *int, offset *int, near *scalar.GeoPoint) ([]*model.TennisCourtSearchResult, error)
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
	near *scalar.GeoPoint,
) ([]model.SearchResult, error) {

	// 1) Convert pointer inputs to actual int values (with defaults).
	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}
	offsetVal := 0
	if offset != nil {
		offsetVal = *offset
	}

	// 2) Get current user ID (for excluding self in user search).
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized or missing user: %w", err)
	}

	// 3) Accumulate final results from each resource type.
	var finalResults []model.SearchResult

	for _, rt := range resourceTypes {
		switch rt {
		// Existing case: searching Users
		case model.ResourceTypeUser:
			userResults, err := s.searchRepo.SearchUsers(ctx, query, currentUserID, limitVal, offsetVal)
			if err != nil {
				return nil, err
			}
			// Append user results to finalResults
			for _, ur := range userResults {
				finalResults = append(finalResults, ur)
			}

		// NEW case: searching Tennis Courts
		case model.ResourceTypeTennisCourts:
			// In your repo, you currently do SearchTennisCourts(ctx, query, lat, lng, radius, limit).
			// If your GraphQL doesn't provide lat/lng, you might supply dummy values or refactor the schema to accept them.
			// For now, let's assume you pass (0, 0) with some large radius (e.g., 25km).
			tennisCourtResults, err := s.searchRepo.SearchTennisCourts(ctx, query, 0, 0, 25000, limitVal, offsetVal)
			if err != nil {
				return nil, err
			}
			// Append court results to finalResults
			for _, tcResult := range tennisCourtResults {
				finalResults = append(finalResults, tcResult)
			}

		default:
			// No-op or return an error if you want to be strict
		}
	}

	return finalResults, nil
}

func (s *searchService) SearchUsers(
	ctx context.Context,
	query string,
	limit *int,
	offset *int,
) ([]*model.UserSearchResult, error) {
	ownerID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Log the incoming parameters (optional).
	log.Printf("[SearchUsers Service] query=%q limit=%v offset=%v", query, limit, offset)

	// Provide defaults if limit/offset are nil.
	l := 10
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}

	// For demonstration, assume we’re not excluding any user, so pass NilObjectID.
	excludeUserID := ownerID

	// Call your repository method.
	// Modify the signature if your real repo has different arguments.
	results, err := s.searchRepo.SearchUsers(ctx, query, excludeUserID, l, o)
	if err != nil {
		return nil, fmt.Errorf("SearchUsers repo error: %w", err)
	}

	// If none found, return empty slice (not nil).
	if len(results) == 0 {
		log.Printf("[SearchUsers Service] No users found for query=%q", query)
		return []*model.UserSearchResult{}, nil
	}

	log.Printf("[SearchUsers Service] Returning %d user(s)", len(results))
	return results, nil
}

func (s *searchService) SearchTennisCourts(
	ctx context.Context,
	query string,
	limit *int,
	offset *int,
	near *scalar.GeoPoint,
) ([]*model.TennisCourtSearchResult, error) {

	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}
	offsetVal := 0
	if offset != nil {
		offsetVal = *offset
	}

	// If the user didn’t supply “near”, we can default lat=0, lng=0 or choose something else
	lat, lng := 0.0, 0.0
	if near != nil {
		lat = near[1]
		lng = near[0]
	}

	// Hard-code a radius (e.g. 25km).
	// Or you could add radius to your schema if you want it to be user-configurable.
	radius := 25000.0

	// Pass it down to the repository
	courts, err := s.searchRepo.SearchTennisCourts(ctx, query, lat, lng, radius, limitVal, offsetVal)
	if err != nil {
		return nil, err
	}

	// offsetVal handling is optional (your repo doesn’t accept offset).
	// You could do manual slicing: courts = courts[offsetVal:]

	return courts, nil

}
