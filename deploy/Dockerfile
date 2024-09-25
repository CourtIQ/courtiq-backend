# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Install git and other dependencies
RUN apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
ARG SERVICE_PORT

# Build the Go project binary
RUN CGO_ENABLED=0 GOOS=linux go build -o service .

# Final stage
FROM alpine:3.14

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/service .

# Install curl for health checks
RUN apk add --no-cache curl

# Ensure the binary is executable
RUN chmod +x service

ARG SERVICE_NAME
ARG SERVICE_PORT

ENV SERVICE_NAME=${SERVICE_NAME}
ENV SERVICE_PORT=${SERVICE_PORT}

EXPOSE ${SERVICE_PORT}

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${SERVICE_PORT}/health || exit 1

CMD ["/bin/sh", "-c", "./service"]