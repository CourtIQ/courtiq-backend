require('dotenv').config();
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const http = require('http');
const { json } = require('body-parser');
const { v4: uuidv4 } = require('uuid');

const config = require('./config');
const gateway = require('./gateway');
const logger = require('./logging/logger');
const { verifyIdToken, initializeFirebase } = require('./auth/firebase');
const { metrics, register } = require('./metrics/metrics');
const { waitForServices, checkAllServices } = require('./health/healthCheck');
const { retry } = require('./utils/retry');

// Track if server is shutting down
let isShuttingDown = false;

// Initialize services
async function initializeServices() {
  try {
    // Initialize Firebase Admin SDK
    initializeFirebase();
    
    // Apply startup delay if configured (useful in containerized environments)
    if (config.STARTUP_DELAY > 0) {
      logger.info(`Applying startup delay of ${config.STARTUP_DELAY}ms`);
      await new Promise(resolve => setTimeout(resolve, config.STARTUP_DELAY));
    }
    
    // Wait for downstream services to be healthy with extended retry
    await waitForServices(config.SERVICE_RETRY_MAX, config.SERVICE_RETRY_INTERVAL);
    
    return true;
  } catch (error) {
    logger.error('Failed to initialize services', { error: error.message });
    return false;
  }
}

// Set up a standalone HTTP server for health checks and metrics
function setupHttpServer() {
  // Create HTTP server for health checks and metrics
  const httpServer = http.createServer(async (req, res) => {
    // Basic routing for health and metrics endpoints
    if (req.url === config.HEALTH_CHECK_PATH) {
      if (isShuttingDown) {
        res.writeHead(503, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'shutting_down' }));
        return;
      }
      
      try {
        const servicesHealth = await checkAllServices();
        const allHealthy = Object.values(servicesHealth).every(status => status);
        
        res.writeHead(allHealthy ? 200 : 503, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          status: allHealthy ? 'healthy' : 'degraded',
          services: servicesHealth
        }));
      } catch (error) {
        logger.error('Health check failed', { error: error.message });
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'error' }));
      }
      return;
    }
    
    if (req.url === config.METRICS_PATH) {
      try {
        res.writeHead(200, { 'Content-Type': register.contentType });
        res.end(await register.metrics());
      } catch (error) {
        logger.error('Failed to get metrics', { error: error.message });
        res.writeHead(500);
        res.end();
      }
      return;
    }
    
    // For all other routes, return 404
    res.writeHead(404);
    res.end('Not Found');
  });

  // Start the HTTP server on a separate port
  const httpPort = config.PORT + 1;
  httpServer.listen(httpPort, () => {
    logger.info(`Health and metrics server ready at http://localhost:${httpPort}`);
    logger.info(`Health checks available at http://localhost:${httpPort}${config.HEALTH_CHECK_PATH}`);
    logger.info(`Metrics available at http://localhost:${httpPort}${config.METRICS_PATH}`);
  });

  return httpServer;
}

// Setup graceful shutdown
function setupGracefulShutdown(httpServer, apolloServer) {
  async function shutdown(signal) {
    logger.info(`Received ${signal}. Starting graceful shutdown...`);
    isShuttingDown = true;
    
    // Set a timeout to force shutdown if graceful shutdown takes too long
    const forceShutdownTimeout = setTimeout(() => {
      logger.error('Forced shutdown due to timeout');
      process.exit(1);
    }, 10000);
    
    try {
      // Shutdown Apollo Server
      await apolloServer.stop();
      logger.info('Apollo Server stopped successfully');
      
      // Close HTTP server - stop accepting new connections
      await new Promise((resolve) => {
        httpServer.close(resolve);
      });
      logger.info('HTTP server closed successfully');
      
      // Clear the force shutdown timeout
      clearTimeout(forceShutdownTimeout);
      
      // Exit gracefully
      logger.info('Server shut down successfully');
      process.exit(0);
    } catch (error) {
      logger.error('Error during server shutdown', { error: error.message });
      process.exit(1);
    }
  }
  
  // Listen for termination signals
  process.on('SIGTERM', () => shutdown('SIGTERM'));
  process.on('SIGINT', () => shutdown('SIGINT'));
  
  // Handle uncaught exceptions
  process.on('uncaughtException', (error) => {
    logger.error('Uncaught exception', { error: error.message, stack: error.stack });
    shutdown('uncaughtException');
  });
  
  // Handle unhandled promise rejections
  process.on('unhandledRejection', (reason, promise) => {
    logger.error('Unhandled promise rejection', {
      reason: reason instanceof Error ? reason.message : reason,
      stack: reason instanceof Error ? reason.stack : undefined
    });
  });
}

