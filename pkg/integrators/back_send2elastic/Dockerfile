# /home/ubuntu/jonobridge/pkg/send2elastic/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum from the project root
COPY go.mod go.sum ./
RUN go mod download

# Copy the send2elastic source code from pkg/send2elastic
COPY ./ /app/send2elastic

WORKDIR /app/send2elastic

# Build the send2elastic binary
RUN CGO_ENABLED=0 GOOS=linux go build -o send2elastic main.go

# Create a minimal image with just the compiled binary
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /
# Copy only the binary from the builder stage
COPY --from=builder /app/send2elastic/send2elastic /send2elastic
# Expose the ports that the application listens on
EXPOSE 1883

ENTRYPOINT ["/send2elastic","-v"]
