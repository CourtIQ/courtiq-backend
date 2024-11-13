function getServiceUrl(serviceName, envUrl) {
    if (envUrl) {
      return envUrl; // Use the provided URL from the environment variable
    }
    // Default to a development URL if envUrl is not provided
    return `http://${serviceName}:8080`; // Replace 8080 with the default port if necessary
}
  
module.exports = getServiceUrl;