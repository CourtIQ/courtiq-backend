package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	timeout := 5 * time.Second
	handler := NewHandler(timeout)
	
	assert.NotNil(t, handler)
	assert.Equal(t, timeout, handler.timeout)
}

func TestAddCheck(t *testing.T) {
	handler := NewHandler(time.Second)
	
	// Add a check
	handler.AddCheck("test", func(ctx context.Context) (Status, error) {
		return StatusUp, nil
	})
	
	// Verify the check was added
	assert.Len(t, handler.checks, 1)
	assert.Contains(t, handler.checks, "test")
}

func TestSetReady(t *testing.T) {
	handler := NewHandler(time.Second)
	assert.False(t, handler.readiness)
	
	handler.SetReady(true)
	assert.True(t, handler.readiness)
	
	handler.SetReady(false)
	assert.False(t, handler.readiness)
}

func TestHandleLiveness(t *testing.T) {
	handler := NewHandler(time.Second)
	
	// Create a request to test the handler
	req, err := http.NewRequest("GET", "/health/live", nil)
	require.NoError(t, err)
	
	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()
	
	// Call the handler
	handler.HandleLiveness(recorder, req)
	
	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	// Check the response body
	var response HealthResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, StatusUp, response.Status)
	assert.Empty(t, response.Components)
}

func TestHandleReadiness(t *testing.T) {
	handler := NewHandler(time.Second)
	
	// Create a request
	req, err := http.NewRequest("GET", "/health/ready", nil)
	require.NoError(t, err)
	
	// Test 1: Not ready
	recorder := httptest.NewRecorder()
	handler.HandleReadiness(recorder, req)
	
	// Should return 503 Service Unavailable when not ready
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	
	// Test 2: Ready with passing checks
	handler.SetReady(true)
	handler.AddCheck("test1", func(ctx context.Context) (Status, error) {
		return StatusUp, nil
	})
	
	recorder = httptest.NewRecorder()
	handler.HandleReadiness(recorder, req)
	
	// Should return 200 OK
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response HealthResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, StatusUp, response.Status)
	assert.Len(t, response.Components, 1)
	assert.Equal(t, "test1", response.Components[0].Name)
	assert.Equal(t, StatusUp, response.Components[0].Status)
	
	// Test 3: Ready with failing checks
	handler.AddCheck("test2", func(ctx context.Context) (Status, error) {
		return StatusDown, errors.New("service down")
	})
	
	recorder = httptest.NewRecorder()
	handler.HandleReadiness(recorder, req)
	
	// Should return 503 Service Unavailable
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, StatusDown, response.Status)
	assert.Len(t, response.Components, 2)
	assert.Equal(t, "service down", response.Components[1].Error)
}