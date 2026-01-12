# Technology Stack

## Backend
*   **Language:** Go (Golang)
*   **Core Engine:** Grule-Rule-Engine
*   **Messaging:** Eclipse Paho MQTT client
*   **Database Driver:** Go-MySQL-Driver
*   **Concurrency:** Standard library (`sync`, `context`, `atomic`), Worker Pool Pattern
*   **Architecture:**
    *   Modular Capability-based architecture (Geofence, Buffer, Metrics, Timing, Alerts)
    *   Declarative Audit System (YAML Manifests, Listener-based capture)

## Frontend
*   **Language:** Python
*   **Web Framework:** Flask
*   **Templating:** Jinja2
*   **UI Framework:** Bootstrap 5
*   **JavaScript Libraries:**
    *   Vue 3 (Dynamic timelines and modals)
    *   jQuery
    *   jqGrid (Free version 4.15.4)
    *   DataTables
*   **Icons:**
    *   Bootstrap Icons
    *   Font Awesome
*   **Styling:** Custom CSS (Gradients, shadows)

## Messaging & Data
*   **Telemetry & Commands:** MQTT
*   **Persistence:** MySQL (Rules, Audit Logs, Progress Tracking, Geofences)

## Infrastructure
*   **Containerization:** Docker
*   **Reverse Proxy:** Nginx (Implied)
