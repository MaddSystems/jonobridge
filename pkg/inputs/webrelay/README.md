# WebRelay Service

This service acts as an API endpoint that retrieves vehicle location data from server1.gpscontrol.com.mx and exposes it via a local HTTP server.

## Setup

Set the following environment variables:

```bash
export GGS_USER="admindesarrollo"
export GGS_PASSWORD="GPSc0ntr0l00"
export APP_ID=424
export WEBRELAY_TOKEN="d655eea7616e05b35dc7b22dd83b6ebc"
export PORTAL_ENDPOINT="test"
```

## Running the Service

Build and run the service:

```bash
go build -o webrelay
./webrelay
```

For verbose logging, use the -v flag:

```bash
./webrelay -v
```

## Testing Locally

### Basic Query (No Authentication)

If WEBRELAY_TOKEN is not set, you can access the endpoint without authentication:

```bash
# Query for vehicle with plates "105"
curl "http://localhost:8081/test?plates=105"

# Query for a different vehicle
curl "http://localhost:8081/test?plates=115"
```

### With Bearer Token Authentication

When WEBRELAY_TOKEN is set, you must include the token in your requests:

```bash
# Query with authentication
curl -H "Authorization: Bearer d655eea7616e05b35dc7b22dd83b6ebc" \
  "http://localhost:8081/test?plates=105"
```

### Formatting Output

For pretty-printed JSON output (if you have jq installed):

```bash
# Pretty print the output
curl "http://localhost:8081/test?plates=105" | jq

# Save response to a file
curl "http://localhost:8081/test?plates=105" -o response.json
```

### Debugging Requests

For more detailed request/response information:

```bash
# Verbose mode
curl -v -H "Authorization: Bearer d655eea7616e05b35dc7b22dd83b6ebc" \
  "http://localhost:8081/test?plates=105"
```

## Example Response

```json
{
  "imei": "867869061831673",
  "plate": "105",
  "altitude": 2351.0,
  "latitude": 19.60464,
  "longitude": -99.209981,
  "speed": 0.0,
  "heading": 264.0,
  "date": "05-19-2025",
  "time": "15:48:08",
  "moving": false,
  "ignitionStatus": true,
  "stoppingDate": "05-19-2025",
  "stopingTime": "00:03"
}
```

## Troubleshooting

1. If you get an `"error": "las placas solicitadas no existen"` response, verify that the plates number exists in the system.

2. If authentication fails, confirm that you're using the correct token in your Authorization header.

3. For connectivity issues, ensure your network allows HTTP connections to server1.gpscontrol.com.mx.
