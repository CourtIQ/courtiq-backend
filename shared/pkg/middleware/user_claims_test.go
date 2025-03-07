package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetMongoID(t *testing.T) {
	// Create a MongoDB ObjectID
	objID := primitive.NewObjectID()
	
	// Test case 1: Valid user claims
	claims := UserClaims{
		MongoID: objID.Hex(),
		Email:   "test@example.com",
	}
	
	claimsJSON, err := json.Marshal(claims)
	require.NoError(t, err)
	
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-Claims", string(claimsJSON))
	
	mongoID, err := GetMongoID(req)
	assert.NoError(t, err)
	assert.Equal(t, objID.Hex(), mongoID)
	
	// Test case 2: Missing header
	req = httptest.NewRequest("GET", "/", nil)
	mongoID, err = GetMongoID(req)
	assert.Error(t, err)
	assert.Empty(t, mongoID)
	
	// Test case 3: Invalid JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-Claims", "invalid json")
	mongoID, err = GetMongoID(req)
	assert.Error(t, err)
	assert.Empty(t, mongoID)
	
	// Test case 4: Missing mongoId in claims
	claims = UserClaims{
		Email: "test@example.com",
		// No MongoID
	}
	
	claimsJSON, err = json.Marshal(claims)
	require.NoError(t, err)
	
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-Claims", string(claimsJSON))
	
	mongoID, err = GetMongoID(req)
	assert.Error(t, err)
	assert.Empty(t, mongoID)
}

func TestWithUserClaims(t *testing.T) {
	// Create a MongoDB ObjectID
	objID := primitive.NewObjectID()
	
	// Create a test handler that checks the mongoId in context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mongoID, err := GetMongoIDFromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mongoID.Hex()))
	})
	
	// Wrap the test handler with our middleware
	handler := WithUserClaims(testHandler)
	
	// Test case 1: Valid claims
	claims := UserClaims{
		MongoID: objID.Hex(),
		Email:   "test@example.com",
	}
	
	claimsJSON, err := json.Marshal(claims)
	require.NoError(t, err)
	
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-User-Claims", string(claimsJSON))
	
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, objID.Hex(), recorder.Body.String())
	
	// Test case 2: Missing claims (should still call handler, but with no mongoId)
	req = httptest.NewRequest("GET", "/", nil)
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetMongoIDFromContext(t *testing.T) {
	// Test case 1: Valid mongoId in context
	objID := primitive.NewObjectID()
	ctx := context.WithValue(context.Background(), mongoIDKey, objID.Hex())
	
	mongoID, err := GetMongoIDFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, objID, mongoID)
	
	// Test case 2: No mongoId in context
	ctx = context.Background()
	mongoID, err = GetMongoIDFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, mongoID)
	
	// Test case 3: Invalid ObjectID in context
	ctx = context.WithValue(context.Background(), mongoIDKey, "not-an-objectid")
	mongoID, err = GetMongoIDFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, mongoID)
	
	// Test case 4: Wrong type in context
	ctx = context.WithValue(context.Background(), mongoIDKey, 123)
	mongoID, err = GetMongoIDFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, primitive.NilObjectID, mongoID)
}