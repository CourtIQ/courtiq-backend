## Local Development

This project uses a Makefile to simplify common development tasks and integrates with 1Password for secure environment variable management.

### Prerequisites

- Docker and Docker Compose
- Make
- 1Password CLI (configured and authenticated)

### Available Commands

To see all available commands, run:

```
make help
```

### Starting the Application

To start the application:

```
make start
```

This command will:
1. Populate environment variables from 1Password using the `populate-env.sh` script.
2. Start all services using Docker Compose.

### Other Useful Commands

- Stop the application: `make stop`
- Restart the application: `make restart`
- View logs: `make logs`
- Rebuild Docker images: `make build`
- Run tests: `make test`
- Clean up generated files and Docker resources: `make clean`
- Manually populate environment variables: `make populate-env`

### Development Workflow

1. Ensure your 1Password CLI is authenticated.
2. Start the application with `make start`
3. Make changes to your code
4. Rebuild and restart the affected service(s) with `make build` followed by `make restart`
5. View logs with `make logs` to check for any issues
6. Run tests with `make test` to ensure everything is working correctly
7. When you're done, stop the application with `make stop`

Remember to never commit the `.env.populated` file to version control.