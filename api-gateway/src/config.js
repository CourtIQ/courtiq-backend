// Environment and application defaults
const ENV = process.env.NODE_ENV || 'development';
const SERVICE_NAME = process.env.SERVICE_NAME || 'api-gateway';
const PORT = parseInt(process.env.PORT, 10) || 3000;
const GRAPHQL_PLAYGROUND = process.env.GRAPHQL_PLAYGROUND === 'true';
const LOG_LEVEL = process.env.LOG_LEVEL || (ENV === 'development' ? 'debug' : 'info');

// Health check configuration
const HEALTH_CHECK_ENABLED = process.env.HEALTH_CHECK_ENABLED !== 'false';
const HEALTH_CHECK_PATH = process.env.HEALTH_CHECK_PATH || '/health';
const METRICS_PATH = process.env.METRICS_PATH || '/metrics';

// Startup and retry configuration
const STARTUP_DELAY = parseInt(process.env.STARTUP_DELAY, 10) || 0;
const SERVICE_RETRY_MAX = parseInt(process.env.SERVICE_RETRY_MAX, 10) || 20;
const SERVICE_RETRY_INTERVAL = parseInt(process.env.SERVICE_RETRY_INTERVAL, 10) || 3000;
const GATEWAY_POLL_INTERVAL = parseInt(process.env.GATEWAY_POLL_INTERVAL, 10) || 30000;

// Security and timeout configuration
const TOKEN_VALIDATION_CACHE_TTL = parseInt(process.env.TOKEN_VALIDATION_CACHE_TTL, 10) || 300000; // 5 minutes
const REQUEST_TIMEOUT = parseInt(process.env.REQUEST_TIMEOUT, 10) || 30000; // 30 seconds

// Database connection URL
const MONGODB_URL = process.env.MONGODB_URL;

// Firebase configuration
const FIREBASE_CONFIG = process.env.FIREBASE_CONFIG ?
  JSON.parse(process.env.FIREBASE_CONFIG) : {};
  
const FIREBASE_SERVICE_ACCOUNT = process.env.FIREBASE_SERVICE_ACCOUNT ?
  JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT) : {};

// Apollo Server options
const APOLLO_OPTIONS = {
  introspection: GRAPHQL_PLAYGROUND,
  csrfPrevention: true,
};

// Microservice endpoints with validation
const SERVICES = {
  userService: {
    name: 'user-service',
    url: process.env.USER_SERVICE_URL || '',
    required: true,
  },
  relationshipService: {
    name: 'relationship-service',
    url: process.env.RELATIONSHIP_SERVICE_URL || '',
    required: true,
  },
  matchupService: {
    name: 'matchup-service',
    url: process.env.MATCHUP_SERVICE_URL || '',
    required: true,
  },
  equipmentService: {
    name: 'equipment-service',
    url: process.env.EQUIPMENT_SERVICE_URL || '',
    required: true,
  },
  searchService: {
    name: 'search-service',
    url: process.env.SEARCH_SERVICE_URL || '',
    required: true,
  },
};

// Validate configuration
function validateConfig() {
  const missingRequiredServices = Object.values(SERVICES)
    .filter(service => service.required && !service.url)
    .map(service => service.name);
  
  if (missingRequiredServices.length > 0) {
    console.warn(`⚠️ Missing URLs for required services: ${missingRequiredServices.join(', ')}`);
    console.warn('The gateway may fail to start if these services cannot be discovered.');
  }
  
  if (!FIREBASE_SERVICE_ACCOUNT || Object.keys(FIREBASE_SERVICE_ACCOUNT).length === 0) {
    console.warn('⚠️ Firebase service account not provided. Authentication will not work.');
  }
}

// Run validation
validateConfig();

module.exports = {
  ENV,
  SERVICE_NAME,
  PORT,
  GRAPHQL_PLAYGROUND,
  LOG_LEVEL,
  HEALTH_CHECK_ENABLED,
  HEALTH_CHECK_PATH,
  METRICS_PATH,
  STARTUP_DELAY,
  SERVICE_RETRY_MAX,
  SERVICE_RETRY_INTERVAL,
  GATEWAY_POLL_INTERVAL,
  TOKEN_VALIDATION_CACHE_TTL,
  REQUEST_TIMEOUT,
  MONGODB_URL,
  FIREBASE_CONFIG,
  FIREBASE_SERVICE_ACCOUNT,
  APOLLO_OPTIONS,
  SERVICES,
};