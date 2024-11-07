#!/bin/bash

# Check if 1Password CLI is installed
if ! command -v op &> /dev/null; then
    echo "1Password CLI not found. Please install it first."
    exit 1
fi

# Check if logged in to 1Password
if ! op account list &> /dev/null; then
    echo "Please login to 1Password CLI first: op signin"
    exit 1
fi

# Get environment argument
ENV=${1:-development}
if [[ ! "$ENV" =~ ^(development|staging|production)$ ]]; then
    echo "Invalid environment. Use: development, staging, or production"
    exit 1
fi

echo "Setting up $ENV environment..."

# Change to scripts directory
cd "$(dirname "$0")"

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    npm install
fi

# Generate environment files
npm run "gen:$ENV"

# Verify all services have .env files
cd ..
services=("api-gateway" "user-service" "relationship-service" "matchup-service" "equipment-service")
for service in "${services[@]}"; do
    if [ ! -f "$service/.env" ]; then
        echo "Error: .env file not generated for $service"
        exit 1
    fi
done

echo "Setup complete! Environment files generated for $ENV"