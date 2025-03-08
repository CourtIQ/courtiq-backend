const winston = require('winston');
const config = require('../config');

// Define log levels
const levels = {
  error: 0,
  warn: 1,
  info: 2,
  http: 3,
  debug: 4,
};

// Define colors for each level
const colors = {
  error: 'red',
  warn: 'yellow',
  info: 'green',
  http: 'magenta',
  debug: 'blue',
};

// Tell winston about our colors
winston.addColors(colors);

// Create a custom format for development that's more readable
const developmentFormat = winston.format.combine(
  winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss:ms' }),
  winston.format.colorize({ all: true }),
  winston.format.printf(({ level, message, timestamp, ...metadata }) => {
    // Extract the most important metadata
    const requestId = metadata.requestId ? `[${metadata.requestId}]` : '';
    const userId = metadata.uid || (metadata.user && metadata.user.uid) || '';
    const userInfo = userId ? `[User: ${userId}]` : '';
    
    // Format other metadata
    const meta = Object.keys(metadata).length > 0
      ? `\n${JSON.stringify(metadata, null, 2)}`
      : '';
    
    return `${timestamp} ${level}: ${requestId}${userInfo} ${message}${meta}`;
  })
);

// Create a production format that's more structured for machine parsing
const productionFormat = winston.format.combine(
  winston.format.timestamp(),
  winston.format.errors({ stack: true }),
  winston.format.json()
);

// Choose format based on environment
const format = config.ENV === 'development' ? developmentFormat : productionFormat;

// Create the logger with console and file transports
const logger = winston.createLogger({
  level: config.LOG_LEVEL || 'info',
  levels,
  format,
  defaultMeta: { service: config.SERVICE_NAME },
  transports: [
    new winston.transports.Console(),
    // Add file transport in production
    ...(config.ENV === 'production'
      ? [
          new winston.transports.File({
            filename: 'logs/error.log',
            level: 'error',
            maxsize: 5242880, // 5MB
            maxFiles: 5,
          }),
          new winston.transports.File({
            filename: 'logs/combined.log',
            maxsize: 5242880, // 5MB
            maxFiles: 5,
          }),
        ]
      : []),
  ],
  exceptionHandlers: [
    new winston.transports.Console(),
    ...(config.ENV === 'production'
      ? [new winston.transports.File({ filename: 'logs/exceptions.log' })]
      : []),
  ],
  rejectionHandlers: [
    new winston.transports.Console(),
    ...(config.ENV === 'production'
      ? [new winston.transports.File({ filename: 'logs/rejections.log' })]
      : []),
  ],
});

// Export a standardized logger interface
module.exports = {
  error: (message, meta = {}) => logger.error(message, meta),
  warn: (message, meta = {}) => logger.warn(message, meta),
  info: (message, meta = {}) => logger.info(message, meta),
  http: (message, meta = {}) => logger.http(message, meta),
  debug: (message, meta = {}) => logger.debug(message, meta),
  
  // Create a child logger with request context
  child: (requestContext) => {
    return {
      error: (message, meta = {}) => logger.error(message, { ...requestContext, ...meta }),
      warn: (message, meta = {}) => logger.warn(message, { ...requestContext, ...meta }),
      info: (message, meta = {}) => logger.info(message, { ...requestContext, ...meta }),
      http: (message, meta = {}) => logger.http(message, { ...requestContext, ...meta }),
      debug: (message, meta = {}) => logger.debug(message, { ...requestContext, ...meta }),
    };
  },
  
  // Log and return promises for async functions
  promise: (promise, message) => {
    return promise
      .then((result) => {
        logger.info(`${message}: success`);
        return result;
      })
      .catch((error) => {
        logger.error(`${message}: failed`, { error: error.message });
        throw error;
      });
  },
};