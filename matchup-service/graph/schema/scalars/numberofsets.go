package scalars

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/99designs/gqlgen/graphql"
)

// NumberOfSets is an int, but only allowed to be 1, 3, 5.
//
// This custom scalar ensures the caller can only provide integers
// that match these discrete values, preventing invalid "number of sets".
type NumberOfSets int

// validSets enumerates all allowable integer values for NumberOfSets.
var validSets = map[int]bool{
	1: true, // e.g. "ONE"
	3: true, // "THREE"
	5: true, // "FIVE"
}

// MarshalNumberOfSets converts our custom scalar to a GraphQL integer
// for the outgoing response.
func MarshalNumberOfSets(n NumberOfSets) graphql.Marshaler {
	return graphql.MarshalInt(int(n))
}

// UnmarshalNumberOfSets parses the raw input into NumberOfSets,
// restricting it to integer values 1, 3, or 5 only.
func UnmarshalNumberOfSets(v interface{}) (NumberOfSets, error) {
	floatVal, err := coerceToFloat64(v)
	if err != nil {
		return 0, fmt.Errorf("NumberOfSets must be numeric: %w", err)
	}

	// Reject fractional values like 2.5, 3.14, etc.
	if floatVal != math.Trunc(floatVal) {
		return 0, fmt.Errorf("NumberOfSets must be an integer, got fractional value %.2f", floatVal)
	}

	intVal := int(floatVal)
	if !validSets[intVal] {
		return 0, fmt.Errorf("NumberOfSets must be one of [1, 3, 5], got %d", intVal)
	}
	return NumberOfSets(intVal), nil
}

// coerceToFloat64 attempts to convert an interface{} to a float64.
// It handles int, int64, float64, and json.Number so we gracefully
// parse any integer-like inputs that different JSON decoders might produce.
//
// We put this helper here so we don't repeat logic across scalars.
func coerceToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case json.Number:
		f, err := val.Float64()
		if err != nil {
			return 0, fmt.Errorf("could not parse json.Number: %v", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("got %T (not a numeric type)", v)
	}
}
