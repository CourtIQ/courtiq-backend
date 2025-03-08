const logger = require('../logging/logger');

/**
 * Retries a function until it succeeds or max retries is reached
 * @param {Function} fn - The function to retry
 * @param {Object} options - Retry options
 * @param {number} [options.maxRetries=3] - Maximum number of retries
 * @param {number} [options.interval=1000] - Interval between retries in ms
 * @param {Function} [options.shouldRetry] - Function to determine if retry should happen
 * @param {string} [options.operationName] - Name of the operation for logging
 * @returns {Promise<*>} The result of the function
 */
async function retry(fn, {
  maxRetries = 3,
  interval = 1000,
  shouldRetry = () => true,
  operationName = 'operation'
} = {}) {
  let retries = 0;
  let lastError = null;
  
  while (retries <= maxRetries) {
    try {
      // If this is a retry, log at debug level
      if (retries > 0) {
        logger.debug(`Retry attempt ${retries}/${maxRetries} for ${operationName}`);
      }
      
      return await fn();
    } catch (error) {
      lastError = error;
      retries++;
      
      // Check if we should retry
      const shouldTryAgain = retries <= maxRetries && shouldRetry(error);
      
      if (!shouldTryAgain) {
        // If it's the last retry or shouldRetry returned false, break out
        break;
      }
      
      logger.warn(`${operationName} failed, retrying...`, {
        error: error.message,
        retry: retries,
        maxRetries,
        nextRetryMs: interval,
      });
      
      // Wait before the next retry
      await new Promise((resolve) => setTimeout(resolve, interval));
    }
  }
  
  // If we got here, all retries failed
  logger.error(`${operationName} failed after ${retries} retries`, {
    error: lastError?.message,
    stack: lastError?.stack,
  });
  
  throw lastError;
}

/**
 * Utility to create a function that retries a specific operation
 * @param {Function} fn - The function to retry
 * @param {Object} options - Retry options as in retry()
 * @returns {Function} A function that will retry the operation
 */
function withRetry(fn, options = {}) {
  return async (...args) => {
    return retry(() => fn(...args), options);
  };
}

module.exports = {
  retry,
  withRetry
};