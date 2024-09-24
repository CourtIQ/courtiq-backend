.PHONY: start stop restart clean logs build test help populate-env

# Default target
.DEFAULT_GOAL := help

# Variables
DEPLOY_DIR = deploy

# Colors for output
YELLOW := \033[1;33m
NC := \033[0m # No Color

populate-env: ## Populate environment variables from 1Password
	@echo "$(YELLOW)Populating environment variables...$(NC)"
	@cd $(DEPLOY_DIR) && ./populate-env.sh

start: populate-env ## Start the application
	@echo "$(YELLOW)Starting the application...$(NC)"
	@cd $(DEPLOY_DIR) && docker-compose up -d

stop: ## Stop the application
	@echo "$(YELLOW)Stopping the application...$(NC)"
	@cd $(DEPLOY_DIR) && docker-compose down

restart: stop start ## Restart the application

clean: ## Remove generated files and stop containers
	@echo "$(YELLOW)Cleaning up...$(NC)"
	@cd $(DEPLOY_DIR) && rm -f .env.populated
	@$(MAKE) stop
	@docker system prune -f

logs: ## View logs of all services
	@cd $(DEPLOY_DIR) && docker-compose logs -f

build: ## Rebuild the Docker images
	@echo "$(YELLOW)Rebuilding Docker images...$(NC)"
	@cd $(DEPLOY_DIR) && docker-compose build

test: ## Run tests (implement per your testing strategy)
	@echo "$(YELLOW)Running tests...$(NC)"
	@echo "Implement your test command here"

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-30s$(NC) %s\n", $$1, $$2}'