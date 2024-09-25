# Variables
DOCKER = docker
DOCKER_COMPOSE = docker-compose
COMPOSE_FILE = deploy/docker-compose.yml
ENV_POPULATOR = deploy/populate-env.sh
ENV_POPULATED = deploy/.env.populated

# Phony targets
.PHONY: build start stop restart logs clean pull status test shell secrets

# Default target
all: start

# Populate secrets from 1Password
secrets:
	@echo "Populating secrets from 1Password..."
	$(ENV_POPULATOR)

# Build or rebuild services
build: secrets
	@echo "Building fresh Docker images..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build --no-cache

# Start services
start: secrets
	@echo "Starting services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

# Stop services
stop:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down
	@echo "Removing populated .env file..."
	rm -f $(ENV_POPULATED)

# Restart services
restart: stop start

# View output from containers
logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Remove stopped containers, networks, volumes, images, and populated env file
clean: stop
	@echo "Cleaning up Docker resources..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down -v --rmi all --remove-orphans
	@echo "Removing populated .env file..."
	rm -f $(ENV_POPULATED)

# Pull latest images
pull:
	@echo "Pulling latest Docker images..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) pull

# Show status of services
status:
	@echo "Showing status of services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) ps

# Run tests (adjust as needed for your test setup)
test: secrets
	@echo "Running tests..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) run --rm test-service npm test

# Enter a specific service's container (usage: make shell SERVICE=service_name)
shell:
	@echo "Entering shell for $(SERVICE)..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) exec $(SERVICE) /bin/sh

# Remove the populated env file
clean-secrets:
	@echo "Removing populated .env file..."
	rm -f $(ENV_POPULATED)