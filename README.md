## Local Development

This project uses a Makefile to simplify common development tasks. Ensure you have Docker, Docker Compose, and Make installed on your system.

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

This command will populate the environment variables from 1Password and start all services using Docker Compose.

### Other Useful Commands

- Stop the application: `make stop`
- Restart the application: `make restart`
- View logs: `make logs`
- Rebuild Docker images: `make build`
- Run tests: `make test`
- Clean up generated files and Docker resources: `make clean`

### Development Workflow

1. Start the application with `make start`
2. Make changes to your code
3. Rebuild and restart the affected service(s) with `make build` followed by `make restart`
4. View logs with `make logs` to check for any issues
5. Run tests with `make test` to ensure everything is working correctly
6. When you're done, stop the application with `make stop`

Remember to never commit the `.env.populated` file to version control.