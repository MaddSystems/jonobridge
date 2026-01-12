# JonoBridge Frontend

JonoBridge is a comprehensive web-based platform for managing GPS tracking services deployed on Kubernetes. It provides an intuitive interface for configuring, deploying, and monitoring GPS tracking infrastructure.

## What It Does

This Flask web application allows users to:

- **Manage Clients**: Create and manage separate Kubernetes namespaces for different clients
- **Configure Services**: Set up GPS tracking workflows using three types of services:
  - **Inputs**: Receive data from GPS trackers (TCP/UDP listeners, HTTP endpoints)
  - **Interpreters**: Process and translate data between different GPS protocols (Meitrack, Queclink, Suntech, etc.)
  - **Integrators**: Send processed data to external systems (MySQL, Elasticsearch, cloud platforms)
- **Deploy to Kubernetes**: Automatically generate and deploy Kubernetes manifests for complete service stacks
- **Monitor Systems**: View pod status, logs, and system health in real-time
- **Admin Functions**: Database management, user authentication, system configuration
- **Tracker Management**: Monitor connected GPS devices and send commands via API

## Key Features

- User authentication with role-based access
- Drag-and-drop service configuration interface
- Automatic Kubernetes manifest generation
- Mosquitto MQTT broker integration for message routing
- MySQL database for configuration and metadata
- REST API for generating GPS protocol commands (GT06, BSJ)
- WhatsApp integration for system alerts
- Comprehensive logging and monitoring

## Technology Stack

- **Backend**: Python Flask
- **Database**: MySQL
- **Container Orchestration**: Kubernetes
- **Messaging**: Mosquitto MQTT
- **Frontend**: HTML/CSS/JavaScript with Jinja2 templates

## Architecture

The application uses a modular service architecture where each GPS tracking component is a Docker container. Services are dynamically loaded from the `services/` directory and can be combined to create custom data processing pipelines for different clients.

## User Interface Overview

Based on the HTML templates, the application provides:

- **Navigation**: Left sidebar with links to Clients, Trackers, and Admin (for admins)
- **Client Management**: Table view of clients with actions (Setup, Deploy, Stop, Status, Delete)
- **Service Configuration**: Visual workflow builder with drag-and-drop interface for arranging Input → Interpreter → Integration services
- **Deployment Interface**: Review and execute Kubernetes deployments with progress tracking
- **Monitoring**: Real-time status of pods, services, and logs
- **Tracker Dashboard**: Table of connected GPS devices with command sending capabilities
- **Admin Panel**: Database management and system monitoring controls

