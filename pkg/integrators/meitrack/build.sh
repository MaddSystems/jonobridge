#!/bin/bash
set -e

echo "Building meitrack integrator..."

# Clean up existing Docker images
echo "Cleaning up previous Docker images..."
docker images --filter=reference="*meitrack*" --format "{{.ID}}" | xargs docker rmi -f 2>/dev/null || true

# Move to the project root to ensure proper context for build
cd /home/ubuntu/jonobridge

# Build new Docker image from the project root with the context pointing to the project root
echo "Building Docker image..."
docker build -t meitrack -f ./pkg/integrators/meitrack/Dockerfile .

# Tag and push the image
echo "Tagging and pushing Docker image..."
docker tag meitrack maddsystems/meitrack:1.0.0
docker push maddsystems/meitrack:1.0.0

echo "Build complete!"

# Add a function to run the container locally for testing
run_local() {
  echo "Running container locally for testing..."
  
  # Parse additional arguments to pass to the container
  CONTAINER_ARGS=""
  if [ "$#" -gt 0 ]; then
    CONTAINER_ARGS="$@"
  fi
  
  # Create and start the container with logs directed to stdout
  docker run --rm -it \
    -e LOG_LEVEL=debug \
    -e LOG_FORMAT=json \
    -e DEBUG=true \
    meitrack $CONTAINER_ARGS
}

# If the script is called with "run" argument, run the container locally
if [ "$1" == "run" ]; then
  shift  # Remove "run" from the arguments
  run_local "$@"  # Pass remaining arguments to run_local
elif [ "$1" == "debug" ]; then
  # Run with debug mode
  shift  # Remove "debug" from the arguments
  docker run --rm -it \
    -e LOG_LEVEL=debug \
    -e LOG_FORMAT=text \
    -e DEBUG=true \
    --entrypoint /bin/sh \
    meitrack
else
  # Usage information
  if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    echo "Usage: $0 [OPTION] [ARGS...]"
    echo "Options:"
    echo "  run [ARGS]  Build and run the container locally for testing, passing ARGS to the application"
    echo "  debug       Run container with interactive shell for debugging"
    echo "  -h, --help  Display this help message"
  fi
fi