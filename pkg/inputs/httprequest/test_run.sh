#!/bin/bash

# Set environment variables
export SKYWAVE_ACCESS_ID="70001184"
export SKYWAVE_PASSWORD="JEUTPKKH"
export SKYWAVE_FROM_ID="13969586728"
export HTTP_POLLING_TIME="300"
export HTTP_URL="https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml"
export MQTT_BROKER_HOST="localhost"

# Build and run with verbose flag to see the output
echo "Building main.go..."
go build -o httprequest main.go

if [ $? -eq 0 ]; then
    echo "Build successful. Running with verbose mode (-v flag)..."
    echo "Press Ctrl+C to stop"
    echo "================================================"
    ./httprequest -v
else
    echo "Build failed!"
    exit 1
fi
