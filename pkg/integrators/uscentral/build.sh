#!/bin/bash
#!/bin/bash
set -e  # Exit on error

# Get the repository root directory (assuming script is run from its location)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# First build locally to check for any compilation errors
echo "Building application locally..."
cd "${REPO_ROOT}/pkg/integrators/uscentral"
go build

# Clean up any existing Docker images with the uscentraltrack name
echo "Removing old Docker images..."
docker images --filter=reference="*uscentral*" --format "{{.ID}}" | xargs -r docker rmi -f

# We need to build from the project root to include common module in the build context
echo "Moving to project root directory..."
cd "${REPO_ROOT}"

# Build the Docker image with context from the project root
echo "Building Docker image..."
docker build -t uscentral -f ./pkg/integrators/uscentral/Dockerfile .

# Tag and push the image
echo "Tagging Docker image..."
docker tag uscentral maddsystems/uscentral:1.0.0
echo "Pushing Docker image..."
docker push maddsystems/uscentral:1.0.0

echo "Build process completed successfully!"
