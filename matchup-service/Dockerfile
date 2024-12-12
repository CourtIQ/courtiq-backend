# matchup-service/Dockerfile (same for other Go services)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Production image
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

# Default command
CMD ["./main"]