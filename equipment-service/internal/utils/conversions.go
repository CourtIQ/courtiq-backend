// internal/utils/conversions.go
package utils

import "github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"

func ConvertStringTensionInput(input *model.StringTensionInput) *model.StringTension {
	if input == nil {
		return nil
	}
	return &model.StringTension{
		Mains:   input.Mains,
		Crosses: input.Crosses,
	}
}
