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

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/court-iq-backend.git
   cd court-iq-backend
