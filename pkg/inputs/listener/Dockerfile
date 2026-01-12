# Use the official Golang image as a build environment
FROM golang:1.23 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the listener application source code
COPY ./ /app/listener/

WORKDIR /app/listener

# Build the listener binary
RUN CGO_ENABLED=0 GOOS=linux go build -o listener .

# Use a minimal base image to run the application
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /

# Copy only the binary from the builder stage
COPY --from=builder /app/listener/listener /listener

# Expose the ports that the application listens on
EXPOSE 1024/tcp
EXPOSE 1024/udp
EXPOSE 8080/tcp

# Run the listener application
ENTRYPOINT ["/listener","-v"]
