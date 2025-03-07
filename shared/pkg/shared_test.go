package pkg

import (
	"testing"

	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/configs"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/health"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/scalar"
	_ "github.com/CourtIQ/courtiq-backend/shared/pkg/validation"
)

// TestSharedPackage is a placeholder test that ensures all package tests are run
func TestSharedPackage(t *testing.T) {
	// This is a placeholder test that doesn't do anything
	// The actual tests are in the imported packages
}