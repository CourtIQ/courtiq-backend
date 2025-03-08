const promClient = require('prom-client');
const logger = require('../logging/logger');

// Create a Registry to register our metrics
const register = new promClient.Registry();

// Add default metrics (CPU, memory, etc.)
promClient.collectDefaultMetrics({ register });

// Define custom metrics
const httpRequestDurationMicroseconds = new promClient.Histogram({
  name: 'http_request_duration_ms',
  help: 'Duration of HTTP requests in ms',
  labelNames: ['route', 'method', 'status'],
  buckets: [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]
});

const httpRequestsTotal = new promClient.Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['route', 'method', 'status']
});

const graphqlOperationsTotal = new promClient.Counter({
  name: 'graphql_operations_total',
  help: 'Total number of GraphQL operations',
  labelNames: ['operation', 'type']
});

const operationDuration = new promClient.Histogram({
  name: 'graphql_operation_duration_ms',
  help: 'Duration of GraphQL operations in ms',
  labelNames: ['operation'],
  buckets: [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]
});

const microserviceResponseTime = new promClient.Histogram({
  name: 'microservice_response_time_ms',
  help: 'Response time of microservices in ms',
  labelNames: ['service'],
  buckets: [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]
});

const microserviceErrors = new promClient.Counter({
  name: 'microservice_errors_total',
  help: 'Total number of microservice errors',
  labelNames: ['service', 'error_type']
});

const graphqlErrors = new promClient.Counter({
  name: 'graphql_errors_total',
  help: 'Total number of GraphQL errors',
  labelNames: ['type']
});

const activeRequests = new promClient.Gauge({
  name: 'active_requests',
  help: 'Number of active requests being processed'
});

// Register all metrics
register.registerMetric(httpRequestDurationMicroseconds);
register.registerMetric(httpRequestsTotal);
register.registerMetric(graphqlOperationsTotal);
register.registerMetric(operationDuration);
register.registerMetric(microserviceResponseTime);
register.registerMetric(microserviceErrors);
register.registerMetric(graphqlErrors);
register.registerMetric(activeRequests);

logger.info('Metrics initialized');

module.exports = {
  register,
  metrics: {
    httpRequestDurationMicroseconds,
    httpRequestsTotal,
    graphqlOperationsTotal,
    operationDuration,
    microserviceResponseTime,
    microserviceErrors,
    graphqlErrors,
    activeRequests
  },
  middleware: {
    // Apollo Plugin for GraphQL metrics
    metricsPlugin: {
      async requestDidStart(requestContext) {
        const startTime = process.hrtime();
        activeRequests.inc();
        
        return {
          async didEncounterErrors({ errors }) {
            microserviceErrors.inc({ service: 'gateway', error_type: 'graphql' }, errors.length);
            graphqlErrors.inc({ type: 'operation' }, errors.length);
          },
          async willSendResponse({ operationName, operation }) {
            activeRequests.dec();
            if (operationName) {
              const operationType = operation?.operation || 'unknown';
              const [seconds, nanoseconds] = process.hrtime(startTime);
              const durationMs = (seconds * 1000) + (nanoseconds / 1000000);
              
              graphqlOperationsTotal.inc({ operation: operationName, type: operationType });
              operationDuration.observe({ operation: operationName }, durationMs);
            }
          }
        };
      }
    }
  }
};