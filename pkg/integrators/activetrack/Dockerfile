# /home/ubuntu/jonobridge/pkg/vecfleet/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Create directory structure
RUN mkdir -p /app/github.com/MaddSystems/jonobridge/common

# Copy only the common module
COPY pkg/common /app/github.com/MaddSystems/jonobridge/common

# Create directory for activetrack and copy its files
WORKDIR /app/activetrack

# Copy only the necessary files for the activetrack module
COPY pkg/integrators/activetrack/go.mod pkg/integrators/activetrack/go.sum ./
COPY pkg/integrators/activetrack/main.go ./
COPY pkg/integrators/activetrack/features ./features

# Modify the replace directive in go.mod
RUN sed -i 's|replace github.com/MaddSystems/jonobridge/common => ../../common|replace github.com/MaddSystems/jonobridge/common => /app/github.com/MaddSystems/jonobridge/common|g' go.mod

# Build the activetrack binary
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o activetrack main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/activetrack/activetrack /activetrack
# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/activetrack","-v"]
