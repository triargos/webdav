# First stage: Build the application
FROM golang:1.22 as builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the application
COPY . .
RUN go build -o /webdav

# Second stage: Create the final runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/webdav .

# Set environment variable if needed
ENV DOCKER_ENABLED="1"

# Define the entry point for the container
ENTRYPOINT ["/app/webdav"]
