#!/bin/bash

# Navigate to the project root (assuming script is in tests/)
cd "$(dirname "$0")/.."

echo "Stopping any existing grule-backend processes..."
pkill -f "grule-backend" || true
pkill -f "go run main.go" || true

echo "Starting Backend with MYSQL_DB=grule..."
cd backend

# Build the executable
echo "Building backend..."
go build -o grule-backend

# Run the executable with verbose flag
env MYSQL_DB=grule ./grule-backend -v

PID=$!
echo "Backend started with PID $PID"
echo "Press Ctrl+C to stop the backend"
