package middleware

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsExtension(t *testing.T) {
	ext := NewMetricsExtension("test-service")
	assert.NotNil(t, ext)
	assert.Equal(t, "test-service", ext.serviceName)
}

func TestExtensionName(t *testing.T) {
	ext := NewMetricsExtension("test-service")
	assert.Equal(t, "Metrics", ext.ExtensionName())
}

func TestValidate(t *testing.T) {
	ext := NewMetricsExtension("test-service")
	// This should not error
	err := ext.Validate(nil)
	assert.NoError(t, err)
}

// Create a mock response handler for testing InterceptResponse
type mockResponseHandler struct {
	called bool
	resp   *graphql.Response
}

func (m *mockResponseHandler) Handle(ctx context.Context) *graphql.Response {
	m.called = true
	return m.resp
}

func TestInterceptResponse(t *testing.T) {
	ext := NewMetricsExtension("test-service")
	
	// Create a mock response
	mockResp := &graphql.Response{
		Data: []byte(`{"test":"data"}`),
	}
	
	// Create a mock response handler
	mockHandler := &mockResponseHandler{
		resp: mockResp,
	}
	
	// Create a context with operation info
	opCtx := &graphql.OperationContext{
		OperationName: "TestQuery",
	}
	ctx := graphql.WithOperationContext(context.Background(), opCtx)
	
	// Call the interceptor
	resp := ext.InterceptResponse(ctx, mockHandler.Handle)
	
	// Verify the handler was called
	assert.True(t, mockHandler.called)
	
	// Verify the response was returned
	assert.Equal(t, mockResp, resp)
}

// Mock resolver for testing InterceptField
type mockResolver struct {
	called bool
	result interface{}
	err    error
}

func (m *mockResolver) Resolve(ctx context.Context) (interface{}, error) {
	m.called = true
	return m.result, m.err
}

func TestInterceptField(t *testing.T) {
	ext := NewMetricsExtension("test-service")
	
	// Create a mock resolver
	expected := "test-result"
	resolver := &mockResolver{
		result: expected,
		err:    nil,
	}
	
	// Call the interceptor
	result, err := ext.InterceptField(context.Background(), resolver.Resolve)
	
	// Verify the resolver was called
	assert.True(t, resolver.called)
	
	// Verify the result was returned
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestGetMetricsConfig(t *testing.T) {
	ext := GetMetricsConfig("test-service")
	assert.NotNil(t, ext)
	assert.IsType(t, &MetricsExtension{}, ext)
}