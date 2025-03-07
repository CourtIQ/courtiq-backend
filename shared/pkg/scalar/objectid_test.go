package scalar

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMarshalObjectID(t *testing.T) {
	// Create an ObjectID
	id := primitive.NewObjectID()

	// Marshal it
	marshaler := MarshalObjectID(id)

	// Convert to string
	builder := &strings.Builder{}
	marshaler.MarshalGQL(builder)

	// Verify the result
	expected := `"` + id.Hex() + `"`
	assert.Equal(t, expected, builder.String())
}

func TestUnmarshalObjectID(t *testing.T) {
	// Create an ObjectID for testing
	id := primitive.NewObjectID()

	// Test valid case
	result, err := UnmarshalObjectID(id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, id, result)

	// Test invalid type case
	result, err = UnmarshalObjectID(123)
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, result)

	// Test invalid format case
	result, err = UnmarshalObjectID("not-an-objectid")
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, result)
}
