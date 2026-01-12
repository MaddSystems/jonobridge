FROM golang:1.24-alpine AS builder

WORKDIR /app

# Create directory structure for the common module
RUN mkdir -p /app/common

# Copy the common module files
COPY pkg/common /app/common

# Create a directory for the backend application and set it as the working directory
WORKDIR /app/backend

# Copy the backend application files
COPY pkg/integrators/grule/backend/ .

# Modify the replace directive in go.mod to use the absolute path in the container
RUN sed -i 's|replace github.com/MaddSystems/jonobridge/common => ../../../common|replace github.com/MaddSystems/jonobridge/common => /app/common|g' go.mod

# Download dependencies and build the application
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o grule-backend

# Start a new stage from scratch for the final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/backend/grule-backend .

EXPOSE 8081

# Run the application
CMD ["./grule-backend"]
