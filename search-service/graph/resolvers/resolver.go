package resolvers

import "github.com/CourtIQ/courtiq-backend/search-service/internal/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	SearchService services.SearchService
}
