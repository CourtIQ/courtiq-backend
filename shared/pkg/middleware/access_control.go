package middleware

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/shared/pkg/access"
)

// AccessControlMiddleware adds access control to the GraphQL context
func AccessControlMiddleware(checker access.Checker) func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		// Add the access checker to the context
		ctx = context.WithValue(ctx, access.CheckerContextKey, checker)
		return next(ctx)
	}
}

// GetAccessChecker retrieves the access checker from the context
func GetAccessChecker(ctx context.Context) (access.Checker, bool) {
	checker, ok := ctx.Value(access.CheckerContextKey).(access.Checker)
	return checker, ok
}

// GetAccessCheckerOrPanic retrieves the access checker or panics if not found
func GetAccessCheckerOrPanic(ctx context.Context) access.Checker {
	checker, ok := GetAccessChecker(ctx)
	if !ok {
		panic("access checker not found in context")
	}
	return checker
}

// CheckOwnerAccess checks if the current user is the owner of a resource
func CheckOwnerAccess(ctx context.Context, ownerID string) bool {
	currentUserID, err := GetMongoIDFromContext(ctx)
	if err != nil {
		return false
	}

	// Convert ObjectID to string for comparison
	return currentUserID.Hex() == ownerID
}

// CheckRoleBasedAccess checks if the user has a role-based access to an entity
func CheckRoleBasedAccess(ctx context.Context, entityID string, requiredRoles []access.Role) (bool, error) {
	checker, ok := GetAccessChecker(ctx)
	if !ok {
		return false, access.ErrAccessDenied
	}

	currentUserID, err := GetMongoIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	// Convert ObjectID to string for GetRoles
	userRoles, err := checker.GetRoles(ctx, currentUserID.Hex(), entityID)
	if err != nil {
		return false, err
	}

	// Check if user has any of the required roles
	userRolesMap := make(map[access.Role]bool)
	for _, role := range userRoles {
		userRolesMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if userRolesMap[requiredRole] {
			return true, nil
		}
	}

	return false, nil
}
