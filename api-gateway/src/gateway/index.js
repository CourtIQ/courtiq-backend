const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { RemoteGraphQLDataSource } = require('@apollo/gateway');
const config = require('../config');
const getServiceUrl = require('../utils/getServiceUrl');
const logger = require('../logging/logger');
const { metrics } = require('../metrics/metrics');
const { retry } = require('../utils/retry');

// Custom DataSource that includes metrics, logging, and better error handling
class InstrumentedDataSource extends RemoteGraphQLDataSource {
  async process({ request, context }) {
    const startTime = process.hrtime();
    const serviceName = this.url.split('/')[2].split(':')[0]; // Extract service name from URL
    
    try {
      // Track active requests
      metrics.activeRequests.inc();
      
      const response = await retry(
        () => super.process({ request, context }),
        {
          maxRetries: 2,
          interval: 500,
          shouldRetry: (error) => {
            // Only retry certain types of errors (e.g., network errors)
            return error.name === 'FetchError' ||
                  error.name === 'AbortError' ||
                  (error.statusCode >= 500 && error.statusCode < 600);
          }
        }
      );
      
      // Record timing metrics
      const [seconds, nanoseconds] = process.hrtime(startTime);
      const durationMs = seconds * 1000 + nanoseconds / 1000000;
      metrics.microserviceResponseTime.observe({ service: serviceName }, durationMs);
      
      logger.debug(`Service ${serviceName} responded`, {
        durationMs,
        operationName: request.operationName,
        requestId: context.requestId
      });
      
      return response;
    } catch (error) {
      // Record error metrics
      metrics.microserviceErrors.inc({
        service: serviceName,
        error_type: error.name || 'UnknownError'
      });
      
      logger.error(`Error from ${serviceName}`, {
        error: error.message,
        operationName: request.operationName,
        requestId: context.requestId,
        errorName: error.name,
        statusCode: error.statusCode
      });
      
      throw error;
    } finally {
      // Decrement active requests counter
      metrics.activeRequests.dec();
    }
  }

  willSendRequest({ request, context }) {
    // Forward user claims to downstream services
    if (context?.user) {
      request.http.headers.set('X-User-Claims', JSON.stringify(context.user));
    }
    
    // Add request ID for tracing (if exists in context)
    if (context?.requestId) {
      request.http.headers.set('X-Request-ID', context.requestId);
    }
    
    // Track request metrics
    metrics.graphqlOperationsTotal.inc({
      operation: request.operationName || 'unknown',
      type: request.query?.definitions?.[0]?.operation || 'unknown'
    });
  }
}

// Create and validate subgraph configurations
const subgraphs = Object.values(config.SERVICES)
  .filter(service => {
    // Only include required services or services with URLs
    const isValid = service.required || service.url;
    if (!isValid) {
      logger.warn(`Skipping optional service ${service.name} - no URL provided`);
    }
    return isValid;
  })
  .map(service => {
    const url = `${getServiceUrl(service.name, service.url)}/graphql`;
    logger.info(`Registering subgraph: ${service.name}`, { url });
    
    return {
      name: service.name,
      url,
    };
  });

if (subgraphs.length === 0) {
  logger.error('No valid subgraphs found. Check service configuration.');
  throw new Error('No valid subgraphs configured');
}

// Create the gateway with instrumentation
const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs,
    pollIntervalInMs: 20000, // Poll every 20 seconds
  }),
  
  buildService({ name, url }) {
    logger.debug(`Building service: ${name}`, { url });
    return new InstrumentedDataSource({
      url,
    });
  },
  
  // Add experimental didFailComposition handler for better error reporting
  experimental_didFailComposition: (error) => {
    logger.error('Failed to compose supergraph schema', {
      error: error.message,
      errors: error.errors?.map(e => e.message)
    });
  },
  
  // Add experimental didUpdateComposition handler for logging
  experimental_didUpdateComposition: ({ supergraphSdl }) => {
    logger.info('Supergraph schema updated successfully', {
      schemaLength: supergraphSdl.length,
    });
  },
});

module.exports = gateway;