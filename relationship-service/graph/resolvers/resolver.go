package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/services"
)


type Resolver struct{
	RelationshipService services.RelationshipService
}
