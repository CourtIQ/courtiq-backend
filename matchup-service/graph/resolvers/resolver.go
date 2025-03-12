package resolvers

import "github.com/CourtIQ/courtiq-backend/matchup-service/internal/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MatchUpServiceInterface services.MatchUpServiceIntf // Note: changed field name to be more idiomatic
}
