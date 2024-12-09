// directive.go (in the satisfies package)
package satisfies

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
)

// RelationshipRepo is a package-level variable that will be set at runtime (from main) and can be overridden in tests.
var RelationshipRepo repository.RelationshipRepository

// GetCurrentUserID is a package-level variable for getting the current user’s ID from context.
// This can be set in main and replaced in tests with a mock.
var GetCurrentUserID func(ctx context.Context) (string, error)

// SatisfiesDirective uses the RelationshipRepo and GetCurrentUserID variables
// to validate conditions. If they are not set, or if conditions fail, it returns an error.
func SatisfiesDirective(ctx context.Context, obj interface{}, next graphql.Resolver, conditions model.SatisfiesConditions) (interface{}, error) {
	if RelationshipRepo == nil {
		return nil, errors.New("relationship repository not set")
	}
	if GetCurrentUserID == nil {
		return nil, errors.New("user ID retrieval function not set")
	}

	roleFilter, err := BuildRoleFilter(ctx, conditions.Roles)
	if err != nil {
		return nil, fmt.Errorf("error building role filter: %w", err)
	}

	existenceFilter, err := BuildExistenceFilter(ctx, conditions.Existence)
	if err != nil {
		return nil, fmt.Errorf("error building existence filter: %w", err)
	}

	nonExistenceFilter, err := BuildNonExistenceFilter(ctx, conditions.NonExistence)
	if err != nil {
		return nil, fmt.Errorf("error building non-existence filter: %w", err)
	}

	// After building roleFilter, existenceFilter, and nonExistenceFilter:
	if conditions.Existence != nil && conditions.NonExistence != nil {
		return nil, errors.New("cannot have both existence and non-existence conditions simultaneously")
	}

	// Combine role and existence filters if needed
	var combinedFilter map[string]interface{}
	if roleFilter != nil && existenceFilter != nil {
		// Merge or use $and operator as shown above
		combinedFilter = map[string]interface{}{
			"$and": []map[string]interface{}{roleFilter, existenceFilter},
		}
	} else if roleFilter != nil {
		combinedFilter = roleFilter
	} else if existenceFilter != nil {
		combinedFilter = existenceFilter
	}

	// Check combined role + existence conditions if any are set
	if combinedFilter != nil {
		fmt.Printf("DEBUG: Combined filter: %+v\n", combinedFilter)
		count, err := RelationshipRepo.Count(ctx, combinedFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to check combined conditions: %w", err)
		}
		if count == 0 {
			return nil, errors.New("combined conditions not met")
		}
	}

	// If there's nonExistenceFilter, handle it separately as it conflicts with existence filters
	if nonExistenceFilter != nil {
		count, err := RelationshipRepo.Count(ctx, nonExistenceFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to check non-existence conditions: %w", err)
		}
		if count > 0 {
			return nil, errors.New("non-existence condition not met")
		}
	}

	return next(ctx)
}
