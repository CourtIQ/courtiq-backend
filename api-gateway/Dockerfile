# Stage 1: Build the Go application
FROM golang:1.23-buster AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-binary .

# Stage 2: Run the application
FROM debian:buster-slim

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/service-binary .

# Expose the port specified by the environment variable
ENV PORT=${PORT:-8080}
EXPOSE ${PORT}

# Command to run the executable
CMD ["./service-binary"]
