# Testing Strategy for User Service

This document outlines the testing approach for the User Service microservice.

## Testing Levels

### Unit Tests

Unit tests focus on individual components in isolation, with dependencies mocked:

- Located in `tests/unit/`
- Fast to run, no external dependencies
- Mock the repository layer
- Test service logic only

### Integration Tests

Integration tests verify components work together with real dependencies:

- Located in `tests/integration/`
- Require a test MongoDB instance
- Test repository and service together
- Verify real database interactions

## Running Tests

### Prerequisites

- Go 1.18+
- MongoDB (for integration tests)
- Copy `.env.test.example` to `.env.test` and configure it

### Commands

Using the Makefile:

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage report
make coverage

# Run fast tests only (skips integration tests)
make test-short
```

Or using Go directly:

```bash
# Run all tests
go test -v ./tests/...

# Run unit tests
go test -v ./tests/unit/...

# Run integration tests
go test -v ./tests/integration/...

# Run with coverage
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

## Test Database

For integration tests, you can use:

1. **Local MongoDB**:
   ```
   TEST_MONGODB_URI=mongodb://localhost:27017/courtiq-test-db
   ```

2. **MongoDB Atlas Test Database**:
   ```
   TEST_MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/courtiq-test-db
   ```

**Important**: Integration tests will clear test collections before running!

## Mock Generation

We use testify/mock for mocking. To generate mocks:

```bash
make generate-mocks
```

## Testing Best Practices

1. **Isolation**: Unit tests should be isolated and not depend on external services
2. **Deterministic**: Tests should produce the same results each run
3. **Maintainable**: Use helpers and utilities to keep tests DRY
4. **Fast**: Unit tests should run quickly, integration tests can be skipped in CI with `-short` flag
5. **Coverage**: Aim for high test coverage, especially for critical paths