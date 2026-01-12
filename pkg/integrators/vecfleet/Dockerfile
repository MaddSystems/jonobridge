# /home/ubuntu/jonobridge/pkg/vecfleet/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum from the project root
COPY go.mod go.sum ./
RUN go mod download

# Copy the vecfleet source code from pkg/vecfleet
COPY ./ /app/vecfleet

WORKDIR /app/vecfleet

# Build the vecfleet binary
RUN CGO_ENABLED=0 GOOS=linux go build -o vecfleet main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/vecfleet/vecfleet /vecfleet
# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/vecfleet","-v"]
