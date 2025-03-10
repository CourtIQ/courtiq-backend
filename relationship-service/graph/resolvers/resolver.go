package resolvers

import (
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// RelationshipService is the service interface for relationship management
	RelationshipService services.RelationshipServiceIntf
}