#!/bin/bash
set -e  # Exit on error

# Get the repository root directory (assuming script is run from its location)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# First build locally to check for any compilation errors
echo "Building application locally..."
cd "${REPO_ROOT}/pkg/interpreters/gpsgatehttp"
go build

# Clean up any existing Docker images with the gpsgatehttp name
echo "Removing old Docker images..."
docker images --filter=reference="*gpsgatehttp*" --format "{{.ID}}" | xargs -r docker rmi -f

# We need to build from the project root to include common module in the build context
echo "Moving to project root directory..."
cd "${REPO_ROOT}"

# Build the Docker image with context from the project root
echo "Building Docker image..."
docker build -t gpsgatehttp -f ./pkg/interpreters/gpsgatehttp/Dockerfile .


# Use a fixed, stable tag
TAG="1.0.0"

echo "Tagging Docker image with fixed tag: ${TAG}"
docker tag gpsgatehttp maddsystems/gpsgatehttp:${TAG}

echo "Pushing Docker image maddsystems/gpsgatehttp:${TAG}..."
docker push maddsystems/gpsgatehttp:${TAG}

echo "Built and pushed: maddsystems/gpsgatehttp:${TAG}"

echo "Build process completed successfully!"
