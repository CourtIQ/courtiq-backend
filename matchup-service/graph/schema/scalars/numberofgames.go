package scalars

import (
	"fmt"
	"math"

	"github.com/99designs/gqlgen/graphql"
)

// NumberOfGames is an int with constraints: 1, 3, 4, 5, 6, 10.
//
// This custom scalar ensures the caller can only provide one of these
// discrete integer values, preventing invalid "number of games" from
// sneaking into the system.
type NumberOfGames int

// validValues enumerates all permissible integer values for NumberOfGames.
var validValues = map[int]bool{
	1:  true, // Could map to "ONE"
	3:  true, // "THREE"
	4:  true, // "FOUR"
	5:  true, // "FIVE"
	6:  true, // "SIX"
	10: true, // "TEN"
}

// MarshalNumberOfGames converts our custom scalar to a GraphQL integer
// for the outgoing response. This is part of the gqlgen custom scalar interface.
func MarshalNumberOfGames(n NumberOfGames) graphql.Marshaler {
	return graphql.MarshalInt(int(n))
}

// UnmarshalNumberOfGames parses incoming data into NumberOfGames,
// ensuring it's a whole integer within our valid set {1,3,4,5,6,10}.
func UnmarshalNumberOfGames(v interface{}) (NumberOfGames, error) {
	floatVal, err := coerceToFloat64(v)
	if err != nil {
		// Provide a clear error message for devs/clients
		return 0, fmt.Errorf("NumberOfGames must be numeric: %w", err)
	}

	// Reject fractional numbers (e.g. 4.5).
	if floatVal != math.Trunc(floatVal) {
		return 0, fmt.Errorf("NumberOfGames must be an integer, got fractional value %.2f", floatVal)
	}

	intVal := int(floatVal)
	if !validValues[intVal] {
		return 0, fmt.Errorf("NumberOfGames must be one of [1,3,4,5,6,10], got %d", intVal)
	}
	return NumberOfGames(intVal), nil
}
