# /home/ubuntu/jonobridge/pkg/vecfleet/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Create directory structure
RUN mkdir -p /app/github.com/MaddSystems/jonobridge/common

# Copy only the common module
COPY pkg/common /app/github.com/MaddSystems/jonobridge/common

# Create directory for recursoconfiable and copy its files
WORKDIR /app/recursoconfiable

# Copy only the necessary files for the recursoconfiable module
COPY pkg/integrators/recursoconfiable/go.mod pkg/integrators/recursoconfiable/go.sum ./
COPY pkg/integrators/recursoconfiable/main.go ./
COPY pkg/integrators/recursoconfiable/features ./features
COPY pkg/integrators/recursoconfiable/utils ./utils

# Modify the replace directive in go.mod
RUN sed -i 's|replace github.com/MaddSystems/jonobridge/common => ../../common|replace github.com/MaddSystems/jonobridge/common => /app/github.com/MaddSystems/jonobridge/common|g' go.mod

# Build the recursoconfiable binary
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o recursoconfiable main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/recursoconfiable/recursoconfiable /recursoconfiable
# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/recursoconfiable","-v"]
