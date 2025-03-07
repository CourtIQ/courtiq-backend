# CourtIQ Shared Package Testing Guide

This document describes how to run and maintain tests for the shared packages.

## Running Tests

### Basic Test Run

To run all tests:

```bash
cd shared
make test
```

### Verbose Output

For more detailed output:

```bash
make test-verbose
```

### Test Coverage

To generate a test coverage report:

```bash
make test-coverage
```

This will create an HTML report at `coverage.html` that you can view in your browser.

## Setting Up Test Environment

Some tests require a MongoDB instance. To run those tests, set the `MONGODB_TEST_URL` environment variable:

```bash
# Using a local MongoDB instance
export MONGODB_TEST_URL="mongodb://localhost:27017"

# Or for CI environments, using an ephemeral test database
export MONGODB_TEST_URL="mongodb://test:password@testhost:27017/test-db?authSource=admin"
```

If this variable is not set, tests that require MongoDB will be skipped.

## Test Conventions

1. **Unit Tests**: All packages should have unit tests with >80% coverage
2. **File Naming**: Tests should be in files named `*_test.go` in the same package as the code being tested
3. **Test Naming**: Test functions should be named `Test<FunctionName>` and follow Go test conventions
4. **Test Independence**: Tests should not depend on external services or state from other tests

## Adding New Tests

When adding a new feature to the shared package:

1. Create a corresponding test file if it doesn't exist
2. Add tests for all public functions and methods
3. Use mocks for external dependencies
4. Run `make test-coverage` to ensure adequate coverage

## Mocking

For tests that require external dependencies like MongoDB, we use mocks. See the repository tests for examples of how to create and use mocks.

## Continuous Integration

Tests are automatically run in the CI pipeline for:
- All pull requests
- All pushes to main branches
- All release tags

If tests fail in CI, the build will be marked as failed.

## Dependencies

Test-specific dependencies are:
- `github.com/stretchr/testify` - For assertions and mocks
- `go.mongodb.org/mongo-driver` - For MongoDB integration tests

These are already included in the `go.mod` file.