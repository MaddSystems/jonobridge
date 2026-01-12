# Environment Variables for HTTPrequest

Http request make a request every polling time to get messages from server
then publisk to mqtt broker using topic http/get

## Required Environment Variables

```bash
# Required: HTTP FEED for example
export HTTP_URL="https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/0BkM9B2i01vF8eigoq3T1XO5HgMfQmfQa/message.xml"

# Required: MQTT Broker Configuration
export MQTT_BROKER_HOST="mosquitto"  
```

## Optional Environment Variables

```bash
# Application Configuration

export HTTP_POLLING_TIME="30"    # Polling interval in seconds (default: 30)
```

## Command Line Flags

The following command line flags are available:

- `-v`: Enable verbose output (debug mode)

## Logging

- Error messages are always logged with timestamps
- Debug messages (with `-v` flag) show:
  - Database operations
  - API requests
  - Message processing
  - Server communications

## Error Handling

The application will:
1. Exit if required environment variables are missing
2. Log errors with timestamps
3. Continue processing other messages if one fails
4. Automatically create required database tables
