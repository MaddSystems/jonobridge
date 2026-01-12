#!/bin/bash

# ORBCOMM MQTT Test Setup Script
# This script sets up the environment and runs the ORBCOMM to MQTT test

echo "Setting up ORBCOMM MQTT Test Environment..."

# Install Python dependencies
echo "Installing Python dependencies..."
pip3 install -r requirements.txt

# Make the Python script executable
chmod +x orbcomm_mqtt_test.py

echo "Setup complete!"
echo ""
echo "Usage examples:"
echo "1. Run once to test:"
echo "   python3 orbcomm_mqtt_test.py --once --verbose"
echo ""
echo "2. Run continuously (every 60 seconds):"
echo "   python3 orbcomm_mqtt_test.py --interval 60"
echo ""
echo "3. Custom MQTT broker:"
echo "   python3 orbcomm_mqtt_test.py --mqtt-host your-broker.com --mqtt-port 1883"
echo ""
echo "4. Different topic:"
echo "   python3 orbcomm_mqtt_test.py --topic tracker/from-udp"
echo ""
echo "To monitor MQTT messages:"
echo "   mosquitto_sub -h localhost -t 'tracker/from-tcp' -v"
echo "   mosquitto_sub -h localhost -t 'tracker/from-tcp/debug' -v"
