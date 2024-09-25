# Use the official Golang image, version 1.23.1, based on Alpine Linux for minimal image size
FROM golang:1.23.1-alpine

# Set the working directory inside the container to /app
WORKDIR /app

# Install git and other dependencies required for building Go projects
RUN apk add --no-cache git

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project source code
COPY . .

# Define build-time variables for service name and service port
ARG SERVICE_NAME
ARG SERVICE_PORT

# Use ENV to set these arguments as environment variables, available during runtime
ENV SERVICE_NAME=${SERVICE_NAME}
ENV SERVICE_PORT=${SERVICE_PORT}

# Build the Go project binary using the SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o ${SERVICE_NAME} .

# Ensure the binary is executable
RUN chmod +x ${SERVICE_NAME}

# Expose the service port specified by the build-time argument
EXPOSE ${SERVICE_PORT}

# Use the environment variable in the CMD for runtime execution
CMD ["sh", "-c", "./${SERVICE_NAME}"]
