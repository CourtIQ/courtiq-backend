package scalar

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalGeoPoint(t *testing.T) {
	// Create a GeoPoint
	gp := GeoPoint{106.80, -6.22}
	
	// Marshal it
	marshaler := MarshalGeoPoint(gp)
	
	// Convert to string
	builder := &strings.Builder{}
	marshaler.MarshalGQL(builder)
	
	// Verify the result
	assert.Contains(t, builder.String(), "106.800000")
	assert.Contains(t, builder.String(), "-6.220000")
}

func TestUnmarshalGeoPoint(t *testing.T) {
	// Test valid case
	result, err := UnmarshalGeoPoint("106.80,-6.22")
	assert.NoError(t, err)
	assert.InDelta(t, 106.80, result[0], 0.0001)
	assert.InDelta(t, -6.22, result[1], 0.0001)
	
	// Test invalid type case
	result, err = UnmarshalGeoPoint(123)
	assert.Error(t, err)
	
	// Test invalid format case
	result, err = UnmarshalGeoPoint("not-a-geopoint")
	assert.Error(t, err)
	
	// Test too many parts
	result, err = UnmarshalGeoPoint("106.80,-6.22,10.0")
	assert.Error(t, err)
	
	// Test too few parts
	result, err = UnmarshalGeoPoint("106.80")
	assert.Error(t, err)
}