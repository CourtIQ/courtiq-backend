// Pseudocode
package scalar

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

// Let's define a type that can hold lat/lng:
type GeoPoint [2]float64 // [lng, lat]

func MarshalGeoPoint(g GeoPoint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// Convert e.g. [106.80, -6.22] -> "106.80,-6.22"
		s := fmt.Sprintf("%.6f,%.6f", g[0], g[1])
		// Wrap it in quotes for JSON
		fmt.Fprintf(w, "%q", s)
	})
}

func UnmarshalGeoPoint(v interface{}) (GeoPoint, error) {
	// Expect a string like "106.80,-6.22"
	s, ok := v.(string)
	if !ok {
		return GeoPoint{}, fmt.Errorf("GeoPoint must be a string")
	}
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return GeoPoint{}, fmt.Errorf("invalid GeoPoint format, want \"lng,lat\"")
	}
	lng, err := strconv.ParseFloat(parts[0], 64)
	if err != nil { /* handle */
	}
	lat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil { /* handle */
	}

	return GeoPoint{lng, lat}, nil
}
