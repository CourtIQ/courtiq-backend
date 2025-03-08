const axios = require('axios');
const logger = require('../logging/logger');
const config = require('../config');
const getServiceUrl = require('../utils/getServiceUrl');
const { retry } = require('../utils/retry');

/**
 * Checks health of a service with retry logic
 * @param {string} serviceName - Name of the service
 * @param {string} serviceUrl - URL of the service
 * @returns {Promise<boolean>} Whether the service is healthy
 */
async function checkServiceHealth(serviceName, serviceUrl) {
  const url = `${getServiceUrl(serviceName, serviceUrl)}/health`;
  
  try {
    // Use the retry utility with shorter timeout for health checks
    const response = await retry(
      () => axios.get(url, {
        timeout: 3000,
        headers: { 'Accept': 'application/json' }
      }),
      {
        maxRetries: 1, // Just one retry for health checks
        interval: 500,
        shouldRetry: (error) => error.code !== 'ECONNREFUSED' // Don't retry if connection refused
      }
    );
    
    const isHealthy = response.status === 200;
    logger.debug(`Health check for ${serviceName}: ${isHealthy ? 'healthy' : 'unhealthy'}`, {
      url,
      status: response.status
    });
    
    return isHealthy;
  } catch (error) {
    logger.warn(`Health check failed for ${serviceName}`, {
      error: error.message,
      url,
      code: error.code
    });
    return false;
  }
}

/**
 * Checks health of all services
 * @returns {Promise<Object>} Health status of all services
 */
async function checkAllServices() {
  const services = Object.values(config.SERVICES);
  const healthStatus = {};
  
  // Run health checks in parallel for faster response
  const healthChecks = services.map(async (service) => {
    const isHealthy = await checkServiceHealth(
      service.name,
      service.url
    );
    
    return { name: service.name, isHealthy };
  });
  
  const results = await Promise.all(healthChecks);
  
  // Convert results to expected format
  results.forEach(result => {
    healthStatus[result.name] = result.isHealthy;
  });
  
  return healthStatus;
}

/**
 * Waits for all services to be healthy with improved retry logic
 * @param {number} maxRetries - Maximum number of retries
 * @param {number} interval - Interval between retries in ms
 * @returns {Promise<boolean>} Whether all services are healthy
 */
async function waitForServices(maxRetries = 20, interval = 3000) {
  const services = Object.values(config.SERVICES);
  let retries = 0;
  
  // Filter required services
  const requiredServices = services.filter(service => service.required);
  
  if (requiredServices.length === 0) {
    logger.warn('No required services found. Gateway may not function correctly.');
    return true;
  }
  
  logger.info('Waiting for required services to be ready...', {
    services: requiredServices.map(s => s.name).join(', '),
    maxRetries,
    interval
  });
  
  while (retries < maxRetries) {
    const results = await Promise.all(
      requiredServices.map(async (service) => {
        const isHealthy = await checkServiceHealth(service.name, service.url);
        return { name: service.name, isHealthy, required: service.required };
      })
    );
    
    const requiredAndHealthy = results
      .filter(result => result.required)
      .every(result => result.isHealthy);
    
    if (requiredAndHealthy) {
      logger.info('All required services are healthy!');
      return true;
    }
    
    const unhealthyServices = results
      .filter(result => result.required && !result.isHealthy)
      .map(result => result.name);
    
    logger.warn(`Waiting for required services: ${unhealthyServices.join(', ')}`, {
      retry: retries + 1,
      maxRetries,
      nextRetryMs: interval
    });
    
    retries++;
    
    if (retries < maxRetries) {
      await new Promise((resolve) => setTimeout(resolve, interval));
    }
  }
  
  const finalResults = await Promise.all(
    requiredServices.map(async (service) => {
      const isHealthy = await checkServiceHealth(service.name, service.url);
      return { name: service.name, isHealthy };
    })
  );
  
  const unhealthyServices = finalResults
    .filter(result => !result.isHealthy)
    .map(result => result.name);
  
  logger.error('Some required services failed to become healthy', {
    unhealthyServices,
    maxRetries,
    interval
  });
  
  return false;
}

module.exports = {
  checkServiceHealth,
  checkAllServices,
  waitForServices
};