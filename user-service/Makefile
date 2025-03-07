.PHONY: test test-unit test-integration coverage generate

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	go test -v ./tests/unit/...

# Run integration tests
test-integration:
	go test -v ./tests/integration/...

# Run tests with short flag (skips integration tests)
test-short:
	go test -v -short ./tests/...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html

# Generate mocks
# Requires mockery: go install github.com/vektra/mockery/v2@latest
generate-mocks:
	mockery --dir=internal/services --name=UserServiceIntf --output=tests/mocks --outpkg=mocks

# Install dependencies
deps:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	go get github.com/stretchr/testify/require
	go get github.com/joho/godotenv