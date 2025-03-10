package access

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RelationshipInfo contains details about a relationship between users
type RelationshipInfo struct {
	ID           primitive.ObjectID `bson:"_id"`
	Type         string             `bson:"type"`
	Status       string             `bson:"status"`
	InitiatorID  primitive.ObjectID `bson:"initiatorId"`
	TargetID     primitive.ObjectID `bson:"targetId"`
	Roles        []RoleAssignment   `bson:"roles"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty"`
}

// RoleAssignment represents a role assigned within a relationship
type RoleAssignment struct {
	UserID    primitive.ObjectID     `bson:"userId"`
	Role      string                 `bson:"role"`
	GrantedAt time.Time              `bson:"grantedAt"`
	GrantedBy primitive.ObjectID     `bson:"grantedBy,omitempty"`
	ExpiresAt *time.Time             `bson:"expiresAt,omitempty"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty"`
}

// RelationshipCheckerConfig holds configuration for the relationship checker
type RelationshipCheckerConfig struct {
	Client            *mongo.Client
	Database          string
	RelCollection     string
	EntityCollection  string
	MemberCollection  string
	CacheSize         int
	CacheTTL          time.Duration
}

// DefaultRelationshipCheckerConfig returns default configuration
func DefaultRelationshipCheckerConfig() RelationshipCheckerConfig {
	return RelationshipCheckerConfig{
		RelCollection:    "relationships",
		EntityCollection: "entities",
		MemberCollection: "memberships",
		CacheSize:        10000,
		CacheTTL:         time.Minute * 5,
	}
}

// RelationshipChecker implements the Checker interface with MongoDB
type RelationshipChecker struct {
	client            *mongo.Client
	database          string
	relCollection     string
	entityCollection  string
	memberCollection  string
	cache             *AccessCache
}

// NewRelationshipChecker creates a relationship checker
func NewRelationshipChecker(config RelationshipCheckerConfig) *RelationshipChecker {
	return &RelationshipChecker{
		client:            config.Client,
		database:          config.Database,
		relCollection:     config.RelCollection,
		entityCollection:  config.EntityCollection,
		memberCollection:  config.MemberCollection,
		cache:             NewAccessCache(config.CacheSize, config.CacheTTL),
	}
}

// relationships returns the relationships collection
func (c *RelationshipChecker) relationships() *mongo.Collection {
	return c.client.Database(c.database).Collection(c.relCollection)
}

