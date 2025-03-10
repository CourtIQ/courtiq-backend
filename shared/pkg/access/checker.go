package access

import (
	"context"
	"errors"
	"time"
)

// Define the context key for storing the access checker
type contextKey string
const CheckerContextKey contextKey = "accessChecker"

// AccessLevel defines the access permission level
type AccessLevel string

// Predefined access levels
const (
	AccessLevelPublic            AccessLevel = "PUBLIC"
	AccessLevelPrivate           AccessLevel = "PRIVATE"
	AccessLevelFriends           AccessLevel = "FRIENDS"
	AccessLevelCoaches           AccessLevel = "COACHES"
	AccessLevelClubMembers       AccessLevel = "CLUB_MEMBERS"
	AccessLevelMatchParticipants AccessLevel = "MATCH_PARTICIPANTS"
	AccessLevelDoublesPartners   AccessLevel = "DOUBLES_PARTNERS"
	AccessLevelFriendsAndCoaches AccessLevel = "FRIENDS_AND_COACHES"
)

// Role defines a user's role in a relationship
type Role string

// Predefined roles
const (
	RoleFriend         Role = "FRIEND"
	RoleCoach          Role = "COACH"
	RoleStudent        Role = "STUDENT"
	RoleClubMember     Role = "CLUB_MEMBER"
	RoleClubAdmin      Role = "CLUB_ADMIN"
	RoleClubCoach      Role = "CLUB_COACH"
	RoleMatchPlayer    Role = "MATCH_PLAYER"
	RoleMatchTracker   Role = "MATCH_TRACKER"
	RoleDoublesPartner Role = "DOUBLES_PARTNER"
	RoleOwner          Role = "OWNER"
	RoleAdmin          Role = "ADMIN"
	RoleViewer         Role = "VIEWER"
)

// Common errors
var (
	ErrAccessDenied    = errors.New("access denied")
	ErrEntityNotFound  = errors.New("entity not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidOwnerID  = errors.New("invalid owner ID")
	ErrInvalidViewerID = errors.New("invalid viewer ID")
)

// AccessResult represents the result of an access check
type AccessResult struct {
	HasAccess    bool        // Whether access is granted
	AccessLevel  AccessLevel // The level of access granted
	Roles        []Role      // Roles the viewer has in relation to the resource
	RelationInfo interface{} // Optional relationship information
	CachedResult bool        // Whether this result came from cache
	ExpiresAt    time.Time   // When this result should expire
}

// CheckConfig holds settings for an access check
type CheckConfig struct {
	RequiredLevel AccessLevel // Minimum level required for access
	AllowedRoles  []Role      // Specific roles that grant access
	EntityID      string      // Optional entity ID for entity-specific checks
	EntityType    string      // Type of entity (e.g., "MATCH", "CLUB")
}

// Checker defines the interface for checking access permissions
type Checker interface {
	// CheckAccess checks if a viewer has required access to an owner's data
	CheckAccess(ctx context.Context, ownerID, viewerID string, config CheckConfig) (*AccessResult, error)

	// HasRole checks if a user has a specific role in relation to an entity
	HasRole(ctx context.Context, userID string, entityID string, role Role) (bool, error)

	// GetRoles gets all roles a user has in relation to an entity
	GetRoles(ctx context.Context, userID string, entityID string) ([]Role, error)

	// ClearCache clears cached access results
	ClearCache(userIDs ...string)
}