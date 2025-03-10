#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print with color
print_green() { echo -e "${GREEN}$1${NC}"; }
print_yellow() { echo -e "${YELLOW}$1${NC}"; }
print_red() { echo -e "${RED}$1${NC}"; }

# Check for required tools
check_requirements() {
  print_green "Checking requirements..."

  # Check if 1Password CLI is installed
  if ! command -v op &> /dev/null; then
    print_red "1Password CLI not found. Please install it first."
    print_yellow "Visit https://1password.com/downloads/command-line/ for installation instructions."
    exit 1
  fi
  print_green "✓ 1Password CLI is installed"

  # Check if logged in to 1Password
  if ! op account list &> /dev/null; then
    print_red "Please login to 1Password CLI first: op signin"
    exit 1
  fi
  print_green "✓ Logged in to 1Password"

  # Check if Node.js is installed
  if ! command -v node &> /dev/null; then
    print_red "Node.js not found. Please install it first."
    exit 1
  fi
  print_green "✓ Node.js is installed"
}

# Get environment argument
ENV=${1:-development}
if [[ ! "$ENV" =~ ^(development|staging|production)$ ]]; then
    print_red "Invalid environment. Use: development, staging, or production"
    exit 1
fi

print_green "Setting up $ENV environment..."

# Check requirements
check_requirements

# Change to scripts directory
cd "$(dirname "$0")"

# Generate environment files
print_green "Generating environment files..."
APP_ENV=$ENV node env.js

# Verify all services have .env files
cd ..
services=("api-gateway" "user-service" "relationship-service" "matchup-service" "equipment-service" "search-service")
for service in "${services[@]}"; do
    if [ ! -f "$service/.env" ]; then
        print_red "Error: .env file not generated for $service"
        exit 1
    else
        print_green "✓ .env file generated for $service"
    fi
done

print_green "Setup complete! Environment files generated for $ENV"
print_yellow "To start the services, run: make start-all"