#!/usr/bin/env node

/**
 * Environment Management Script for Court IQ
 *
 * This script generates environment files for local development and CI/CD
 * using configuration from config/ and secrets from 1Password.
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Configuration
const ENV = process.env.APP_ENV || 'development';
const CONFIG_PATH = path.join(__dirname, '../config/config.json');
const SECRETS_PATH = path.join(__dirname, '../config/secrets.json');
const ROOT_DIR = path.join(__dirname, '..');

// Load configuration
const config = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
const secretsConfig = JSON.parse(fs.readFileSync(SECRETS_PATH, 'utf8'));

// Validate environment
if (!['development', 'staging', 'production'].includes(ENV)) {
  console.error(`Invalid environment: ${ENV}. Use development, staging, or production.`);
  process.exit(1);
}

/**
 * Get service URL based on environment
 */
function getServiceUrl(serviceName, port, internalPort) {
  if (ENV === 'development') {
    return `http://${serviceName}:${internalPort}`;
  }
  return config.environments[ENV].base_url.replace('{service}', serviceName);
}

/**
 * Generate base environment variables for a service
 */
function generateBaseEnv(serviceKey, service) {
  const serviceUrls = Object.entries(config.services).reduce((acc, [key, svc]) => ({
    ...acc,
    [`${key.toUpperCase()}_URL`]: getServiceUrl(svc.name, svc.port, svc.internal_port)
  }), {});

  const envVars = {
    // Basic service info
    SERVICE_NAME: service.name,
    PORT: service.internal_port,
    DOCKER_ENV: 'true',
    
    // Environment-specific variables from config
    ...service.env_variables[ENV],
    // Convert boolean to string
    GRAPHQL_PLAYGROUND: service.env_variables[ENV].GRAPHQL_PLAYGROUND.toString(),
    
    // Add service URLs for all services
    ...serviceUrls
  };

  return Object.entries(envVars)
    .filter(([_, value]) => value !== undefined)
    .map(([key, value]) => `${key}=${value}`)
    .join('\n');
}

/**
 * Generate Docker root environment file
 */
function generateDockerEnv() {
  return Object.entries(config.services)
    .map(([_, service]) => {
      const servicePrefix = service.name.toUpperCase().replace(/-/g, '_');
      return [
        `${servicePrefix}_PORT=${service.port}`,
        `${servicePrefix}_INTERNAL_PORT=${service.internal_port}`
      ].join('\n');
    })
    .join('\n');
}

/**
 * Get secrets from 1Password
 */
async function getSecrets() {
  try {
    // Check if 1Password CLI is installed
    try {
      execSync('op --version', { stdio: 'ignore' });
    } catch (error) {
      throw new Error('1Password CLI not installed. Please install it first.');
    }

    // Check if logged in
    try {
      execSync('op account list', { stdio: 'ignore' });
    } catch (error) {
      throw new Error('Not logged into 1Password CLI. Run: op signin');
    }

    const opConfig = secretsConfig.environments[ENV].onepassword;
    console.log(`Fetching secrets from 1Password vault "${opConfig.vault}", item "${opConfig.item}"...`);

    const result = execSync(
      `op item get "${opConfig.item}" --vault "${opConfig.vault}" --format json`
    ).toString();

    const item = JSON.parse(result);
    
    // Get all fields marked as passwords or concealed
    return item.fields
      .filter(field =>
        field.type === "CONCEALED" &&
        field.value !== undefined &&
        field.label !== "password" // Skip the main password field
      )
      .map(field => `${field.label}=${field.value}`)
      .join('\n');

  } catch (error) {
    console.error('Error getting secrets:', error.message);
    process.exit(1);
  }
}

/**
 * Main function to generate environment files
 */
async function generateEnvironments() {
  try {
    console.log(`Generating ${ENV} environment configuration...`);

    // Get secrets from 1Password
    const secrets = await getSecrets();

    // Generate Docker root .env file
    const rootEnvContent = generateDockerEnv();
    fs.writeFileSync(path.join(ROOT_DIR, '.env'), rootEnvContent);
    console.log('Generated root .env file for Docker');

    // Generate environment files for each service
    for (const [serviceKey, service] of Object.entries(config.services)) {
      const baseEnv = generateBaseEnv(serviceKey, service);
      const envContent = `# Generated for ${ENV} environment\n\n${baseEnv}\n\n# Secrets\n${secrets}`;
      
      const envPath = path.join(ROOT_DIR, service.name, '.env');
      fs.writeFileSync(envPath, envContent);
      console.log(`Generated .env for ${service.name}`);
    }

    console.log('\nEnvironment Configuration Summary:');
    console.log(`Environment: ${ENV}`);
    console.log('Service URLs:');
    Object.entries(config.services).forEach(([_, service]) => {
      console.log(`  - ${service.name}: ${getServiceUrl(service.name, service.port, service.internal_port)}`);
    });

    console.log('\nEnvironment generation complete!');

  } catch (error) {
    console.error('Failed to generate environments:', error.message);
    process.exit(1);
  }
}

// Execute the script
generateEnvironments();