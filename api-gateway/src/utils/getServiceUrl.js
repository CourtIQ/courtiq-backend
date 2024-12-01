function getServiceUrl(serviceName, envUrl) {
  return envUrl || `http://${serviceName}:8080`;
}

module.exports = getServiceUrl;
