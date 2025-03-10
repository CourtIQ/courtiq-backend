package server

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/access"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
)

// DirectiveFunc is a function that implements a GraphQL directive
type DirectiveFunc func(ctx context.Context, obj interface{}, next graphql.Resolver, args map[string]interface{}) (interface{}, error)

// DirectivesMap is a map of directive names to their implementations
type DirectivesMap map[string]DirectiveFunc

// CreateFieldDirectiveMiddleware creates middleware for a specific directive
func CreateFieldDirectiveMiddleware(directiveName string, fn DirectiveFunc) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		fieldCtx := graphql.GetFieldContext(ctx)
		directive := fieldCtx.Field.Definition.Directives.ForName(directiveName)
		if directive == nil {
			return next(ctx)
		}

		// Convert directive arguments to map
		args := make(map[string]interface{})
		for _, arg := range directive.Arguments {
			args[arg.Name] = arg.Value.Value
		}

		return fn(ctx, fieldCtx.Object, next, args)
	}
}

// GetDefaultDirectives returns common directives used across services
func GetDefaultDirectives() DirectivesMap {
	return DirectivesMap{
		"accessControl": AccessControlDirective,
	}
}

// AccessControlDirective implements the @accessControl directive
func AccessControlDirective(ctx context.Context, obj interface{}, next graphql.Resolver, args map[string]interface{}) (interface{}, error) {
	// Get access checker from context
	checker, ok := middleware.GetAccessChecker(ctx)
	if !ok {
		// If no checker available, allow access
		return next(ctx)
	}

	// Get required level from directive args
	requiredLevel, ok := args["requiredLevel"].(string)
	if !ok {
		// Default to public if not specified
		requiredLevel = string(access.AccessLevelPublic)
	}

	// Get owner field from directive args
	ownerField, ok := args["ownerField"].(string)
	if !ok {
		ownerField = "id" // Default to "id"
	}

	// Get owner ID from object
	ownerID := ""
	if objMap, ok := obj.(map[string]interface{}); ok {
		if id, ok := objMap[ownerField]; ok {
			ownerID = id.(string)
		}
	}

	// Get current user ID from context
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		// No user in context - default to no access
		redactWithNull, _ := args["redactWithNull"].(bool)
		if redactWithNull {
			return nil, nil
		}
		return nil, access.ErrAccessDenied
	}

	// Check if user has access
	config := access.CheckConfig{
		RequiredLevel: access.AccessLevel(requiredLevel),
	}

	// Check access
	result, err := checker.CheckAccess(ctx, ownerID, currentUserID.Hex(), config)
	if err != nil || !result.HasAccess {
		// Access denied
		redactWithNull, _ := args["redactWithNull"].(bool)
		if redactWithNull {
			return nil, nil
		}
		return nil, access.ErrAccessDenied
	}

	// Access granted
	return next(ctx)
}