package scalar

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarshalObjectID converts MongoDB ObjectID to GraphQL string
func MarshalObjectID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, `"`+id.Hex()+`"`)
	})
}

// UnmarshalObjectID converts GraphQL string to MongoDB ObjectID
func UnmarshalObjectID(v interface{}) (primitive.ObjectID, error) {
	str, ok := v.(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ObjectID must be a string")
	}

	objID, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid ObjectID format")
	}

	return objID, nil
}
