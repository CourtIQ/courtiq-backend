import { readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { execSync } from 'child_process';

interface ServiceConfig {
  name: string;
  port: number;
  internal_port: number;
  graphql_path: string;
  env_variables: {
    [key: string]: {
      NODE_ENV?: string;
      GO_ENV?: string;
      GRAPHQL_PLAYGROUND: boolean;
      LOG_LEVEL: string;
      GIN_MODE?: string;
    };
  };
}

interface Config {
  app: {
    name: string;
    domain: string;
  };
  services: {
    [key: string]: ServiceConfig;
  };
  environments: {
    [key: string]: {
      base_url: string;
    };
  };
}

// Load environment from command line or default to development
const env = process.env.APP_ENV || 'development';

// Load configuration
const config: Config = JSON.parse(
  readFileSync(join(__dirname, '../config/config.json'), 'utf8')
);

// Function to get service URL based on environment
function getServiceUrl(serviceName: string, port: number, internal_port: number): string {
  if (env === 'development') {
    // In Docker development, use service names
    return `http://${serviceName}:${internal_port}`;
  }
  return config.environments[env].base_url.replace('{service}', serviceName);
}

// Function to generate base environment variables for a service
function generateBaseEnv(serviceKey: string, service: ServiceConfig): string {
  const envVars: Record<string, any> = {
    // Basic service info
    SERVICE_NAME: service.name,
    PORT: service.internal_port,
    DOCKER_ENV: 'true',
    
    // Environment-specific variables from config
    ...service.env_variables[env],
    // Convert boolean to string for GRAPHQL_PLAYGROUND
    GRAPHQL_PLAYGROUND: service.env_variables[env].GRAPHQL_PLAYGROUND.toString(),
    
    // Add service URLs for API Gateway if it's the gateway service
    ...(serviceKey === 'api_gateway' ? {
      ...Object.entries(config.services).reduce((acc, [key, svc]) => ({
        ...acc,
        [`${key.toUpperCase()}_URL`]: getServiceUrl(svc.name, svc.port, svc.internal_port)
      }), {})
    } : {})
  };

  return Object.entries(envVars)
    .filter(([_, value]) => value !== undefined)
    .map(([key, value]) => `${key}=${value}`)
    .join('\n');
}

// Function to generate root Docker environment file
function generateDockerEnv(): string {
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

// Function to get secrets from 1Password
async function getSecrets(): Promise<string> {
  try {
    const secrets = JSON.parse(
      readFileSync(join(__dirname, '../config/secrets.json'), 'utf8')
    );

    const result = execSync(
      `op item get "${secrets.onepassword.items[env].name}" --vault "${secrets.onepassword.vault}" --format json`
    ).toString();

    const item = JSON.parse(result);
    return item.fields
      .filter((field: any) => 
        field.type === "CONCEALED" && 
        field.value !== undefined &&
        field.label !== "password"
      )
      .map((field: any) => `${field.label}=${field.value}`)
      .join('\n');

  } catch (error) {
    console.error('Error getting secrets:', error);
    throw error;
  }
}

// Main function to generate environment files
async function generateEnvironments() {
  try {
    // Check 1Password CLI
    try {
      execSync('op --version');
    } catch {
      throw new Error('1Password CLI not installed or not in PATH');
    }

    // Check 1Password authentication
    try {
      execSync('op account list');
    } catch {
      throw new Error('Not logged into 1Password CLI. Run: op signin');
    }

    console.log(`Generating ${env} environment configuration...`);

    // Get secrets first
    const secrets = await getSecrets();

    // Generate Docker root .env file
    const rootEnvContent = generateDockerEnv();
    writeFileSync(join(__dirname, '..', '.env'), rootEnvContent);
    console.log('Generated root .env file for Docker');

    // Generate environment files for each service
    for (const [serviceKey, service] of Object.entries(config.services)) {
      const baseEnv = generateBaseEnv(serviceKey, service);
      const envContent = `# Generated for ${env} environment\n\n${baseEnv}\n\n# Secrets\n${secrets}`;
      
      const envPath = join(__dirname, '..', service.name, '.env');
      writeFileSync(envPath, envContent);
      console.log(`Generated .env for ${service.name}`);
    }

    // Log configuration summary
    console.log('\nEnvironment Configuration Summary:');
    console.log(`Environment: ${env}`);
    console.log('Service URLs:');
    Object.entries(config.services).forEach(([key, service]) => {
      console.log(`  - ${service.name}: ${getServiceUrl(service.name, service.port, service.internal_port)}`);
    });

    console.log('\nEnvironment generation complete!');

  } catch (error) {
    console.error('Failed to generate environments:', error);
    process.exit(1);
  }
}

// Execute the script
generateEnvironments();