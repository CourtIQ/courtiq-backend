# CourtIQ Backend

CourtIQ Backend is a microservices-based application designed to provide robust and scalable APIs for court-related data and functionalities.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Credential Management](#credential-management)
  - [MongoDB Credentials Setup](#mongodb-credentials-setup)
  - [Firebase Configuration](#firebase-configuration)
  - [Populating Environment Variables](#populating-environment-variables)
  - [Updating Credentials](#updating-credentials)
- [Setup Instructions](#setup-instructions)
  - [Clone the Repository](#clone-the-repository)
  - [Install Dependencies](#install-dependencies)
- [Running Locally](#running-locally)
- [Deployment](#deployment)
- [Makefile Commands](#makefile-commands)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

---

## Prerequisites

- **Docker** installed and running
- **Docker Compose** installed
- **1Password CLI** installed and authenticated
- **Go** installed (if you need to build services locally)
- **Make** installed

## Credential Management

This project uses 1Password for secure credential management. Credentials are stored in a shared vault called **"CourtIQ Backend"**.

### MongoDB Credentials Setup

1. In the 1Password **"CourtIQ Backend"** vault, there should be an item named **"courtiq-mongodb-user"** with:
   - A **"username"** field containing the MongoDB username
   - A **"password"** field containing the MongoDB password

2. The `.env` file references these credentials as:

   ```dotenv
   MONGO_USERNAME=op:courtiq-mongodb-user/username
   MONGO_PASSWORD=op:courtiq-mongodb-user/password
