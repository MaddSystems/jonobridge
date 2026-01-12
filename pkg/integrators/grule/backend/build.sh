#!/bin/bash
set -e

# Get the repository root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# First build locally to check for any compilation errors
echo "Building backend locally..."
CGO_ENABLED=0 go build -o grule-backend

echo "Backend built successfully."

IMAGE_NAME="maddsystems/grule"
IMAGE_TAG="1.0.0"
FULL_IMAGE_NAME="${IMAGE_NAME}:${IMAGE_TAG}"

# Clean up any existing Docker images with the grule name
echo "Removing old Docker images..."
docker images --filter=reference="*grule*" --format "{{.ID}}" | xargs -r docker rmi -f

# Build from project root to include common module in the build context
echo "Moving to project root directory..."
cd "${REPO_ROOT}"

echo "Building Docker image: ${FULL_IMAGE_NAME}"

# Docker build from repo root
docker build -f pkg/integrators/grule/backend/Dockerfile -t "${FULL_IMAGE_NAME}" .

echo "Docker image built successfully: ${FULL_IMAGE_NAME}"

# Tag and push the image
echo "Tagging Docker image..."
docker tag "${FULL_IMAGE_NAME}" "${FULL_IMAGE_NAME}"
echo "Pushing Docker image..."
docker push "${FULL_IMAGE_NAME}"

echo "Build process completed successfully!"
