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

	// Check allowed combinations
	if conditions.Existence != nil && conditions.NonExistence != nil {
		return nil, errors.New("cannot have both existence and non-existence conditions simultaneously")
	}

	// Validate role conditions
	if conditions.Roles != nil && roleFilter != nil {
		fmt.Printf("DEBUG: Role filter: %+v\n", roleFilter) // Print the role filter
		count, err := RelationshipRepo.Count(ctx, roleFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to check role conditions: %w", err)
		}
		if count == 0 {
			return nil, errors.New("role conditions not satisfied")
		}
	}

	// Validate existence conditions
	if conditions.Existence != nil && existenceFilter != nil {
		fmt.Printf("DEBUG: Existence filter: %+v\n", existenceFilter) // Print the existence filter
		count, err := RelationshipRepo.Count(ctx, existenceFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to check existence conditions: %w", err)
		}
		if count == 0 {
			return nil, errors.New("existence condition not met")
		}
	}

	// Validate non-existence conditions
	if conditions.NonExistence != nil && nonExistenceFilter != nil {
		fmt.Printf("DEBUG: Non-Existence filter: %+v\n", nonExistenceFilter) // Print the non-existence filter
		count, err := RelationshipRepo.Count(ctx, nonExistenceFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to check non-existence conditions: %w", err)
		}
		if count > 0 {
			return nil, errors.New("non-existence condition not met")
		}
	}

	// If we reach here, conditions are satisfied
	return next(ctx)
}
