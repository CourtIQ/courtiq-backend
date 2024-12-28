package resolvers

import "github.com/CourtIQ/courtiq-backend/matchup-service/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MatchupService service.MatchUpServiceInterface // Note: changed field name to be more idiomatic
}
