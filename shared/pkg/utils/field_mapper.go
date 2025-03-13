package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// BSONFieldName returns the BSON field name for a given struct field
// Usage:
//   - structType: Pass in a struct type using reflect.TypeOf(YourStruct{})
//   - fieldName: The name of the field as defined in the struct
//
// Returns:
//   - The bson tag value for the field, or the field name if no tag is present
//   - Empty string and error if the field doesn't exist
func BSONFieldName(structType reflect.Type, fieldName string) (string, error) {
	// Check if we're dealing with a pointer and get the underlying type
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	// Make sure we're dealing with a struct
	if structType.Kind() != reflect.Struct {
		return "", fmt.Errorf("type must be a struct, got %s", structType.Kind())
	}

	// Try to find the field
	field, found := structType.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field '%s' not found in struct %s", fieldName, structType.Name())
	}

	// Get the bson tag
	bsonTag := field.Tag.Get("bson")
	
	// If no tag is present, return the field name with first letter lowercase
	if bsonTag == "" {
		return strings.ToLower(fieldName[:1]) + fieldName[1:], nil
	}

	// If bson tag has options (like omitempty), extract just the field name
	parts := strings.Split(bsonTag, ",")
	return parts[0], nil
}

// BSONFieldNameFromInstance is a convenience function that gets the BSON field name
// from an instance of a struct rather than requiring the caller to use reflect.TypeOf
func BSONFieldNameFromInstance(structInstance interface{}, fieldName string) (string, error) {
	structType := reflect.TypeOf(structInstance)
	return BSONFieldName(structType, fieldName)
}

// GetBSONFieldMap returns a map of struct field names to BSON field names for a given struct type
func GetBSONFieldMap(structType reflect.Type) map[string]string {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]string)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		bsonTag := field.Tag.Get("bson")
		
		// Skip fields without bson tags or with "-" as the tag
		if bsonTag == "" || bsonTag == "-" {
			continue
		}
		
		// Handle bson tags with options
		parts := strings.Split(bsonTag, ",")
		result[field.Name] = parts[0]
	}
	
	return result
}