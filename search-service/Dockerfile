FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Now copy the service's files
COPY search-service/go.mod search-service/go.sum ./
RUN go mod download

# Copy the rest of service code
COPY search-service/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Production image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]