FROM golang:1.23-alpine AS builder

# Add git for private dependencies if needed
RUN apk add --no-cache git

WORKDIR /app

# Install dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main ./cmd/main.go

# Use distroless as minimal base image 
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env file if needed
COPY --from=builder /app/.env .

# Use non-root user for security
USER nonroot:nonroot

# Expose port if needed
EXPOSE 8080

# Default command
CMD ["./main"]