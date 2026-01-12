# Initial Concept
A GPS Fleet Rules Engine ("Grule Engine") system. It consists of a high-performance Go backend that processes MQTT messages using the Grule rule engine and a Python/Flask frontend for rule management and audit visualization.

# Product Guide

## 1. Project Overview
The **Grule Engine** is a universal, high-performance rule engine designed specifically for GPS fleet data. It serves as a real-time telematics security and anti-theft decision engine, sitting between raw telematics protocols and fleet management operations. Unlike simple alert systems, Grule Engine emphasizes forensics, developer-friendly extensibility, and active defense capabilities.

The system ingests GPS data from multiple vendors, normalized into a single **Jono Protocol** via MQTT. It processes this data using a stateful, per-device pipeline that supports complex detection strategies—such as jamming detection, fuel theft patterns, and geofence violations—while providing a detailed, auditable timeline of every decision made.

It is designed with extremely strong forensic capabilities, where every decision can be replayed frame-by-frame. The worker-per-IMEI model ensures high availability and survives partial failures, while the battle-tested concurrency model prevents race conditions even under high packet rates. Furthermore, the system is open for extension, allowing new detection strategies to be added via simple `.grl` files.

## 2. Target Users
*   **Fleet Managers:** Monitor vehicle alerts, status, and security incidents in real-time.
*   **System Administrators:** Manage rule logic, configurations, and system health.
*   **Developers & Analysts:** Debug rule execution flows, analyze performance, and create new detection logic without altering the core codebase.
*   **Security Operations:** React to critical alerts with active interventions (e.g., cutting the engine).

## 3. Key Features
*   **Real-Time Decision Engine:**
    *   **Per-Device Sequential Processing:** Dedicated worker goroutines ensure strict ordering and eliminate race conditions for each IMEI.
    *   **Hot-Reloadable Logic:** Business logic is decoupled into `.grl` files (Grule), allowing for updates without recompiling the Go backend.
    *   **Stateful Analysis:** Maintains a memory buffer (e.g., last 10 valid positions) and persistent counters per device to detect complex patterns over time (e.g., 90-minute windows).
*   **Advanced Security & Anti-Theft:**
    *   **Active Defense:** Capable of triggering active responses, such as engine cuts or relay activations, via MQTT commands.
    *   **Geofence Integration:** Supports safe-zone exclusions (workshops, secure ports) to minimize false positives.
*   **Comprehensive Observability & Audit:**
    *   **Declarative Audit System:** Audit logic is decoupled from business rules using YAML manifests, ensuring rules remain pure while providing rich, automated execution tracking.
    *   **Universal Alert Deduplication:** Prevents duplicate alerts (intra-packet and inter-packet) using atomic check-and-set guards and local cycle caching, ensuring exactly one alert per event.
    *   **"Movie Mode" Timeline:** A granular, frame-by-frame timeline of rule executions, buffer states, and decision contexts for forensic analysis.
    *   **Visual Audit Tools:** Web-based interface (using jqGrid) for inspecting execution history and verifying why an alert was or wasn't triggered.
*   **Scalable Architecture:**
    *   **Protocol Abstraction:** Consumes the vendor-neutral **Jono Protocol** via MQTT.
    *   **Horizontal Scalability:** Worker-per-IMEI model supports distribution across multiple instances.

## 4. Technical Architecture
*   **Backend (Engine):** Written in **Go (Golang)** for high performance and concurrency. It handles MQTT ingestion, state management, and rule execution using the **Grule-Rule-Engine**.
*   **Frontend (UI):** Built with **Python** and **Flask**, providing a user-friendly interface for managing rules and visualizing audit logs.
*   **Communication:**
    *   **Input:** MQTT for receiving normalized Jono Protocol GPS data.
    *   **Output:** MQTT for sending commands; Multi-channel alerting via Telegram, Webhooks, etc.
*   **Persistence:** **MySQL** is used for storing rule definitions, geofences, and audit logs.
*   **Separation of Concerns:** Strict separation between the Go backend (engine) and the Python frontend (UI), ensuring modularity and easier maintenance.
