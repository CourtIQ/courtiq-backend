# Court IQ Makefile

# Default environment
ENV ?= development

# Paths
SCRIPTS_DIR := scripts
NODE_BIN := node_modules/.bin
TS_NODE := $(NODE_BIN)/ts-node

# Colors for output
YELLOW := \033[1;33m
GREEN := \033[1;32m
RED := \033[1;31m
NC := \033[0m

# Logging functions
define log
echo "$(GREEN)>>> $(1)$(NC)"
endef

define warn
echo "$(YELLOW)>>> WARNING: $(1)$(NC)"
endef

define error
echo "$(RED)>>> ERROR: $(1)$(NC)"
endef

# ====================
# Environment Setup
# ====================

# Initialize environment and generate .env files
.PHONY: init
init:
	$(call log, "Initializing $(ENV) environment...")
	./$(SCRIPTS_DIR)/setup.sh $(ENV)

# Shortcuts for initializing specific environments
.PHONY: init-dev init-staging init-prod
init-dev:
	$(MAKE) init ENV=development

init-staging:
	$(MAKE) init ENV=staging

init-prod:
	$(MAKE) init ENV=production

# Verify all required .env files are present; if not, run init
.PHONY: verify-env
verify-env:
	$(call log, "Verifying environment files")
	@services=("api-gateway" "user-service" "relationship-service" "matchup-service" "equipment-service"); \
	for service in $${services[@]}; do \
		if [ ! -f "$$service/.env" ]; then \
			$(call warn, "$$service missing .env file. Running initialization..."); \
			$(MAKE) init ENV=$(ENV); \
		else \
			$(call log, "$$service .env file found"); \
		fi; \
	done

# ====================
# Docker Commands
# ====================

# Build a specific service, verifying .env files first
.PHONY: build
build: verify-env
	@if [ -z "$(SERVICE)" ]; then \
		$(call error, "SERVICE parameter is required. Usage: make build SERVICE=<service-name>"); \
		exit 1; \
	fi
	@if [ ! -f "$(SERVICE)/Dockerfile" ]; then \
		$(call error, "Dockerfile not found for $(SERVICE)"); \
		exit 1; \
	fi
	$(call log, "Building $(SERVICE) for $(ENV) environment")
	docker compose build $(SERVICE)

# Build all services after verifying .env files
.PHONY: build-all
build-all: verify-env
	$(call log, "Building all services for $(ENV) environment")
	docker compose build

# Start a specific service
.PHONY: start-service
start-service: verify-env
	@if [ -z "$(SERVICE)" ]; then \
		$(call error, "SERVICE parameter is required. Usage: make start-service SERVICE=<service-name>"); \
		exit 1; \
	fi
	$(call log, "Starting $(SERVICE) in $(ENV) environment")
	docker compose up -d $(SERVICE)

# Start all services
.PHONY: start-all
start-all: verify-env
	$(call log, "Starting all services in $(ENV) environment")
	docker compose up -d

# Stop a specific service
.PHONY: stop-service
stop-service:
	@if [ -z "$(SERVICE)" ]; then \
		$(call error, "SERVICE parameter is required. Usage: make stop-service SERVICE=<service-name>"); \
		exit 1; \
	fi
	$(call log, "Stopping $(SERVICE)")
	docker compose stop $(SERVICE)

# Stop all services
.PHONY: stop-all
stop-all:
	$(call log, "Stopping all services")
	docker compose down

# Restart a specific service
.PHONY: restart-service
restart-service: verify-env
	@if [ -z "$(SERVICE)" ]; then \
		$(call error, "SERVICE parameter is required. Usage: make restart-service SERVICE=<service-name>"); \
		exit 1; \
	fi
	$(call log, "Restarting $(SERVICE)")
	docker compose restart $(SERVICE)

# ====================
# Logs & Monitoring
# ====================

# View logs for a specific service
.PHONY: logs
logs:
	@if [ -z "$(SERVICE)" ]; then \
		$(call error, "SERVICE parameter is required. Usage: make logs SERVICE=<service-name>"); \
		exit 1; \
	fi
	docker compose logs -f $(SERVICE)

# View logs for all services
.PHONY: logs-all
logs-all:
	$(call log, "Showing logs for all services")
	docker compose logs -f

# Show Docker container status
.PHONY: ps
ps:
	docker compose ps

# ====================
# Cleanup Commands
# ====================

# Clean environment files
.PHONY: clean
clean:
	$(call log, "Cleaning environment files")
	find . -name ".env" -type f -delete
	$(call log, "Clean complete")

# Clean Docker resources
.PHONY: docker-clean
docker-clean:
	$(call log, "Cleaning Docker resources")
	docker compose down -v --remove-orphans
	docker system prune -f
	$(call log, "Docker clean complete")

# Clean both environment files and Docker resources
.PHONY: clean-all
clean-all: clean docker-clean
	$(call log, "Full cleanup complete")

# ====================
# Help
# ====================

# Display help for available commands
.PHONY: help
help:
	@echo "$(YELLOW)Court IQ Makefile Commands:$(NC)"
	@echo ""
	@echo "Environment Setup:"
	@echo "  make init ENV=<env>       - Initialize environment"
	@echo "  make init-dev             - Initialize development environment"
	@echo "  make init-staging         - Initialize staging environment"
	@echo "  make init-prod            - Initialize production environment"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build SERVICE=<name> - Build a specific service"
	@echo "  make build-all            - Build all services"
	@echo "  make start-service SERVICE=<name> - Start a specific service"
	@echo "  make stop-service SERVICE=<name>  - Stop a specific service"
	@echo "  make restart-service SERVICE=<name> - Restart a specific service"
	@echo "  make start-all            - Start all services"
	@echo "  make stop-all             - Stop all services"
	@echo ""
	@echo "Logs & Monitoring:"
	@echo "  make logs SERVICE=<name>  - View logs for a service"
	@echo "  make logs-all             - View logs for all services"
	@echo "  make ps                   - Show running containers"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean                - Clean environment files"
	@echo "  make docker-clean         - Clean Docker resources"
	@echo "  make clean-all            - Clean everything"
	@echo ""
	@echo "Usage:"
	@echo "  - Set environment with ENV=<development|staging|production>"
	@echo "  - Example: make start-all ENV=development"
