# GEMINI.md - Meitrack Protocol Interpreter

## Project Overview

This project is a Go-based service that acts as a protocol interpreter for Meitrack GPS trackers. It listens for incoming data from trackers on MQTT topics, parses the proprietary Meitrack protocol, transforms it into a standardized JSON format (referred to as "Jono protocol"), and publishes the normalized data to another MQTT topic.

The service is designed to be robust, with features like a circuit breaker to prevent cascading failures, health monitoring, and graceful shutdown.

### Key Technologies

*   **Go:** The primary programming language.
*   **MQTT:** The messaging protocol used for communication with the trackers and other services. The `github.com/eclipse/paho.mqtt.golang` library is used.
*   **Docker:** The project includes a `Dockerfile` for containerization.

### Architecture

The project is structured into three main components:

1.  **`main.go`:** The application's entry point. It manages the MQTT client, including connection, subscriptions, and message handling. It orchestrates the data processing pipeline.
2.  **`features/meitrack_protocol`:** This package is responsible for parsing the raw, proprietary Meitrack protocol. It contains logic to handle different Meitrack command types (e.g., `AAA`, `CCE`, `CCC`).
3.  **`features/jono`:** This package takes the parsed data from the `meitrack_protocol` package and normalizes it into a structured, consistent JSON format.

## Building and Running

### Prerequisites

*   Go
*   Docker
*   An MQTT broker (e.g., Mosquitto)

### Environment Variables

*   `MQTT_BROKER_HOST`: The hostname or IP address of the MQTT broker.

### Building and Running with Docker

The `README.md` provides instructions for building and running the service as a Docker container.

```bash
# Build the Docker image
docker build -t meitrackprotocol -f ./Dockerfile .

# Run the Docker container
docker run -e MQTT_BROKER_HOST=<your_mqtt_broker_host> meitrackprotocol
```

### Building and Running Locally

The `build.sh` script can be used to build the application.

```bash
# Build the executable
./build.sh

# Run the executable
MQTT_BROKER_HOST=<your_mqtt_broker_host> ./meitrackprotocol
```

### Testing

The `README.md` provides example `mosquitto_pub` and `mosquitto_sub` commands for testing the service.

**Publishing a raw tracker message:**

```bash
mosquitto_pub -h localhost -t tracker/raw -m '$$L172,864765045580768,AAA,35,19.400860,-98.927070,240122233943,A,9,16,0,298,0.8,2260,106904870,84774230,334|20|32D2|050ACCB5,0000,0000|0000|0000|0191|04E4,00000001,,3,,,0,0*6C'
```

**Subscribing to the normalized output:**

```bash
mosquitto_sub -h localhost -t tracker/jonoprotocol
```

## Development Conventions

*   **Logging:** The service uses the standard `log` package. Verbose logging can be enabled with the `-v` flag.
*   **Error Handling:** Errors are generally handled by logging them and, in some cases, incrementing an error counter for health monitoring.
*   **Code Structure:** The code is organized into feature packages (`meitrack_protocol`, `jono`) with a clear separation of concerns. Each feature has its own `models` and `usecases` sub-packages.
