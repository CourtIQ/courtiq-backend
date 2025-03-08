# CourtIQ API Gateway

The API Gateway serves as the central entry point for the CourtIQ backend microservices architecture, providing a unified GraphQL API using Apollo Federation.

## Features

- **Apollo Federation**: Creates a unified GraphQL API from multiple microservices
- **Authentication**: Firebase authentication with token validation and claims forwarding
- **Observability**: Structured logging and Prometheus metrics for monitoring
- **Resilience**: Automatic retry logic for service communication
- **Graceful Shutdown**: Proper shutdown handling with timeout protection

## Getting Started

### Prerequisites

- Node.js 18 or higher
- Docker and Docker Compose
- Firebase account with proper configuration

### Installation

1. Clone the repository
2. Create a `.env` file based on `.env.template`
3. Start the services with Docker Compose:
   ```
   docker-compose up -d
   ```

## Configuration

The API Gateway can be configured through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_NAME` | Name of the service | `api-gateway` |
| `PORT` | Port to listen on | `3000` |
| `NODE_ENV` | Environment (development/production) | `development` |
| `LOG_LEVEL` | Logging level | `info` (prod), `debug` (dev) |
| `GRAPHQL_PLAYGROUND` | Enable GraphQL playground | `false` |
| `METRICS_PATH` | Path for metrics endpoint | `/metrics` |
| `REQUEST_TIMEOUT` | Timeout for service requests | `30000` |
| `*_SERVICE_URL` | URLs for each microservice | See `.env.template` |

## Metrics and Monitoring

### Available Metrics

The API Gateway exposes Prometheus metrics at `http://localhost:3001/metrics` for monitoring:

- **HTTP request metrics**: Duration and counts by route, method, and status
- **GraphQL operation metrics**: Timing and counts by operation type
- **Microservice metrics**: Response times and error counts by service
- **System metrics**: Memory, CPU, and other Node.js metrics

### Viewing Metrics

1. **Directly**: Access `http://localhost:3001/metrics` in your browser to see raw Prometheus metrics
2. **Using Prometheus**: Configure Prometheus to scrape this endpoint
3. **Using Grafana**: Set up Grafana dashboards with Prometheus as the data source

#### Setting up Prometheus and Grafana

You can add Prometheus and Grafana to your docker-compose.yml:

```yaml
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus:/etc/prometheus
  ports:
    - 9090:9090
  depends_on:
    - api-gateway

grafana:
  image: grafana/grafana
  ports:
    - 3100:3000
  volumes:
    - grafana-storage:/var/lib/grafana
  depends_on:
    - prometheus

volumes:
  grafana-storage:
```

Then create a `prometheus/prometheus.yml` file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api-gateway'
    scrape_interval: 5s
    static_configs:
      - targets: ['api-gateway:3001']
```

Access Grafana at http://localhost:3100 (default credentials: admin/admin), add Prometheus as a data source (URL: http://prometheus:9090), and create dashboards.

## Authentication Flow

1. Clients authenticate with Firebase and obtain an ID token
2. The token is sent to the API Gateway in the Authorization header:
   ```
   Authorization: Bearer <firebase-token>
   ```
3. The API Gateway verifies the token with Firebase
4. User claims are forwarded to downstream services in the `X-User-Claims` header

## Troubleshooting

### Common Issues

#### Cannot Access Metrics Endpoint

If you cannot access the metrics endpoint at http://localhost:3001/metrics:

1. **Port Exposure**: Make sure port 3001 is exposed in the docker-compose.yml file
2. **Container Running**: Verify the gateway container is running with `docker ps`
3. **Logs**: Check logs for any errors with `docker logs courtiq-api-gateway`
4. **Network**: Ensure no other service is using port 3001 on your host machine

#### Service Connection Issues

If the API Gateway can't connect to services:

1. **Docker Network**: Make sure all services are on the same Docker network
2. **Service Names**: Verify service hostnames match the names in docker-compose.yml
3. **Service Health**: Check if all services are running and healthy

#### Authentication Issues

If authentication isn't working:

1. **Firebase Config**: Verify Firebase configuration in environment variables
2. **Token Format**: Check that tokens are being sent with correct format
3. **Logs**: Look for Firebase verification errors in logs

## Development

### Local Development

For local development:

1. Install dependencies: `npm install`
2. Start in dev mode: `npm run dev`
3. Run tests: `npm test`

### Project Structure

- `/src`: Source code
  - `/auth`: Authentication handling
  - `/gateway`: Apollo Gateway configuration
  - `/health`: Health check implementation
  - `/logging`: Logging configuration
  - `/metrics`: Metrics collection
  - `/utils`: Utility functions

## License

Copyright © 2024 CourtIQ