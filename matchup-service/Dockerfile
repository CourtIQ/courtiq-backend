# Use the official Golang image, version 1.23.1, based on Alpine Linux for minimal image size
# Alpine is a small Linux distribution commonly used for containers
FROM golang:1.23.1-alpine

# Set the working directory inside the container to /app
WORKDIR /app

# Install git as it is required for go mod to download dependencies from repositories
RUN apk add --no-cache git

# Copy the go.mod and go.sum files to the /app directory in the container
# These files are used to manage the Go dependencies for your project
COPY go.mod go.sum ./

# Download the Go module dependencies based on go.mod and go.sum
RUN go mod download

# Copy the entire project code from your local machine to the /app directory in the container
COPY . .

# Define build-time variables for the service name and service port
# SERVICE_NAME will be used to name the final binary, and SERVICE_PORT will be the port the service listens on
ARG SERVICE_NAME
ARG SERVICE_PORT

# Build the Go project using the provided SERVICE_NAME
# CGO_ENABLED=0 disables CGO (which allows Go to interact with C libraries)
# GOOS=linux sets the operating system target to Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o ${SERVICE_NAME} .

# Expose the SERVICE_PORT environment variable as the port that will be exposed by the container
EXPOSE ${SERVICE_PORT}

# Define the command that will be run when the container starts
# This command will execute the binary built earlier, using the SERVICE_NAME
CMD ["/bin/sh", "-c", "./${SERVICE_NAME}"]
