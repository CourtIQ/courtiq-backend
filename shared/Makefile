.PHONY: test test-verbose test-coverage

test:
	@echo "Running tests..."
	@go test ./pkg/... -race

test-verbose:
	@echo "Running tests with verbose output..."
	@go test ./pkg/... -v -race

test-coverage:
	@echo "Running tests with coverage report..."
	@go test ./pkg/... -race -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-update:
	@echo "Updating test dependencies..."
	@go get -t -u ./pkg/...
	@go mod tidy