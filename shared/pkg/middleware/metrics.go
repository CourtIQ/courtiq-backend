package middleware

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	gqlRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_requests_total",
			Help: "Total number of GraphQL requests",
		},
		[]string{"operation", "service"},
	)

	gqlRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "graphql_request_duration_seconds",
			Help:    "Duration of GraphQL requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "service"},
	)

	gqlErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_errors_total",
			Help: "Total number of GraphQL errors",
		},
		[]string{"operation", "service"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(gqlRequestsTotal)
	prometheus.MustRegister(gqlRequestDuration)
	prometheus.MustRegister(gqlErrorsTotal)
}

// MetricsExtension implements graphql.HandlerExtension to collect metrics
type MetricsExtension struct {
	serviceName string
}

// NewMetricsExtension creates a new MetricsExtension
func NewMetricsExtension(serviceName string) *MetricsExtension {
	return &MetricsExtension{
		serviceName: serviceName,
	}
}

var _ graphql.HandlerExtension = &MetricsExtension{}
var _ graphql.ResponseInterceptor = &MetricsExtension{}
var _ graphql.FieldInterceptor = &MetricsExtension{}

// ExtensionName returns the extension name
func (m *MetricsExtension) ExtensionName() string {
	return "Metrics"
}

// Validate validates the extension
func (m *MetricsExtension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse records metrics about the GraphQL response
func (m *MetricsExtension) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	// Get operation name
	operationContext := graphql.GetOperationContext(ctx)
	operationName := operationContext.OperationName
	if operationName == "" {
		operationName = "anonymous"
	}

	// Increment request counter
	gqlRequestsTotal.WithLabelValues(operationName, m.serviceName).Inc()

	// Start timer
	start := time.Now()

	// Execute the next handler
	resp := next(ctx)

	// Record request duration
	duration := time.Since(start).Seconds()
	gqlRequestDuration.WithLabelValues(operationName, m.serviceName).Observe(duration)

	// Check for errors
	if len(resp.Errors) > 0 {
		gqlErrorsTotal.WithLabelValues(operationName, m.serviceName).Add(float64(len(resp.Errors)))
	}

	return resp
}

// InterceptField implements field-level metrics if needed
func (m *MetricsExtension) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	return next(ctx)
}

func GetMetricsConfig(serviceName string) graphql.HandlerExtension {
	return NewMetricsExtension(serviceName)
}
