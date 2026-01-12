# /home/ubuntu/jonobridge/pkg/httpinput/Dockerfile
FROM golang:1.23 AS builder

# Set the base working directory
WORKDIR /src

# Copy Go module files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the directories preserving the structure the Go imports expect
COPY ./pkg/ ./pkg/
COPY ./common/ ./common/

# Set working directory to httpinput component
WORKDIR /src/pkg/inputs/httpinput

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o httpinput

# Create minimal runtime image
FROM alpine:latest

# Copy the binary from builder
COPY --from=builder /src/pkg/inputs/httpinput/httpinput /httpinput

# Expose HTTP port
EXPOSE 8081

# Run with verbose flag enabled
ENTRYPOINT ["/httpinput", "-v"]
