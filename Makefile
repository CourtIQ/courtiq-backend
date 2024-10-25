# Variables
DOCKER = docker
DOCKER_COMPOSE = docker-compose
COMPOSE_FILE = deploy/docker-compose.yml
ENV_POPULATOR = deploy/populate-env.sh
ENV_POPULATED = deploy/.env.populated

# Phony targets
.PHONY: build start stop restart logs clean pull status test shell secrets refresh refresh-service

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

# Complete refresh cycle (stop, rebuild, start)
refresh: 
	@echo "🔄 Starting complete refresh cycle..."
	@echo "🛑 Stopping all services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down
	@echo "🧹 Cleaning up..."
	rm -f $(ENV_POPULATED)
	@echo "🔑 Regenerating secrets..."
	$(ENV_POPULATOR)
	@echo "🏗️  Rebuilding services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build --no-cache
	@echo "🚀 Starting services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d
	@echo "📝 Showing logs..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Refresh specific service
refresh-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Error: SERVICE not specified. Usage: make refresh-service SERVICE=service-name"; \
		exit 1; \
	fi
	@echo "🔄 Refreshing service: $(SERVICE)"
	@echo "🛑 Stopping service..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) stop $(SERVICE)
	@echo "🏗️  Rebuilding service..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build --no-cache $(SERVICE)
	@echo "🚀 Starting service..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d $(SERVICE)
	@echo "📝 Showing service logs..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f $(SERVICE)

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

	# Deploy to development
.PHONY: deploy-dev
deploy-dev: secrets
	@echo "Deploying to development environment..."
	./deploy/scripts/deploy.sh development

# Deploy to production
.PHONY: deploy-prod
deploy-prod: secrets
	@echo "Deploying to production environment..."
	./deploy/scripts/deploy.sh production

# Get cluster credentials
.PHONY: cluster-credentials
cluster-credentials:
	gcloud container clusters get-credentials court-iq-cluster --region us-central1

# Show cluster status
.PHONY: cluster-status
cluster-status: cluster-credentials
	@echo "Cluster Status:"
	kubectl get nodes
	@echo "\nPods Status:"
	kubectl get pods -n court-iq
	@echo "\nServices Status:"
	kubectl get services -n court-iq
	@echo "\nIngress Status:"
	kubectl get ingress -n court-iq