// Start the Apollo Server
async function startApolloServer() {
  try {
    // Initialize services first with retry
    const servicesInitialized = await retry(
      initializeServices,
      {
        maxRetries: 5,
        interval: 5000,
        shouldRetry: (err) => !isShuttingDown
      }
    );
    
    if (!servicesInitialized) {
      logger.error('Cannot start Apollo Server due to service initialization failure');
      process.exit(1);
    }
    
    // Set up HTTP server for health and metrics
    const httpServer = setupHttpServer();
    
    // Define plugins upfront
    const plugins = [{
      async serverWillStart() {
        logger.info('Apollo Server is starting up');
        return {
          async serverWillStop() {
            logger.info('Apollo Server is shutting down');
          }
        };
      },
      async requestDidStart() {
        const requestStartTime = process.hrtime();
        
        return {
          async didEncounterErrors({ errors }) {
            metrics.graphqlErrors.inc({ type: 'gateway' }, errors.length);
          },
          async willSendResponse({ request, response }) {
            const [seconds, nanoseconds] = process.hrtime(requestStartTime);
            const durationMs = (seconds * 1000) + (nanoseconds / 1000000);
            
            metrics.operationDuration.observe(
              { operation: request.operationName || 'unknown' },
              durationMs
            );
          }
        };
      }
    }];
    
    // Create Apollo Server with plugins defined upfront
    const server = new ApolloServer({
      gateway,
      ...config.APOLLO_OPTIONS,
      plugins
    });
    
    // IMPORTANT: Do not call server.start() here!
    // startStandaloneServer will handle starting the server
    
    // Now configure standalone server with context function
    const { url } = await startStandaloneServer(server, {
      context: async ({ req }) => {
        // Generate a unique request ID for tracing
        const requestId = uuidv4();
        
        // Extract and verify auth token
        const authHeader = req.headers.authorization || '';
        const token = authHeader.replace('Bearer ', '');
        let user = null;
        
        if (token) {
          try {
            user = await verifyIdToken(token);
            if (user) {
              logger.debug('User authenticated', {
                uid: user.uid,
                requestId
              });
            }
          } catch (error) {
            logger.warn('Failed to verify token', {
              error: error.message,
              requestId
            });
          }
        }
        
        return {
          user,
          requestId,
          startTime: Date.now() // For request duration tracking
        };
      },
      listen: { port: config.PORT }
    });
    
    logger.info(`🚀 Server ready at ${url}`);
    
    // Setup graceful shutdown
    setupGracefulShutdown(httpServer, server);
    
    return { server, httpServer };
  } catch (error) {
    logger.error('Failed to start Apollo Server', { error: error.message, stack: error.stack });
    throw error; // Re-throw to trigger retry
  }
}

// Start the server with retry logic for better fault tolerance
retry(
  startApolloServer,
  {
    maxRetries: 5,
    interval: 5000,
    shouldRetry: (err) => !isShuttingDown,
    operationName: 'Apollo Server startup'
  }
).catch(error => {
  logger.error('Failed to start server after multiple attempts', {
    error: error.message,
    stack: error.stack
  });
  process.exit(1);
});