// CheckAccess implements the Checker interface
func (c *RelationshipChecker) CheckAccess(
	ctx context.Context,
	ownerID,
	viewerID string,
	config CheckConfig,
) (*AccessResult, error) {
	// Check cache first
	if cachedResult := c.cache.Get(ownerID, viewerID, config); cachedResult != nil {
		return cachedResult, nil
	}
	
	// If owner and viewer are the same, always allow access
	if ownerID == viewerID {
		result := &AccessResult{
			HasAccess:   true,
			AccessLevel: AccessLevelPrivate,
			Roles:       []Role{RoleOwner},
			ExpiresAt:   time.Now().Add(24 * time.Hour), // Owner access doesn't change frequently
		}
		c.cache.Set(ownerID, viewerID, config, result)
		return result, nil
	}
	
	// Public data is always accessible
	if config.RequiredLevel == AccessLevelPublic {
		result := &AccessResult{
			HasAccess:   true,
			AccessLevel: AccessLevelPublic,
			Roles:       []Role{RoleViewer},
			ExpiresAt:   time.Now().Add(24 * time.Hour), // Public access doesn't change frequently
		}
		c.cache.Set(ownerID, viewerID, config, result)
		return result, nil
	}
	
	// Private data is only accessible to the owner
	if config.RequiredLevel == AccessLevelPrivate {
		result := &AccessResult{
			HasAccess:   false,
			AccessLevel: AccessLevelPublic, // Default to public level
			Roles:       []Role{},
			ExpiresAt:   time.Now().Add(time.Hour), // Cache negative result for less time
		}
		c.cache.Set(ownerID, viewerID, config, result)
		return result, nil
	}
	
	// Convert string IDs to ObjectIDs
	ownerObjID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, ErrInvalidOwnerID
	}
	
	viewerObjID, err := primitive.ObjectIDFromHex(viewerID)
	if err != nil {
		return nil, ErrInvalidViewerID
	}
	
	// Build the query to find relationship between users
	query := bson.M{
		"$or": []bson.M{
			{
				"initiatorId": ownerObjID,
				"targetId": viewerObjID,
				"status": "ACTIVE",
			},
			{
				"initiatorId": viewerObjID,
				"targetId": ownerObjID,
				"status": "ACTIVE",
			},
		},
	}
	
	// Find all active relationships between these users
	relationshipsCursor, err := c.relationships().Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}
	defer relationshipsCursor.Close(ctx)
	
	relationships := []RelationshipInfo{}
	if err := relationshipsCursor.All(ctx, &relationships); err != nil {
		return nil, fmt.Errorf("failed to decode relationships: %w", err)
	}
	
	// Process relationships to determine access
	result := &AccessResult{
		HasAccess:   false,
		AccessLevel: AccessLevelPublic,
		Roles:       []Role{},
		ExpiresAt:   time.Now().Add(time.Hour),
	}
	
	// Check for blocked relationship first
	for _, rel := range relationships {
		if rel.Status == "BLOCKED" {
			// If either user has blocked the other, deny access
			result.HasAccess = false
			c.cache.Set(ownerID, viewerID, config, result)
			return result, nil
		}
	}
	
	// Extract roles from relationships
	userRoles := make(map[Role]bool)
	for _, rel := range relationships {
		for _, roleAssign := range rel.Roles {
			if roleAssign.UserID == viewerObjID {
				userRoles[Role(roleAssign.Role)] = true
			}
		}
		
		// Infer roles from relationship type
		switch rel.Type {
		case "FRIENDSHIP":
			userRoles[RoleFriend] = true
		case "COACHSHIP":
			// Determine if viewer is coach or student
			if rel.InitiatorID == viewerObjID {
				for _, roleAssign := range rel.Roles {
					if roleAssign.UserID == viewerObjID && roleAssign.Role == string(RoleCoach) {
						userRoles[RoleCoach] = true
					}
				}
			} else {
				for _, roleAssign := range rel.Roles {
					if roleAssign.UserID == viewerObjID && roleAssign.Role == string(RoleStudent) {
						userRoles[RoleStudent] = true
					}
				}
			}
		case "CLUB_MEMBERSHIP":
			userRoles[RoleClubMember] = true
			// Check for admin role
			for _, roleAssign := range rel.Roles {
				if roleAssign.UserID == viewerObjID &&
				(roleAssign.Role == string(RoleClubAdmin) || roleAssign.Role == string(RoleClubCoach)) {
					userRoles[Role(roleAssign.Role)] = true
				}
			}
		case "MATCH_PARTICIPATION":
			userRoles[RoleMatchPlayer] = true
		case "DOUBLES_PARTNERSHIP":
			userRoles[RoleDoublesPartner] = true
		}
	}
	
	// Convert role map to slice
	var roles []Role
	for role := range userRoles {
		roles = append(roles, role)
	}
	result.Roles = roles
	
	// Check if user has any of the allowed roles
	if len(config.AllowedRoles) > 0 {
		for _, allowedRole := range config.AllowedRoles {
			if userRoles[allowedRole] {
				result.HasAccess = true
				break
			}
		}
	}
	
	// Check against required access level
	switch config.RequiredLevel {
	case AccessLevelFriends:
		result.HasAccess = userRoles[RoleFriend]
	case AccessLevelCoaches:
		result.HasAccess = userRoles[RoleCoach]
	case AccessLevelClubMembers:
		result.HasAccess = userRoles[RoleClubMember]
	case AccessLevelMatchParticipants:
		result.HasAccess = userRoles[RoleMatchPlayer] || userRoles[RoleMatchTracker]
	case AccessLevelDoublesPartners:
		result.HasAccess = userRoles[RoleDoublesPartner]
	case AccessLevelFriendsAndCoaches:
		result.HasAccess = userRoles[RoleFriend] || userRoles[RoleCoach]
	}
	
	// Store result in cache
	if result.HasAccess {
		result.AccessLevel = config.RequiredLevel
	}
	c.cache.Set(ownerID, viewerID, config, result)
	
	return result, nil
}

// HasRole checks if a user has a specific role in relation to an entity
func (c *RelationshipChecker) HasRole(
	ctx context.Context,
	userID string,
	entityID string,
	role Role,
) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, ErrInvalidOwnerID
	}
	
	entityObjID, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return false, ErrInvalidOwnerID
	}
	
	// Query for relationship with this role
	query := bson.M{
		"$or": []bson.M{
			{"initiatorId": userObjID, "targetId": entityObjID},
			{"initiatorId": entityObjID, "targetId": userObjID},
		},
		"status": "ACTIVE",
		"roles": bson.M{
			"$elemMatch": bson.M{
				"userId": userObjID,
				"role": string(role),
			},
		},
	}
	
	count, err := c.relationships().CountDocuments(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to query relationships: %w", err)
	}
	
	return count > 0, nil
}

// GetRoles gets all roles a user has in relation to an entity
func (c *RelationshipChecker) GetRoles(
	ctx context.Context,
	userID string,
	entityID string,
) ([]Role, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrInvalidOwnerID
	}
	
	entityObjID, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return nil, ErrInvalidOwnerID
	}
	
	// Query for relationships
	query := bson.M{
		"$or": []bson.M{
			{"initiatorId": userObjID, "targetId": entityObjID},
			{"initiatorId": entityObjID, "targetId": userObjID},
		},
		"status": "ACTIVE",
	}
	
	relationshipsCursor, err := c.relationships().Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}
	defer relationshipsCursor.Close(ctx)
	
	relationships := []RelationshipInfo{}
	if err := relationshipsCursor.All(ctx, &relationships); err != nil {
		return nil, fmt.Errorf("failed to decode relationships: %w", err)
	}
	
	// Extract roles
	roleMap := make(map[Role]bool)
	for _, rel := range relationships {
		for _, roleAssign := range rel.Roles {
			if roleAssign.UserID == userObjID {
				roleMap[Role(roleAssign.Role)] = true
			}
		}
	}
	
	// Convert to slice
	var roles []Role
	for role := range roleMap {
		roles = append(roles, role)
	}
	
	return roles, nil
}

// ClearCache clears cached access results
func (c *RelationshipChecker) ClearCache(userIDs ...string) {
	c.cache.Clear(userIDs...)
}