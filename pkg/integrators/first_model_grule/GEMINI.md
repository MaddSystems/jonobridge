This is a Go-based universal audit system for a Grule Rule Engine, with a Flask frontend. It uses MySQL for data storage and MQTT for message passing.

The backend is a Go application that exposes a REST API for managing rules and viewing audit data. The core of the backend is a worker pool that processes messages for each IMEI concurrently. For each message, it creates a `PacketWrapper` and executes the rules against it. The `PacketWrapper` contains the data from the message, as well as flags that can be used to control the execution flow of the rules.

The frontend is a Flask application that provides a web interface for the Grule audit system. It communicates with the Go backend via a REST API. The application allows users to view and manage rules, as well as to view audit data.

### Building and Running

**Backend (Go)**

To build the backend, run the following command:

```
go build -o grule-engine
```

To run the backend, you need to set the following environment variables:

-   `GRULE_AUDIT_ENABLED`: `Y` or `N`
-   `GRULE_AUDIT_LEVEL`: `ALL` or `ALERT_ONLY`
-   `API_PORT`: The port for the backend API (default: `8080`)
-   `MYSQL_HOST`: The MySQL host
-   `MYSQL_USER`: The MySQL user
-   `MYSQL_PASSWORD`: The MySQL password
-   `MYSQL_DATABASE`: The MySQL database
-   `TELEGRAM_BOT_TOKEN`: The Telegram bot token
-   `MQTT_BROKER_HOST`: The MQTT broker host

Then, you can run the backend with the following command:

```
./grule-engine
```

**Frontend (Flask)**

To run the frontend, first install the dependencies:

```
cd external-web
pip install -r requirements.txt
```

Then, set the following environment variables:

-   `GRULE_API_URL`: The URL of the backend API
-   `FLASK_PORT`: The port for the frontend application (default: `5001`)

Finally, run the frontend with the following command:

```
python main.py
```

### Project Manual

**Location:** `Manual/`

This directory contains the complete documentation for the Grule Rule Engine, including the Grule Rule Language (GRL) syntax, available functions, and advanced concepts like the RETE algorithm.

**It is critical to read and understand the documents in this directory before attempting to modify any rules or Go code that interacts with the rule engine.** This manual is the primary source of truth for the engine's behavior and conventions. Failure to consult it may lead to incorrect assumptions and bugs.

### Development Conventions

The codebase is well-structured and follows standard Go and Python conventions. The code is well-documented with comments. The project uses a multi-stage Dockerfile for building and deploying the application. It also includes a Kubernetes deployment file.
