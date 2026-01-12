#!/bin/bash
#!/bin/bash
set -e  # Exit on error

# Get the repository root directory (assuming script is run from its location)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# First build locally to check for any compilation errors
echo "Building application locally..."
cd "${REPO_ROOT}/pkg/integrators/unigis"
go build

# Clean up any existing Docker images with the unigis name
echo "Removing old Docker images..."
docker images --filter=reference="*unigis*" --format "{{.ID}}" | xargs -r docker rmi -f

# We need to build from the project root to include common module in the build context
echo "Moving to project root directory..."
cd "${REPO_ROOT}"

# Build the Docker image with context from the project root
echo "Building Docker image..."
docker build -t unigis -f ./pkg/integrators/unigis/Dockerfile .

# Tag and push the image
echo "Tagging Docker image..."
docker tag unigis maddsystems/unigis:1.0.0
echo "Pushing Docker image..."
docker push maddsystems/unigis:1.0.0

echo "Build process completed successfully!"
