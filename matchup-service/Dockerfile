# Stage 1: Build Go 1.23.1 from source
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o service-binary .

# Stage 2: Run the application
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/service-binary .

# Expose the port specified by the environment variable
ENV PORT=${PORT:-8084}
EXPOSE ${PORT}

# Command to run the executable
CMD ["./service-binary"]