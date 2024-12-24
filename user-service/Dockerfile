FROM golang:1.23-alpine AS builder

WORKDIR /app

# First copy and build shared module
WORKDIR /app/shared
COPY shared/ .
RUN go mod download
RUN go build ./...

# Now build the service
WORKDIR /app/service
COPY user-service/ .
RUN go mod download
RUN go build -o main ./cmd/main.go

# Production image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service/main .
CMD ["./main"]