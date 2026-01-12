# Use the official Golang image as a build environment
FROM golang:1.23 AS builder

# Set the working directory in the container
WORKDIR /app

# Install swag CLI for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the proxy application source code
COPY ./ /app/proxy/

# Run swag init to generate Swagger docs
WORKDIR /app/proxy
RUN /go/bin/swag init

# Build the proxy binary
RUN CGO_ENABLED=0 GOOS=linux go build -o proxy .

# Use a minimal base image to run the application
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /

# Copy only the binary from the builder stage
COPY --from=builder /app/proxy/proxy /proxy

# Copy the generated Swagger docs
COPY --from=builder /app/proxy/docs /docs

# Expose the ports that the application listens on
EXPOSE 1024 8080

# Run the proxy application
ENTRYPOINT ["/proxy", "-v"]
