# /home/ubuntu/jonobridge/pkg/integrators/meitrack/Dockerfile
FROM golang:1.23 AS builder

# Set working directory for the build
WORKDIR /go/src/github.com/MaddSystems/jonobridge

# Copy the entire project
COPY . .

# Build the meitrack binary
WORKDIR /go/src/github.com/MaddSystems/jonobridge/pkg/integrators/meitrack

# First display the go.mod content for debugging
RUN cat go.mod 

# Build a static binary with no CGO dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o meitrack

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /go/src/github.com/MaddSystems/jonobridge/pkg/integrators/meitrack/meitrack .

# Make sure the binary is executable
RUN chmod +x /app/meitrack

# Set environment variables for logging
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json
ENV DEBUG=false

# Set the entrypoint to run the binary with verbose flag enabled by default
ENTRYPOINT ["/app/meitrack", "-v"]
