package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sample struct for testing
type TestUser struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email,omitempty"`
	NoTag        string
	IgnoredField string             `bson:"-"`
	CustomName   string             `bson:"customFieldName"`
}

func TestBSONFieldName(t *testing.T) {
	userType := reflect.TypeOf(TestUser{})

	// Test basic field with bson tag
	name, err := BSONFieldName(userType, "Name")
	assert.NoError(t, err)
	assert.Equal(t, "name", name)

	// Test field with bson tag and options
	email, err := BSONFieldName(userType, "Email")
	assert.NoError(t, err)
	assert.Equal(t, "email", email)

	// Test field with no bson tag
	noTag, err := BSONFieldName(userType, "NoTag")
	assert.NoError(t, err)
	assert.Equal(t, "noTag", noTag)

	// Test field with ignored bson tag
	ignored, err := BSONFieldName(userType, "IgnoredField")
	assert.NoError(t, err)
	assert.Equal(t, "-", ignored)

	// Test field with custom name
	custom, err := BSONFieldName(userType, "CustomName")
	assert.NoError(t, err)
	assert.Equal(t, "customFieldName", custom)

	// Test non-existent field
	_, err = BSONFieldName(userType, "NonExistentField")
	assert.Error(t, err)

	// Test with pointer type
	ptrType := reflect.TypeOf(&TestUser{})
	idPtr, err := BSONFieldName(ptrType, "ID")
	assert.NoError(t, err)
	assert.Equal(t, "_id", idPtr)
}

func TestBSONFieldNameFromInstance(t *testing.T) {
	testUser := TestUser{
		ID:           primitive.NewObjectID(),
		Name:         "Test User",
		Email:        "test@example.com",
		NoTag:        "No tag value",
		IgnoredField: "Ignored",
		CustomName:   "Custom",
	}

	// Test getting field from instance
	name, err := BSONFieldNameFromInstance(testUser, "Name")
	assert.NoError(t, err)
	assert.Equal(t, "name", name)

	// Test with pointer instance
	ptrUser := &testUser
	namePtr, err := BSONFieldNameFromInstance(ptrUser, "Name")
	assert.NoError(t, err)
	assert.Equal(t, "name", namePtr)
}

func TestGetBSONFieldMap(t *testing.T) {
	userType := reflect.TypeOf(TestUser{})
	fieldMap := GetBSONFieldMap(userType)

	// Check expected mappings
	assert.Equal(t, "_id", fieldMap["ID"])
	assert.Equal(t, "name", fieldMap["Name"])
	assert.Equal(t, "email", fieldMap["Email"])
	assert.Equal(t, "customFieldName", fieldMap["CustomName"])
	
	// Check fields that shouldn't be in the map
	_, hasNoTag := fieldMap["NoTag"]
	assert.False(t, hasNoTag)
	
	_, hasIgnored := fieldMap["IgnoredField"]
	assert.False(t, hasIgnored)
}