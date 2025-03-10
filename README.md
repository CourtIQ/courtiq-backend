
# Court IQ Backend

This is the backend system for Court IQ, built with a microservices architecture and managed through Docker and Makefile commands. This document explains how to set up, build, and manage the environment for development, staging, and production.

## Table of Contents
1. [Project Structure](#project-structure)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Environment Setup](#environment-setup)
5. [Docker Commands](#docker-commands)
6. [Logs & Monitoring](#logs--monitoring)
7. [Cleanup Commands](#cleanup-commands)
8. [Usage Examples](#usage-examples)

---

## Project Structure

The main directories and files for this project:
- **api-gateway**: The API Gateway for routing requests across microservices.
- **user-service, relationship-service, matchup-service, equipment-service**: Microservices that make up the backend.
- **scripts**: Contains setup scripts for initializing environment files.
- **Makefile**: Provides easy-to-use commands to build, manage, and monitor services.

---

## Requirements

Ensure you have the following installed:
- **Docker**: For containerization and deployment.
- **Make**: To use the `Makefile` commands.
- **Node.js**: For running and managing the API Gateway.
- **Go**: For building and running Go-based microservices.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/court-iq-backend.git
   cd court-iq-backend
   ```

2. **Install Dependencies**:
   Install Node.js dependencies in each microservice directory if required. Docker will also handle this during the build phase.

---

## Environment Setup

Before running the services, you need to initialize environment files. The `Makefile` provides commands to set up environments for development, staging, and production.

### Initialize Environment

To initialize environment files:
```bash
make init ENV=development
```

Alternatively, you can use:
```bash
make init-dev       # Initializes the development environment
make init-staging   # Initializes the staging environment
make init-prod      # Initializes the production environment
```

---

## Docker Commands

The Makefile includes Docker commands to manage services.

### Build a Specific Service

```bash
make build SERVICE=<service-name>
```

For example, to build the `api-gateway` service:
```bash
make build SERVICE=api-gateway
```

### Build All Services

```bash
make build-all
```

### Start a Specific Service

```bash
make start-service SERVICE=<service-name>
```

### Start All Services

```bash
make start-all
```

### Stop a Specific Service

```bash
make stop-service SERVICE=<service-name>
```

### Stop All Services

```bash
make stop-all
```

---

## Logs & Monitoring

### View Logs for a Specific Service

```bash
make logs SERVICE=<service-name>
```

### View Logs for All Services

```bash
make logs-all
```

### Show Docker Container Status

```bash
make ps
```

---

## Cleanup Commands

### Clean Environment Files

```bash
make clean
```

### Clean Docker Resources

```bash
make docker-clean
```

### Full Cleanup (Environment Files and Docker Resources)

```bash
make clean-all
```

---

## Usage Examples

To start the development environment with all services running:
```bash
make start-all ENV=development
```

To build and start a single service in staging:
```bash
make build SERVICE=api-gateway ENV=staging
make start-service SERVICE=api-gateway ENV=staging
```

To view logs for `user-service`:
```bash
make logs SERVICE=user-service
```

For more information on each command, refer to the `Makefile`.
