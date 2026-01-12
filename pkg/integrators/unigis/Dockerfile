# /home/ubuntu/jonobridge/pkg/vecfleet/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Create directory structure
RUN mkdir -p /app/github.com/MaddSystems/jonobridge/common

# Copy only the common module
COPY pkg/common /app/github.com/MaddSystems/jonobridge/common

# Create directory for unigis and copy its files
WORKDIR /app/unigis

# Copy only the necessary files for the unigis module
COPY pkg/integrators/unigis/go.mod pkg/integrators/unigis/go.sum ./
COPY pkg/integrators/unigis/main.go ./
COPY pkg/integrators/unigis/features ./features

# Modify the replace directive in go.mod
RUN sed -i 's|replace github.com/MaddSystems/jonobridge/common => ../../common|replace github.com/MaddSystems/jonobridge/common => /app/github.com/MaddSystems/jonobridge/common|g' go.mod

# Build the unigis binary
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o unigis main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/unigis/unigis /unigis
# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/unigis","-v"]
