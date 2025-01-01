package scalars

import (
	"fmt"
	"math"

	"github.com/99designs/gqlgen/graphql"
)

// TiebreakPoints is an integer scalar restricted to {5, 6, 7, 8, 9, 10}.
type TiebreakPoints int

// validTiebreakPoints enumerates all permissible integer values.
var validTiebreakPoints = map[int]bool{
	5:  true,
	6:  true,
	7:  true,
	8:  true,
	9:  true,
	10: true,
}

// MarshalTiebreakPoints converts our custom scalar to an integer
// for the GraphQL response.
func MarshalTiebreakPoints(tp TiebreakPoints) graphql.Marshaler {
	return graphql.MarshalInt(int(tp))
}

// UnmarshalTiebreakPoints ensures the user-provided value is a numeric integer
// and one of [5, 6, 7, 8, 9, 10]. Rejects floats (e.g., 7.5) or out-of-range values.
func UnmarshalTiebreakPoints(v interface{}) (TiebreakPoints, error) {
	floatVal, err := coerceToFloat64(v)
	if err != nil {
		return 0, fmt.Errorf("TiebreakPoints must be numeric: %w", err)
	}

	// Reject fractional numbers, e.g., 7.25
	if floatVal != math.Trunc(floatVal) {
		return 0, fmt.Errorf("TiebreakPoints must be an integer, got fractional value %.2f", floatVal)
	}

	intVal := int(floatVal)
	if !validTiebreakPoints[intVal] {
		return 0, fmt.Errorf("TiebreakPoints must be one of [5,6,7,8,9,10], got %d", intVal)
	}

	return TiebreakPoints(intVal), nil
}
