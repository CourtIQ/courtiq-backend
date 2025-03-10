# Court IQ Backend

A microservices-based backend for Court IQ, providing GraphQL APIs for tennis players to manage matches, equipment, and relationships.

## Architecture

The system is built with several microservices:

- **API Gateway**: Node.js-based API gateway that routes requests to appropriate services
- **User Service**: Go service for user management
- **Relationship Service**: Go service for managing relationships between users
- **Matchup Service**: Go service for tennis match management
- **Equipment Service**: Go service for managing tennis equipment
- **Search Service**: Go service for search functionality

All services use MongoDB for persistence and communicate via GraphQL with Apollo Federation.

## Environment Setup

The project uses 1Password for secure credential management in all environments.

### Prerequisites

- Docker and Docker Compose
- Node.js (v18+)
- Go (v1.23+)
- 1Password CLI (`op`) installed and configured
- Make

### Initial Setup

1. Log in to 1Password CLI:
   ```
   op signin
   ```

2. Initialize your development environment:
   ```
   make init-dev
   ```

   This command will:
   - Fetch secrets from 1Password
   - Generate `.env` files for all services
   - Configure Docker environment

3. Build and start all services:
   ```
   make build-all
   make start-all
   ```

4. Access the GraphQL playground at http://localhost:3000/graphql

### Environment Management

- **Development**: `make init-dev`
- **Staging**: `make init-staging`
- **Production**: `make init-prod`

### Common Commands

```
# Start all services
make start-all

# Start a specific service
make start-service SERVICE=api-gateway

# View logs
make logs-all
make logs SERVICE=user-service

# Stop all services
make stop-all

# Clean environment
make clean-all
```

See `make help` for more commands.

## Secret Management

All secrets are stored in 1Password. The required secret structure in 1Password vault items:

- `MONGODB_URL`: MongoDB connection string
- `FIREBASE_CONFIG`: Firebase configuration JSON
- `FIREBASE_SERVICE_ACCOUNT`: Firebase service account JSON
- `GOOGLE_PLACES_API_KEY`: Google Places API key

## CI/CD Pipeline

The project uses GitHub Actions for CI/CD:

1. **Pull Request**: Runs tests and checks builds
2. **Merge to Main**: Runs tests
3. **Tag Release**: Builds and deploys to Google Cloud Run

Secrets for CI/CD are loaded from 1Password into GitHub Actions.

## License

Proprietary, all rights reserved.