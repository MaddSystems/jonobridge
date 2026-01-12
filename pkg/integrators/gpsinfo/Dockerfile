# /home/ubuntu/jonobridge/pkg/vecfleet/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Install timezone data for Debian-based golang image
RUN apt-get update && apt-get install -y --no-install-recommends tzdata

# Create directory structure
RUN mkdir -p /app/github.com/MaddSystems/jonobridge/common

# Copy only the common module
COPY pkg/common /app/github.com/MaddSystems/jonobridge/common

# Create directory for gpsinfo and copy its files
WORKDIR /app/gpsinfo

# Copy only the necessary files for the gpsinfo module
COPY pkg/integrators/gpsinfo/go.mod pkg/integrators/gpsinfo/go.sum ./
COPY pkg/integrators/gpsinfo/main.go ./
COPY pkg/integrators/gpsinfo/features ./features

# Modify the replace directive in go.mod
RUN sed -i 's|replace github.com/MaddSystems/jonobridge/common => ../../common|replace github.com/MaddSystems/jonobridge/common => /app/github.com/MaddSystems/jonobridge/common|g' go.mod

# Build the gpsinfo binary
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o gpsinfo main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Install tzdata in the Alpine final image
RUN apk add --no-cache tzdata

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/gpsinfo/gpsinfo /gpsinfo

# Copy the timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set the timezone environment variable
ENV TZ=UTC

# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/gpsinfo","-v"]