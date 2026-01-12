# GPS Gate HTTP Service

This service listens for MQTT messages from GPS Gate, processes them, and forwards the data to Elasticsearch. It also sends Telegram notifications for alerts.

## Environment Variables

### Required Environment Variables

```bash
# MQTT Configuration
export MQTT_BROKER_HOST="your-mqtt-broker-host"  # Required: MQTT broker hostname

# Elasticsearch Configuration
export ELASTIC_DOC_NAME="gpsgate_index"          # Elasticsearch index name (default: gpsgate_default)
```

### Optional Environment Variables

```bash
# Telegram Notification Configuration
export TELEGRAM_BOT_TOKEN="your-bot-token"       # Telegram bot token for notifications
export TELEGRAM_API_URL="https://your-api.com"   # API URL to get chat_id from app_id
export TELEGRAM_MESSAGE_HEADER="🚨 *Alerta: %s*"  # Message header format (default: "🚨 *Alerta: %s*")
export TELEGRAM_ADDITIONAL_FIELDS="GEOFENCE_NAME,POS_ADDRESS,GEOFENCE_TAG_ID"  # Comma-separated list of additional fields to include in Telegram messages

# Test Mode
export GPSGATE_TEST_TELEGRAM="Y"                 # Enable test mode for Telegram (uses test chat_id)
export GPSGATE_TEST_CHAT_ID="-1002135388607"    # Test chat_id for Telegram notifications

# EGO API (for correcting APP_ID)
export EGO_API_URL="https://your-ego-api.com"    # API URL to get correct APP_ID from IMEI

# Application Configuration
export HOSTNAME="your-pod-name"                  # Pod/container hostname
```

## TELEGRAM_MESSAGE_HEADER Format

The `TELEGRAM_MESSAGE_HEADER` environment variable allows you to customize the header of the Telegram notification message.

### Format
- A string that can include `%s` as a placeholder for the rule name
- Default: `"🚨 *Alerta: %s*"`
- Example: `"⚠️ *Warning: %s*"` or `"🚨 *Alert: %s*"`

### Behavior
- If not set, uses the default header
- The `%s` will be replaced with the `RULE_NAME` from the GPS data
- Must include `%s` if you want the rule name in the header

### Example Message Output
With `TELEGRAM_MESSAGE_HEADER="⚠️ *Warning: %s*"`:

```
⚠️ *Warning: CorteDeCorrienteExterna*

*Unidad:* R2-FOURFINAL250955
*Fecha y hora:* 2025-09-30 15:13:10 CDT
```

## TELEGRAM_ADDITIONAL_FIELDS Format

The `TELEGRAM_ADDITIONAL_FIELDS` environment variable allows you to add extra fields from the GPS data to the Telegram notification message.

### Format
- Comma-separated list of field names
- Example: `GEOFENCE_NAME,POS_ADDRESS,GEOFENCE_TAG_ID`

### Behavior
- Each field name corresponds to a key in the JSON data stored in Elasticsearch
- If the field exists and has a non-empty string value, it will be added to the message
- If the field doesn't exist or is empty, it will be skipped
- Fields are added in the order specified

### Example Message Output
With `TELEGRAM_ADDITIONAL_FIELDS="GEOFENCE_NAME,POS_ADDRESS"`:

```
🚨 *Alerta: CorteDeCorrienteExterna*

*Unidad:* R2-FOURFINAL250955
*Fecha y hora:* 2025-09-30 15:13:10 CDT
*GEOFENCE_NAME:* GPScontrol
*POS_ADDRESS:* T399E_ATT_P2-355
```

### Available Fields
Based on the JSON structure, common fields you might want to include:
- `GEOFENCE_NAME`
- `POS_ADDRESS`
- `GEOFENCE_TAG_NAME`
- `SIGNAL_SPEED`
- `SIGNAL_BATTERYLEVEL`
- `EVENT_DURATION`
- And any other field present in your GPS data

## Command Line Flags

- `-v`: Enable verbose logging

## Usage

```bash
# Basic run
./gpsgatehttp

# With verbose logging
./gpsgatehttp -v
```

## Health Check

The service provides a health check endpoint at `http://localhost:8080/health` that returns JSON with:
- Service status
- Uptime
- Time since last message
- Total messages processed

## Features

- MQTT message processing with panic recovery
- Dynamic field decoding (percent-encoded strings)
- Elasticsearch integration
- Telegram notifications with customizable additional fields
- Health monitoring and automatic reconnection
- Comprehensive logging and error handling
