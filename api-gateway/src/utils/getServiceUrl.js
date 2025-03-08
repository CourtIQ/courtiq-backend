const logger = require('../logging/logger');

/**
 * Gets the URL for a service with validation
 * @param {string} serviceName - Name of the service
 * @param {string} envUrl - URL from environment variables
 * @returns {string} The URL to use for the service
 */
function getServiceUrl(serviceName, envUrl) {
  if (!serviceName) {
    throw new Error('Service name is required');
  }
  
  // If environment URL is provided, use it
  if (envUrl) {
    // Validate URL format
    try {
      // Add protocol if missing
      const urlWithProtocol = envUrl.startsWith('http') ?
        envUrl : `http://${envUrl}`;
        
      // Test URL validity
      new URL(urlWithProtocol);
      
      logger.debug(`Using environment URL for ${serviceName}`, { url: envUrl });
      return envUrl;
    } catch (error) {
      logger.warn(`Invalid URL format for ${serviceName}, falling back to default`, {
        providedUrl: envUrl,
        error: error.message
      });
    }
  }
  
  // Default URL for Kubernetes/Docker service discovery
  const url = `http://${serviceName}:8080`;
  
  logger.debug(`Using default URL for ${serviceName}`, {
    url,
    fromEnv: false
  });
  
  return url;
}

module.exports = getServiceUrl;