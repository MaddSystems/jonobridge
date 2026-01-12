# SkyWave/ORBCOMM Environment Variables - Part 1

This document describes the new environment variables added for SkyWave/ORBCOMM integration.

## Environment Variables

### Required SkyWave Variables (all must be set to enable SkyWave mode):

- **SKYWAVE_ACCESS_ID**: SkyWave account access ID (e.g., `70001184`)
- **SKYWAVE_PASSWORD**: SkyWave API password (e.g., `JEUTPKKH`)
- **SKYWAVE_FROM_ID**: Starting message ID for retrieval (e.g., `13969586728`)

### Existing Variables (still required):

- **MQTT_BROKER_HOST**: MQTT broker hostname/IP
- **HTTP_POLLING_TIME**: Polling interval in seconds (default: 30)

### Optional Variables (for fallback HTTP mode):

- **HTTP_URL**: URL for HTTP GET requests (only needed when SkyWave mode is disabled)

## Operation Modes

### SkyWave Mode (when all SkyWave variables are set):
- Fetches data from SkyWave API at `https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/`
- Publishes raw XML data (as hex) to MQTT topic: `skywave/xml`
- Uses the configured polling interval

### HTTP Mode (fallback when SkyWave variables are missing):
- Fetches data from HTTP_URL
- Publishes data (as hex) to MQTT topic: `http/get`
- Uses the configured polling interval

## Usage Examples

### Enable SkyWave Mode:
```bash
export SKYWAVE_ACCESS_ID=70001184
export SKYWAVE_PASSWORD=JEUTPKKH
export SKYWAVE_FROM_ID=13969586728
export MQTT_BROKER_HOST=localhost
export HTTP_POLLING_TIME=180

./httprequest -v
```

### HTTP Mode (existing behavior):
```bash
export HTTP_URL=https://api.example.com/data
export MQTT_BROKER_HOST=localhost
export HTTP_POLLING_TIME=30

./httprequest -v
```

## MQTT Topics

- `skywave/xml`: Raw SkyWave XML data (hex encoded) - for Part 2 processing
- `http/get`: Raw HTTP response data (hex encoded) - existing functionality

## Next Steps (Part 2)

Part 2 will involve creating another Go program that:
1. Subscribes to the `skywave/xml` MQTT topic
2. Decodes the hex data back to XML
3. Parses the XML to extract position data
4. Converts coordinates to MVT366 format
5. Sends MVT366 messages via UDP to tracking server



export SKYWAVE_ACCESS_ID="70001184"
export SKYWAVE_PASSWORD="JEUTPKKH"
export SKYWAVE_FROM_ID="13969586728"
export HTTP_POLLING_TIME="300"
export HTTP_URL="https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml"
export MQTT_BROKER_HOST="localhost"